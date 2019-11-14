package awsutil

import (
	"log"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func localstackResolver(service string, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
	if service == endpoints.S3ServiceID {
		return endpoints.ResolvedEndpoint{
			URL: "http://localhost:4572",
		}, nil
	}

	return endpoints.DefaultResolver().EndpointFor(service, region, optFns...)
}

func getLocalstackSession() *session.Session {
	localstackSession, err := session.NewSession(&aws.Config{
		Region:           aws.String("us-east-1"),
		EndpointResolver: endpoints.ResolverFunc(localstackResolver),
		S3ForcePathStyle: aws.Bool(true),
	})

	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	return localstackSession
}

func setupS3TestScaffold(t *testing.T) func(t *testing.T) {
	t.Log("setting up S3 scaffold for testing")
	sess := getLocalstackSession()
	bucket := aws.String("bucket")
	testObject := aws.String("test")
	svc := s3.New(sess)

	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: bucket,
	})
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	_, err = svc.PutObject(&s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(strings.NewReader("testing")),
		Bucket: bucket,
		Key:    testObject,
	})
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	return destroyS3TestScaffold
}

func destroyS3TestScaffold(t *testing.T) {
	t.Log("destroying S3 scaffold")
	sess := getLocalstackSession()
	bucket := aws.String("bucket")
	svc := s3.New(sess)

	result, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: bucket,
	})
	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	for _, object := range result.Contents {
		_, err = svc.DeleteObject(&s3.DeleteObjectInput{
			Bucket: bucket,
			Key:    object.Key,
		})
		if err != nil {
			t.Fatalf("%s\n", err.Error())
		}
	}

	_, err = svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: bucket,
	})
	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	t.Log("S3 scaffold destroyed")
}
