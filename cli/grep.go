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

// Grep in objects in a thisS3 bucket
func Grep(session *config.AWSSession, bucketName string, query string, ignoreCase bool) {
	svc := s3.New(session.Session)

	objects, err := thisS3.ListObjects(svc, bucketName)

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
func partitionS3Objects(objects []string, desiredPartitionNum int) [][]string {
	var divided [][]string
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
func grepInObjectContent(session *config.AWSSession, bucketName string, objects []string,
	query string, ignoreCase bool, results chan<- grepResult, done chan<- int) {
	for _, object := range objects {
		content, numBytes, err := thisS3.GetObjectContent(session, bucketName, object)
		if err != nil {
			fmt.Printf("%s:%s\n", err, object)
		} else if numBytes > 0 {
			for i, line := range bytes.Split(content, []byte("\n")) {
				if caseAwareContains(line, []byte(query), ignoreCase) {
					results <- grepResult{
						Key:     object,
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
	index := bytes.Index(text, query)
	from := int(math.Max(float64(index)-10, 0))
	to := index + len(query) + 10

	return text[from:to]
}

// A case aware contains function for byte arrays
func caseAwareContains(b []byte, sub []byte, ignoreCase bool) bool {
	if ignoreCase {
		return bytes.Contains(bytes.ToUpper(b), bytes.ToUpper(sub))
	} else {
		return bytes.Contains(b, sub)
	}
}
