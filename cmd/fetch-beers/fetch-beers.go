package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// queryBeersReturn matches the structure of the JSON returned by the BreweryDB /beers endpoint.
// It intentionally omits most of the returned fields to ensure only those fields we care about
// are retained
type queryBeersReturn struct {
	CurrentPage   int `json:"currentPage"`
	NumberOfPages int `json:"numberOfPages"`
	Data          []struct {
		NameDisplay string `json:"nameDisplay"`
		Style       struct {
			ShortName string `json:"shortName"`
		} `json:"style"`
	} `json:"data"`
}

// beerDetails is a flat struct of only the information we care about for creating the corpus.
type beerDetails struct {
	Name  string `json:"name"`
	Style string `json:"style"`
}

// breweryDBConfig contains configuration details for accessing the BreweryDB API.
type breweryDBConfig struct {
	Endpoint string
	APIKey   string
}

// loadBreweryDBConfig fetches config values and returns object.
func loadBreweryDBConfig() (breweryDBConfig, error) {
	log.Printf("fetching environment variables")
	var config breweryDBConfig

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

// queryBeersAPI is a basic wrapper around the BreweryDB API.
// It queries the /beers endpoint and returns the body of the HTTP request.
func queryBeersAPI(pageNumber int, endpoint string, apiKey string) ([]byte, error) {
	queryString := fmt.Sprintf("%s/beers?key=%s&p=%d", endpoint, apiKey, pageNumber)
	log.Printf("sending query: %s", queryString)

	response, err := http.Get(queryString)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

// parseQueryBeersAPI parses the byte slice returned by queryBeersAPI
// It returns the relevant details of each beer, the page those results come from
// and the maximum page
func parseQueryBeersAPI(data []byte) ([]beerDetails, int, int, error) {
	var parsedResponse queryBeersReturn
	var beers []beerDetails
	if err := json.Unmarshal(data, &parsedResponse); err != nil {
		// In all cases where we return an error we have to assume that any return value
		// other than the error is unusable
		return beers, 0, 0, err
	}

	for i := range parsedResponse.Data {
		beers = append(beers, beerDetails{parsedResponse.Data[i].NameDisplay, parsedResponse.Data[i].Style.ShortName})
	}

	return beers, parsedResponse.CurrentPage, parsedResponse.NumberOfPages, nil
}

func main() {
	config, err := loadBreweryDBConfig()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	// beers will contain the final, concatenated results
	var beers []beerDetails
	// results will contain just the results from any particular query
	var results []beerDetails
	var currPage int
	// maxPage is intentionally set here because we don't know how many pages the query will have until we return
	maxPage := 1
	for currPage = 1; currPage <= maxPage; currPage++ {
		rawResults, err := queryBeersAPI(currPage, config.Endpoint, config.APIKey)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		// maxPage gets set with the correct number here, this also accounts for the number of pages changing
		results, currPage, maxPage, err = parseQueryBeersAPI(rawResults)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		beers = append(beers, results...)

		// BreweryDB's API has a limit of 10 requests per second, this is a dumb implementation to honour that
		// TODO: implement smarter version that takes into account time since query
		time.Sleep(120 * time.Millisecond)
	}

	log.Printf("Received results:\n%+v", beers)
}
