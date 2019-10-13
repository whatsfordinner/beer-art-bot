package main

import (
	"log"
	"os"
	"testing"
)

var configTests = map[string]struct {
	setEndpoint bool
	setAPIKey   bool
	shouldErr   bool
}{
	"both set":          {true, true, false},
	"only endpoint set": {true, false, true},
	"only API Key set":  {false, true, true},
	"neither set":       {false, false, true},
}

func TestLoadBreweryDBConfig(t *testing.T) {
	for test, tt := range configTests {
		t.Run(test, func(t *testing.T) {
			if tt.setEndpoint {
				err := os.Setenv("BREWERYDB_ENDPOINT", "fake endpoint")
				if err != nil {
					log.Fatalf("cannot set environment variable BREWERYDB_ENDPOINT")
				}
			} else {
				err := os.Unsetenv("BREWERYDB_ENDPOINT")
				if err != nil {
					log.Fatalf("cannot set environment variable BREWERYDB_ENDPOINT")
				}
			}

			if tt.setAPIKey {
				err := os.Setenv("BREWERYDB_APIKEY", "fake API key")
				if err != nil {
					log.Fatalf("cannot set environment variable BREWERYDB_APIKEY")
				}
			} else {
				err := os.Unsetenv("BREWERYDB_APIKEY")
				if err != nil {
					log.Fatalf("cannot set environment variable BREWERYDB_APIKEY")
				}
			}

			testConfig, testErr := loadBreweryDBConfig()

			if tt.shouldErr && testErr == nil {
				t.Errorf("expected error but got no error\n")
			}

			if !tt.shouldErr && testErr != nil {
				t.Errorf("expected no error but got %s\n", testErr.Error())
			}

			if tt.setEndpoint && tt.setAPIKey && testErr == nil {
				if testConfig.Endpoint != "fake endpoint" {
					t.Errorf("mismatched endpoint. Expected %s but got %s\n", "fake endpoint", testConfig.Endpoint)
				}

				if testConfig.APIKey != "fake API key" {
					t.Errorf("mismatched API key. Exptected %s but got %s\n", "fake API key", testConfig.APIKey)
				}
			}
		})
	}
}

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
			testResult := sliceContains(tt.testSlice, tt.testElement)
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
			resultSlice := appendIfUnique(tt.inputSlice, tt.inputElement)

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
