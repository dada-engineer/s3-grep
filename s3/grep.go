package s3

import (
	"bytes"
	"fmt"
	"math"
	"runtime"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dabdada/s3-grep/config"
)

type GrepResult struct {
	Key string
	LineNum int
	Excerpt []byte
}

// Grep in objects in a thisS3 bucket
func Grep(session *config.AWSSession, bucketName string, query string) {
	svc := s3.New(session.Session)

	objects, err := listObjects(svc, bucketName)

	if err != nil {
		fmt.Errorf("%s\n", err)
		return
	}

	results := make(chan GrepResult)
	done := make(chan int)

	dividedObjects := divideObjects(objects)
	for _, chunk := range dividedObjects {
		go grepInObjectContent(session, bucketName, chunk, query, results, done)
	}

	finished := 0
	for {
		select {
		case result := <-results:
			fmt.Printf("%s %s:%d\n", result.Key, result.Excerpt, result.LineNum)
		case i := <- done:
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

// Divide list of objects in bucket to NumCPU same sized chunks, for concurrent processing
func divideObjects(objects []string) [][]string {
	var divided [][]string
	numCPU := runtime.NumCPU() - 1
	numObjects := len(objects)
	chunkSize := (numObjects + numCPU - 1) / numCPU

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
						 query string, results chan<- GrepResult, done chan<- int) {
	for _, object := range objects {
		content, numBytes, err := getObjectContent(session, bucketName, object)
		if err != nil {
			fmt.Errorf("%s:%s\n", err, object)
		} else if numBytes > 0 {
			for i, line := range bytes.Split(content, []byte("\n")) {
				if bytes.Contains(line, []byte(query)) {
					results <- GrepResult{
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

func getContentExcerpt(text []byte, query []byte) []byte {
	index := bytes.Index(text, query)
	from := int(math.Max(float64(index) - 10, 0))
	to := index + len(query) + 10

	return text[from:to]
}
