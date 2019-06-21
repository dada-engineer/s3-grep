package s3

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/dabdada/s3-grep/config"
)

// StoredObject provides an interface for Objects in a Cloud Storage
type StoredObject interface {
	GetKey() string
	GetContent(*config.AWSSession, string) ([]byte, int64, error)
}

// Object provides an Object with a Key
type Object struct {
	Key string
}

// GetKey is a getter method to get the Key of the Object
func (o Object) GetKey() string {
	return o.Key
}

// GetContent loads the content of a S3 object key into a buffer
func (o Object) GetContent(session *config.AWSSession, bucketName string) ([]byte, int64, error) {
	if o.Key == "" {
		return []byte{}, 0, errors.New("Object has no Key")
	}
	buff := &aws.WriteAtBuffer{}
	downloader := s3manager.NewDownloader(session.Session)
	numBytes, err := downloader.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(o.Key),
	})

	return buff.Bytes(), numBytes, err
}

// NewObject is a constructor for Objects
func NewObject(key string) StoredObject {
	return Object{Key: key}
}

// ListObjects lists all objects in the specified bucket
func ListObjects(svc s3iface.S3API, bucketName string, prefix string) ([]StoredObject, error) {
	var objects []StoredObject
	prefixExpression, err := regexp.Compile(fmt.Sprintf("^%s", strings.Trim(strings.TrimSpace(prefix), "/")))
	if err != nil {
		return nil, errors.New("The provided prefix is not a valid regular expression")
	}
	err = svc.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	}, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {
		for _, obj := range p.Contents {
			if prefixExpression.MatchString(*obj.Key) {
				objects = append(objects, NewObject(*obj.Key))
			}
		}
		return true
	})
	if err != nil {
		return []StoredObject{}, err
	}

	return objects, nil
}
