package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticache"
)

// ListCache gets all the elasticache instances and returns all their attributes
func (cacheSess *Client) ListCache() ([]*elasticache.CacheCluster, error) {

	svc := elasticache.New(cacheSess.session)
	c := &elasticache.DescribeCacheClustersInput{
		MaxRecords: aws.Int64(100), // get max number of records possible
	}

	resp, err := svc.DescribeCacheClusters(c)
	if err != nil {
		return nil, err
	}
	return resp.CacheClusters, err
}

// DescribeCache returns all attributes of cache cluster given cluster id
func (cacheSess *Client) DescribeCache(id *string) ([]*elasticache.CacheCluster, error) {

	svc := elasticache.New(cacheSess.session)
	c := &elasticache.DescribeCacheClustersInput{
		CacheClusterId: aws.String(*id),
	}

	resp, err := svc.DescribeCacheClusters(c)
	if err != nil {
		return nil, err
	}
	return resp.CacheClusters, err
}

//DeleteCache deltes an elasticache cluster
func (cacheSess *Client) DeleteCache(id *string) error {

	svc := elasticache.New(cacheSess.session)

	c := &elasticache.DeleteCacheClusterInput{
		CacheClusterId: aws.String(*id), // Required
	}
	resp, err := svc.DeleteCacheCluster(c)
	if err != nil {
		return err
	}
	fmt.Println(resp)
	return err
}

//ListCacheTags returns tags on Cache
func (cacheSess *Client) ListCacheTags(id *string) ([]*elasticache.Tag, error) {

	svc := elasticache.New(cacheSess.session)

	params := &elasticache.ListTagsForResourceInput{
		ResourceName: aws.String(*id), // Required
	}

	resp, err := svc.ListTagsForResource(params)
	if err != nil {
		return nil, err
	}
	return resp.TagList, err
}
