package awsutil

import (
	"testing"
)

var objectExistsTests = map[string]struct {
	k   string
	b   string
	e   bool
	err bool
}{
	"bucket doesn't exist": {"test", "foobar", false, true},
	"object doesn't exist": {"foobar", "bucket", false, false},
	"object exists":        {"test", "bucket", true, false},
}

func TestObjectExistsInBucket(t *testing.T) {
	teardown := setupS3TestScaffold(t)
	defer teardown(t)
	testSession := getLocalstackSession()
	for test, tt := range objectExistsTests {
		t.Run(test, func(t *testing.T) {
			result, err := ObjectExistsInBucket(tt.b, tt.k, testSession)

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
