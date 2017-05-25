package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// ListEC2 will return a string slice of all the instance ids in the account
func (ec2Sess *Client) ListEC2() ([]*ec2.Reservation, error) {

	svc := ec2.New(ec2Sess.session)
	serv := &ec2.DescribeInstancesInput{
		MaxResults: aws.Int64(500),
	}

	// have to make this call multiple times to get them all
	resp, err := svc.DescribeInstances(serv)
	if err != nil {
		return nil, err
	}

	//return all the attributes of the ec2 instances
	return resp.Reservations, err
}

//DeleteEC2 deletes instance based on instance id
func (ec2Sess *Client) DeleteEC2(id *string) (err error) {

	svc := ec2.New(ec2Sess.session)
	params := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(*id),
		},
	}
	resp, err := svc.TerminateInstances(params)
	if err != nil {
		return err
	}
	fmt.Println(resp.TerminatingInstances)
	return err
}
