package brewerydb

import (
	"errors"
	"log"
	"os"
)

// breweryDBConfig contains configuration details for accessing the BreweryDB API.
type breweryDBConfig struct {
	Endpoint string
	APIKey   string
}

// loadBreweryDBConfig fetches config values and returns object.
func loadBreweryDBConfig() (breweryDBConfig, error) {
	log.Printf("fetching environment variables")
	var config breweryDBConfig

	// BREWERYDB_ENDPOINT is not necessarily a static thing: BreweryDB offers a sandbox API as well as a prod one
	endpoint, exists := os.LookupEnv("BREWERYDB_ENDPOINT")
	if !exists {
		err := errors.New("environment variable BREWERYDB_ENDPOINT not found")
		return config, err
	}

	apiKey, exists := os.LookupEnv("BREWERYDB_APIKEY")
	if !exists {
		err := errors.New("environment variable BREWERYDB_APIKEY not found")
		return config, err
	}

	config.Endpoint = endpoint
	config.APIKey = apiKey

	return config, nil
}
