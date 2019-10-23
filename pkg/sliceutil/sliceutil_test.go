package sliceutil

import (
	"testing"
)

var containsTests = map[string]struct {
	s      []string
	e      string
	expect bool
}{
	"empty slice":                            {[]string{}, "test", false},
	"populated slice, element doesn't exist": {[]string{"foo", "bar", "baz"}, "qux", false},
	"populated slice, element does exist":    {[]string{"foo", "bar", "baz"}, "foo", true},
}

func TestSliceContains(t *testing.T) {
	for test, tt := range containsTests {
		t.Run(test, func(t *testing.T) {
			result := SliceContains(tt.s, tt.e)
			if result != tt.expect {
				t.Errorf("mismatched contains. Expected %v but got %v\n", tt.expect, result)
			}
		})
	}
}

var equalTests = map[string]struct {
	a      []string
	b      []string
	expect bool
}{
	"empty slices":        {[]string{}, []string{}, true},
	"equal slices":        {[]string{"foo", "bar", "baz"}, []string{"foo", "bar", "baz"}, true},
	"uneven length":       {[]string{"foo"}, []string{"foo", "bar"}, false},
	"equal but unordered": {[]string{"foo", "bar", "baz"}, []string{"bar", "foo", "baz"}, false},
	"equal but different": {[]string{"foo", "bar", "baz"}, []string{"foo", "qux", "baz"}, false},
}

func TestSlicesEqual(t *testing.T) {
	for test, tt := range equalTests {
		t.Run(test, func(t *testing.T) {
			if tt.expect != SlicesEqual(tt.a, tt.b) {
				t.Errorf("expected %v but got %v for slices:\n%v\n%v\n", tt.expect, !tt.expect, tt.a, tt.b)
			}
		})
	}
}

var appendTests = map[string]struct {
	s      []string
	e      string
	expect []string
}{
	"empty string":          {[]string{"foo", "bar"}, "", []string{"foo", "bar"}},
	"element exists":        {[]string{"foo", "bar"}, "foo", []string{"foo", "bar"}},
	"element doesn't exist": {[]string{"foo", "bar"}, "baz", []string{"foo", "bar", "baz"}},
}

func TestAppendIfUnique(t *testing.T) {
	for test, tt := range appendTests {
		t.Run(test, func(t *testing.T) {
			result := AppendIfUnique(tt.s, tt.e)
			if !SlicesEqual(result, tt.expect) {
				t.Errorf("result not equal to expected:\nresult:\t%v\nexpect:\t%v\n", result, tt.expect)
			}
		})
	}
}

var xorTests = map[string]struct {
	a      []string
	b      []string
	expect []string
}{
	"empty slices":       {[]string{}, []string{}, []string{}},
	"no unique elements": {[]string{"foo", "bar", "baz"}, []string{"foo", "bar", "baz"}, []string{}},
	"unique in first":    {[]string{"foo", "bar", "baz"}, []string{"foo", "bar"}, []string{"baz"}},
	"unique in second":   {[]string{"bar", "baz"}, []string{"foo", "bar", "baz"}, []string{"foo"}},
	"unique in both":     {[]string{"foo", "bar"}, []string{"bar", "baz"}, []string{"foo", "baz"}},
}

func TestGetMutuallyExclusiveElements(t *testing.T) {
	for test, tt := range xorTests {
		t.Run(test, func(t *testing.T) {
			result := GetMutuallyExclusiveElements(tt.a, tt.b)
			if !SlicesEqual(result, tt.expect) {
				t.Errorf("result not equal to expected:\nresult:\t%v\nexpect:\t%v\n", result, tt.expect)
			}
		})
	}
}
