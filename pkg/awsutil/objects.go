package awsutil

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// ObjectExistsInBucket returns true if the key provided matches a single object in the bucket
// provided or false otherwise
func ObjectExistsInBucket(bucket string, key string, sess *session.Session) (bool, error) {
	svc := s3.New(sess)
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}

	result, err := svc.ListObjectsV2(input)
	if err != nil {
		return false, err
	}

	for _, object := range result.Contents {
		if *object.Key == key {
			return true, nil
		}
	}

	return false, nil
}
