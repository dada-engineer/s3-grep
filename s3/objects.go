package s3

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/dabdada/s3-grep/config"
)

// Object safes Key and Content of a object in S3
type Object struct {
	Key      string
	Content  []byte
	NumBytes int64
	Error    error
}

// NewObject is a constructor for Objects
func NewObject(key string) Object {
	return Object{Key: key}
}

// ListObjects lists all objects in the specified bucket
func ListObjects(svc s3iface.S3API, bucketName string) ([]Object, error) {
	var objects []Object
	err := svc.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	}, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {
		for _, obj := range p.Contents {
			objects = append(objects, NewObject(*obj.Key))
		}
		return true
	})
	if err != nil {
		return []Object{}, err
	}

	return objects, nil
}

// GetObjectContent loads the content of a S3 object key into a buffer
func (o Object) GetObjectContent(session *config.AWSSession, bucketName string) {
	if o.Key == "" {
		o.Error = errors.New("Object has no Key")
		return
	}
	buff := &aws.WriteAtBuffer{}
	downloader := s3manager.NewDownloader(session.Session)
	numBytes, err := downloader.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(o.Key),
	})

	o.Content = buff.Bytes()
	o.NumBytes = numBytes
	o.Error = err
}
