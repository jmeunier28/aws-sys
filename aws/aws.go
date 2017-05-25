package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

//Client is a AWS Client Session
type Client struct {
	session *session.Session
}

// NewClient creates new session
func NewClient(accessKey, secretKey, region string) *Client {
	conf := &aws.Config{Region: aws.String(region)}
	conf = conf.WithCredentials(credentials.NewStaticCredentials(accessKey, secretKey, ""))
	c := &Client{
		session: session.Must(session.NewSessionWithOptions(session.Options{
			Config: aws.Config(*conf),
		})),
	}
	return c
}
