package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// AWSSession saves config to to create an aws service clients
type AWSSession struct {
	Session *session.Session
}

// NewAWSParams creates a new AWSSession object
func NewAWSSession(awsProfile string) (*AWSSession, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String("us-east-2")},
		Profile:           awsProfile,
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		return nil, err
	}

	return &AWSSession{
		Session: sess,
	}, nil

}
