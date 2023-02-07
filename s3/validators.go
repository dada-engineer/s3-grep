package s3

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dabdada/s3-grep/config"
)

// Get bucket region: https://github.com/aws/aws-sdk-go/issues/720
func GetBucketRegion(bucketName string) (string, error) {
	bucketUri := fmt.Sprintf("https://%s.s3.amazonaws.com", bucketName)
	res, err := http.Head(bucketUri)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	header := "X-Amz-Bucket-Region"

	if len(res.Header[header]) == 0 {
		return "", fmt.Errorf("header '%s' not found in response for bucket: %s", header, bucketUri)
	}

	return res.Header.Get(header), nil
}

// IsBucket validator to ensure bucket is available in profile
func IsBucket(session config.AWSSession, bucketName string) (bool, error) {
	svc := s3.New(session.Session)
	headInput := &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	}
	_, err := svc.HeadBucket(headInput)

	if err != nil {
		return false, err
	}
	return true, err
}
