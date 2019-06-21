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

type grepResult struct {
	Key     string
	LineNum int
	Excerpt []byte
}

// Grep in objects in a S3 bucket
func Grep(session *config.AWSSession, bucketName string, prefix string, query string, ignoreCase bool) {
	svc := s3.New(session.Session)

	objects, err := thisS3.ListObjects(svc, bucketName, prefix)

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	results := make(chan grepResult)
	done := make(chan int)

	dividedObjects := partitionS3Objects(objects, runtime.NumCPU()-1)
	for _, chunk := range dividedObjects {
		go grepInObjectContent(session, bucketName, chunk, query, ignoreCase, results, done)
	}

	finished := 0
	for {
		select {
		case result := <-results:
			fmt.Printf("s3://%s/%s %s:%d\n", bucketName, result.Key, result.Excerpt, result.LineNum)
		case i := <-done:
			finished += i
		default:
			if finished == len(dividedObjects) {
				close(results)
				close(done)
				return
			}
		}
	}

}

// Divide list of objects in bucket to desiredPartitionNum same sized chunks, for concurrent processing
func partitionS3Objects(objects []thisS3.StoredObject, desiredPartitionNum int) [][]thisS3.StoredObject {
	var divided [][]thisS3.StoredObject
	numObjects := len(objects)
	chunkSize := (numObjects + desiredPartitionNum - 1) / desiredPartitionNum

	for i := 0; i < numObjects; i += chunkSize {
		end := i + chunkSize

		if end > numObjects {
			end = numObjects
		}

		divided = append(divided, objects[i:end])
	}

	return divided
}

// Grep within the content of a single S3 object
func grepInObjectContent(session *config.AWSSession, bucketName string, objects []thisS3.StoredObject,
	query string, ignoreCase bool, results chan<- grepResult, done chan<- int) {
	for _, object := range objects {
		content, numBytes, err := object.GetContent(session, bucketName)
		if err != nil {
			fmt.Printf("%s:%s\n", err, object.GetKey())
		} else if numBytes > 0 {
			for i, line := range bytes.Split(content, []byte("\n")) {
				if caseAwareContains(line, []byte(query), ignoreCase) {
					results <- grepResult{
						Key:     object.GetKey(),
						LineNum: i + 1,
						Excerpt: getContentExcerpt(line, []byte(query)),
					}
				}
			}
		}
	}
	done <- 1
}

// Get a small Excerpt of a byte array
//
// 10 chars before and after the substring
func getContentExcerpt(text []byte, query []byte) []byte {
	queryLength := float64(len(query))
	index := float64(bytes.Index(text, query))
	from := int(math.Max(index-10, 0))
	to := int(math.Min(float64(index+queryLength+10), float64(len(text))))

	return text[from:to]
}

// A case aware contains function for byte arrays
func caseAwareContains(b []byte, sub []byte, ignoreCase bool) bool {
	if ignoreCase {
		return bytes.Contains(bytes.ToUpper(b), bytes.ToUpper(sub))
	}
	return bytes.Contains(b, sub)
}
