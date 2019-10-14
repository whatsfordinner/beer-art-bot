package sliceutil

import (
	"testing"
)

var containsTests = map[string]struct {
	testSlice        []string
	testElement      string
	shouldReturnTrue bool
}{
	"empty slice":                            {[]string{}, "test", false},
	"populated slice, element doesn't exist": {[]string{"foo", "bar", "baz"}, "qux", false},
	"populated slice, element does exist":    {[]string{"foo", "bar", "baz"}, "foo", true},
}

func TestSliceContains(t *testing.T) {
	for test, tt := range containsTests {
		t.Run(test, func(t *testing.T) {
			testResult := SliceContains(tt.testSlice, tt.testElement)
			if testResult != tt.shouldReturnTrue {
				t.Errorf("mismatched contains. Expected %v but got %v\n", tt.shouldReturnTrue, testResult)
			}
		})
	}
}

var appendTests = map[string]struct {
	inputSlice          []string
	inputElement        string
	expectedOutputSlice []string
}{
	"empty string":          {[]string{"foo", "bar"}, "", []string{"foo", "bar"}},
	"element exists":        {[]string{"foo", "bar"}, "foo", []string{"foo", "bar"}},
	"element doesn't exist": {[]string{"foo", "bar"}, "baz", []string{"foo", "bar", "baz"}},
}

func TestAppendIfUnique(t *testing.T) {
	for test, tt := range appendTests {
		t.Run(test, func(t *testing.T) {
			resultSlice := AppendIfUnique(tt.inputSlice, tt.inputElement)

			if len(resultSlice) != len(tt.expectedOutputSlice) {
				t.Errorf("length mismatch. Expected %d elements but got %d elements\n", len(tt.expectedOutputSlice), len(resultSlice))
			}

			for i := range resultSlice {
				if resultSlice[i] != tt.expectedOutputSlice[i] {
					t.Errorf("element mismatch. Expected %s but got %s\n", resultSlice[i], tt.expectedOutputSlice[i])
				}
			}
		})
	}
}
