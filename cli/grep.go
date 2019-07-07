package cli

import (
	"bytes"
	"fmt"
	"math"
	"runtime"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dabdada/s3-grep/config"
	thisS3 "github.com/dabdada/s3-grep/s3"
)

var maxExcerptLength = 120

type grepResult struct {
	Key     string
	LineNum int
	Excerpt []byte
}

// Grep in objects in a S3 bucket
func Grep(session *config.AWSSession, bucketName string, prefix string, query string, ignoreCase bool) {
	svc := s3.New(session.Session)
	objects := make(chan thisS3.StoredObject, runtime.NumCPU())
	listObjectsErrors := make(chan error)
	listObjectsDone := make(chan bool)
	grepResults := make(chan *grepResult, runtime.NumCPU())
	objectProcessed := make(chan bool)

	objectsCount := 0
	objectsProcessed := 0
	allObjectsListed := false

	go thisS3.ListObjects(svc, bucketName, prefix, objects, listObjectsErrors, listObjectsDone)

	for {
		select {
		case object := <-objects:
			objectsCount++
			go grepInObjectContent(session, bucketName, object, query, ignoreCase, grepResults, objectProcessed)
		case err := <-listObjectsErrors:
			fmt.Printf("%s\n", err)
			return
		case <-listObjectsDone:
			allObjectsListed = true
		case grepResult := <-grepResults:
			fmt.Printf("s3://%s/%s %d:%s\n", bucketName, grepResult.Key, grepResult.LineNum, grepResult.Excerpt)
		case <-objectProcessed:
			objectsProcessed++
		default:
			if (objectsCount == objectsProcessed) && allObjectsListed {
				close(listObjectsErrors)
				close(objects)
				close(grepResults)
				close(objectProcessed)
				return
			}
		}
	}
}

// Grep within the content of a single S3 object
func grepInObjectContent(session *config.AWSSession, bucketName string, object thisS3.StoredObject,
	query string, ignoreCase bool, results chan<- *grepResult, processed chan<- bool) {
	content, numBytes, err := object.GetContent(session, bucketName)
	if err != nil {
		fmt.Printf("%s:%s\n", err, object.GetKey())
	} else if numBytes > 0 {
		for i, line := range bytes.Split(content, []byte("\n")) {
			if caseAwareContains(line, []byte(query), ignoreCase) {
				results <- &grepResult{
					Key:     object.GetKey(),
					LineNum: i + 1,
					Excerpt: getContentExcerpt(line, []byte(query)),
				}
			}
		}
	}
	processed <- true
}

// Get a Excerpt of a byte array
//
// If the line is not maxExcerptLength long, the whole text will be returned.
// Otherwise a 120 char excerpt is returned.
func getContentExcerpt(text []byte, query []byte) []byte {
	textLenght := len(text)
	if textLenght <= maxExcerptLength {
		return text
	}
	queryLength := len(query)
	excerptLengthLeftAndRight := (maxExcerptLength - queryLength) / 2
	index := bytes.Index(text, query)
	from := int(math.Max(float64(index-excerptLengthLeftAndRight), 0))

	// Do not cut in the middle of words.
	if text[from] == byte(' ') {
		from++
	} else if from != 0 {
		from = bytes.Index(text[from:textLenght], []byte(" ")) + 1 + from
	}

	to := int(math.Min(float64(index+queryLength+excerptLengthLeftAndRight), float64(textLenght)))
	if to != textLenght {
		offset := bytes.Index(text[to:textLenght], []byte(" "))
		if offset < 0 {
			to = textLenght
		} else {
			to += offset
		}
	}

	return text[from:to]
}

// A case aware contains function for byte arrays
func caseAwareContains(b []byte, sub []byte, ignoreCase bool) bool {
	if ignoreCase {
		return bytes.Contains(bytes.ToUpper(b), bytes.ToUpper(sub))
	}
	return bytes.Contains(b, sub)
}
