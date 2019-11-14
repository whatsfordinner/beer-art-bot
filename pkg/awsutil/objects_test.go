package awsutil

import (
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
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

var objectExistsTests = map[string]struct {
	o   string
	b   string
	e   bool
	err bool
}{
	"bucket doesn't exist": {"test", "foobar", false, true},
	"object doesn't exist": {"foobar", "bucket", false, false},
	"object exists":        {"test", "bucket", true, false},
}

func TestObjectExistsInBucket(t *testing.T) {
	testSession := getLocalstackSession()
	for test, tt := range objectExistsTests {
		t.Run(test, func(t *testing.T) {
			result, err := ObjectExistsInBucket(tt.b, tt.o, testSession)

			if tt.err && (err == nil) {
				t.Errorf("expected error but got no error\n")
			}

			if !tt.err && (err != nil) {
				t.Errorf("expected no error but got %s\n", err.Error())
			}

			if !tt.err && (err == nil) {
				if tt.e != result {
					t.Errorf("expected %v but got %v\n", tt.e, result)
				}
			}
		})
	}
}
