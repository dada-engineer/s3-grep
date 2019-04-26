package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dabdada/s3-grep/config"
)

// Validator to ensure bucket is available in profile
func IsBucket(session config.AWSSession, bucketName string) bool {
	svc := s3.New(session.Session)
	headInput := &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	}
	_, err := svc.HeadBucket(headInput)

	if err != nil {
		return true
	}
	return false
}
