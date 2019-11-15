package brewerydb

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/whatsfordinner/beer-art-bot/pkg/awsutil"
)

var testBeerOutput = BeerOutput{
	BeerData: []string{"foo", "bar", "baz"},
}

var emptyBeerOutput = BeerOutput{
	BeerData: []string{},
}

func beerOutputToByteSlice(bo BeerOutput, t *testing.T) []byte {
	blob, err := json.Marshal(bo)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	return blob
}

func byteSliceToBeerOutput(blob []byte, t *testing.T) BeerOutput {
	var out BeerOutput
	err := json.Unmarshal(blob, &out)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	return out
}

func compareBeerOutput(a, b BeerOutput) bool {
	if len(a.BeerData) != len(b.BeerData) {
		return false
	}

	for i := range a.BeerData {
		if a.BeerData[i] != b.BeerData[i] {
			return false
		}
	}

	return true
}

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
		S3ForcePathStyle: aws.Bool(true), //localstack only supports path-style S3 names
	})

	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	return localstackSession
}

func setupS3TestScaffold(t *testing.T) func(t *testing.T) {
	t.Log("setting up S3 scaffold for testing")
	sess := getLocalstackSession()
	testContent := beerOutputToByteSlice(testBeerOutput, t)
	svc := s3.New(sess)

	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String("bucket"),
	})
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	err = awsutil.WriteByteSliceToS3("bucket", "test.json", testContent, sess)
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

var getBeerOutputTests = map[string]struct {
	b   string
	k   string
	m   bool
	e   BeerOutput
	err bool
}{
	"bucket doesn't exist":               {"foobar", "bazqux", false, testBeerOutput, true},
	"mandatory object doesn't exist":     {"bucket", "foobar", true, testBeerOutput, true},
	"non-mandatory object doesn't exist": {"bucket", "foobar", false, emptyBeerOutput, false},
	"object exists":                      {"bucket", "test.json", false, testBeerOutput, false},
}

func TestGetBeerOutputFromS3(t *testing.T) {
	teardown := setupS3TestScaffold(t)
	defer teardown(t)
	testSession := getLocalstackSession()
	for test, tt := range getBeerOutputTests {
		t.Run(test, func(t *testing.T) {
			result, err := GetBeerOutputFromS3(tt.b, tt.k, tt.m, testSession)

			if tt.err && (err == nil) {
				t.Errorf("expected error but got no error\n")
			}

			if !tt.err && (err != nil) {
				t.Errorf("expected no error but got %s\n", err.Error())
			}

			if !tt.err && (err == nil) {
				if !compareBeerOutput(tt.e, result) {
					t.Errorf("BeerOutput mismatch.\nexpected:\n%+v\ngot:\n%+v", tt.e, result)
				}
			}
		})
	}
}

var writeBeerOutputTests = map[string]struct {
	b   string
	k   string
	d   BeerOutput
	err bool
}{
	"bucket doesn't exist":        {"foobar", "bazqux", testBeerOutput, true},
	"writing new object":          {"foobar", "new_test.json", testBeerOutput, true},
	"overwriting existing object": {"foobar", "test.json", testBeerOutput, true},
}

func TestWriteBeerOutputToS3(t *testing.T) {
	teardown := setupS3TestScaffold(t)
	defer teardown(t)
	testSession := getLocalstackSession()
	for test, tt := range writeBeerOutputTests {
		t.Run(test, func(t *testing.T) {
			err := WriteBeerOutputToS3(tt.b, tt.k, tt.d, testSession)

			if tt.err && (err == nil) {
				t.Errorf("expected error but go no error\n")
			}

			if !tt.err && (err != nil) {
				t.Errorf("expected no error but got %s\n", err.Error())
			}

			if !tt.err && (err == nil) {
				result, err := GetBeerOutputFromS3(tt.b, tt.k, true, testSession)
				if err != nil {
					t.Errorf("%s\n", err.Error())
				} else if !compareBeerOutput(tt.d, result) {
					t.Errorf("BeerOutput mismatch.\nexpected:\n%+v\ngot:\n%+v", tt.d, result)
				}
			}
		})
	}
}
