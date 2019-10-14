package brewerydb

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
