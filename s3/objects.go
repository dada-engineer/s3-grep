package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// list all objects in the specified bucket
func listObjects(svc s3iface.S3API, bucketName string) ([]string, error) {
	var objects []string
	err := svc.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	}, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {
		for _, obj := range p.Contents {
			objects = append(objects, *obj.Key)
		}
		return true
	})
	if err != nil {
		return []string{}, err
	}

	return objects, nil
}