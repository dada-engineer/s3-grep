package s3

import (
	"fmt"
	"runtime"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/dabdada/s3-grep/config"
	thisS3 "github.com/dabdada/s3-grep/s3"
)

// Grep in objects in a thisS3 bucket
func Grep(session *config.AWSSession, bucketName string, query string) {
	svc := s3.New(session.Session)

	objects, err := thisS3.listObjects(svc, bucketName)

	if err != nil {
		fmt.Errorf("%s\n", err)
		return
	}

	for _, chunk := range divideObjects(objects) {
		go
	}
}

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

func grepInObjectContent(svc s3iface.S3API, objects []string, query string) {
	for _, object := range objects {
		content := getObjectContent(svc, object)
	}
}
