package awsutil

import "testing"

var getByteSliceTests = map[string]struct {
	k   string
	b   string
	e   string
	err bool
}{
	"bucket doesn't exist": {"test", "foobar", "testing", true},
	"object doesn't exist": {"foobar", "bucket", "testing", true},
	"object exists":        {"test", "bucket", "testing", false},
}

func TestGetByteSliceFromS3(t *testing.T) {
	teardown := setupS3TestScaffold(t)
	defer teardown(t)
	testSession := getLocalstackSession()
	for test, tt := range getByteSliceTests {
		t.Run(test, func(t *testing.T) {
			result, err := GetByteSliceFromS3(tt.b, tt.k, testSession)

			if tt.err && (err == nil) {
				t.Errorf("expected error but got no error\n")
			}

			if !tt.err && (err != nil) {
				t.Errorf("expected no error but got %s\n", err.Error())
			}

			if !tt.err && (err == nil) {
				if tt.e != string(result) {
					t.Errorf("exptected %s but got %s\n", tt.e, string(result))
				}
			}
		})
	}
}

var writeByteSliceTests = map[string]struct {
	k   string
	b   string
	c   string
	err bool
}{
	"bucket doesn't exist":        {"qux", "foobar", "sometestext", true},
	"writing a new object":        {"qux", "bucket", "sometesttext", false},
	"overwriting existing object": {"qux", "bucket", "sometesttext", false},
}

func TestWriteByteSliceToS3(t *testing.T) {
	teardown := setupS3TestScaffold(t)
	defer teardown(t)
	testSession := getLocalstackSession()
	for test, tt := range writeByteSliceTests {
		t.Run(test, func(t *testing.T) {
			content := []byte(tt.c)
			err := WriteByteSliceToS3(tt.b, tt.k, content, testSession)

			if tt.err && (err == nil) {
				t.Errorf("expected error but got no error\n")
			}

			if !tt.err && (err != nil) {
				t.Errorf("expected no error but got %s\n", err.Error())
			}

			if !tt.err && (err == nil) {
				result, err := GetByteSliceFromS3(tt.b, tt.k, testSession)
				if err != nil {
					t.Errorf("%s\n", err.Error())
				}

				if tt.c != string(result) {
					t.Errorf("expected %s but got %s\n", tt.c, string(result))
				}
			}
		})
	}
}
