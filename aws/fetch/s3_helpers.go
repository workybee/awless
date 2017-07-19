package awsfetch

import (
	"context"
	"sync"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/wallix/awless/aws/conv"
	"github.com/wallix/awless/cloud/rdf"
	"github.com/wallix/awless/fetch"
	"github.com/wallix/awless/graph"
)

func forEachBucketParallel(ctx context.Context, cache fetch.Cache, api *s3.S3, f func(b *s3.Bucket) error) error {
	var buckets []*s3.Bucket
	if cached, ok := cache.Get("getBucketsPerRegion").([]*s3.Bucket); ok && cached != nil {
		buckets = cached
	} else {
		res, err := getBucketsPerRegion(ctx, api)
		if err != nil {
			return err
		}
		buckets = res
		cache.Store("getBucketsPerRegion", res)
	}

	errc := make(chan error)
	var wg sync.WaitGroup

	for _, output := range buckets {
		wg.Add(1)
		go func(b *s3.Bucket) {
			defer wg.Done()
			if err := f(b); err != nil {
				errc <- err
			}
		}(output)
	}
	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			return err
		}
	}

	return nil
}

func fetchObjectsForBucket(ctx context.Context, api *s3.S3, bucket *s3.Bucket, resources *[]*graph.Resource) error {
	out, err := api.ListObjects(&s3.ListObjectsInput{Bucket: bucket.Name})
	if err != nil {
		return err
	}

	for _, output := range out.Contents {
		res, err := awsconv.NewResource(output)
		if err != nil {
			return err
		}
		res.Properties["Bucket"] = awssdk.StringValue(bucket.Name)
		*resources = append(*resources, res)
		parent, err := awsconv.InitResource(bucket)
		if err != nil {
			return err
		}
		parent.Relations[rdf.ParentOf] = append(parent.Relations[rdf.ParentOf], res)
		*resources = append(*resources, parent)
	}

	return nil
}

func getBucketsPerRegion(ctx context.Context, api *s3.S3) ([]*s3.Bucket, error) {
	var buckets []*s3.Bucket
	out, err := api.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return buckets, err
	}

	bucketc := make(chan *s3.Bucket)
	errc := make(chan error)

	var wg sync.WaitGroup

	for _, bucket := range out.Buckets {
		wg.Add(1)
		go func(b *s3.Bucket) {
			defer wg.Done()
			loc, err := api.GetBucketLocation(&s3.GetBucketLocationInput{Bucket: b.Name})
			if err != nil {
				errc <- err
				return
			}

			region, _ := ctx.Value("region").(string)
			switch awssdk.StringValue(loc.LocationConstraint) {
			case "":
				if region == "us-east-1" {
					bucketc <- b
				}
			case region:
				bucketc <- b
			}
		}(bucket)
	}
	go func() {
		wg.Wait()
		close(bucketc)
	}()

	for {
		select {
		case err := <-errc:
			if err != nil {
				return buckets, err
			}
		case b, ok := <-bucketc:
			if !ok {
				return buckets, nil
			}
			buckets = append(buckets, b)
		}
	}
}
