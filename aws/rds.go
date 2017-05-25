package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
)

//ListRDS shows all the rds instances in the account via their identifier
func (rdsSess *Client) ListRDS() ([]*rds.DBInstance, error) {

	svc := rds.New(rdsSess.session)
	db := &rds.DescribeDBInstancesInput{}
	resp, err := svc.DescribeDBInstances(db)
	if err != nil {
		return nil, err
	}

	// return all the attributes about each rds instance
	return resp.DBInstances, err
}

//DeleteRDS nukes a single rds instance
func (rdsSess *Client) DeleteRDS(identifier *string) error {
	svc := rds.New(rdsSess.session)

	db := &rds.DeleteDBInstanceInput{
		DBInstanceIdentifier: aws.String(*identifier),
	}
	resp, err := svc.DeleteDBInstance(db)
	if err != nil {
		return err
	}

	fmt.Println(resp.DBInstance)

	return err

}

//ListRdsTags list all the tags for a resource
func (rdsSess *Client) ListRdsTags(id *string) ([]*rds.Tag, error) {
	svc := rds.New(rdsSess.session)
	db := &rds.ListTagsForResourceInput{
		ResourceName: aws.String(*id), // ze arn
	}
	resp, err := svc.ListTagsForResource(db)
	if err != nil {
		return nil, err
	}

	// return all the attributes about each rds instance tags
	return resp.TagList, err
}
