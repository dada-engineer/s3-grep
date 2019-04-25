package config

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

// AWSSession saves config to to create an aws service clients
type AWSSession struct {
	Session *session.Session
}

// NewAWSParams creates a new AWSSession object
func NewAWSSession(awsProfile string) (*AWSSession, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile:           awsProfile,
		SharedConfigState: session.SharedConfigEnable,
	}))

	return &AWSSession{
		Session: sess,
	}, nil
}
