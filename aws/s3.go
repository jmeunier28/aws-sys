package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

//ListS3 list all dem buckets in the account
func (s3Sess *Client) ListS3() (*s3.ListBucketsOutput, error) {

	svc := s3.New(s3Sess.session)
	buck := &s3.ListBucketsInput{}
	resp, err := svc.ListBuckets(buck)
	if err != nil {
		return nil, err
	}
	//return all attributes about the s3 buckets
	return resp, err
}

//DeleteS3 deletes bucket based on bucket name
func (s3Sess *Client) DeleteS3(name *string) error {
	svc := s3.New(s3Sess.session)
	buck := &s3.DeleteBucketInput{
		Bucket: aws.String(*name),
	}
	_, err := svc.DeleteBucket(buck)
	if err != nil {
		return err
	}
	return err
}

// CreateBucket creates an S3 bucket
func (s3Sess *Client) CreateBucket(name *string) error {

	svc := s3.New(s3Sess.session)

	buck := &s3.CreateBucketInput{
		Bucket: aws.String(*name),
	}
	_, err := svc.CreateBucket(buck)
	if err != nil {
		return err
	}
	return err

}
