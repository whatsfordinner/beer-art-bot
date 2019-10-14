package main

import (
	"encoding/json"
	"log"
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

type beerOutput struct {
	BeerData []string `json:"beerData"`
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

	for _, i := range parsedResponse.Data {
		beers = append(beers, beerDetails{i.NameDisplay, i.Style.ShortName})
	}

	return beers, parsedResponse.CurrentPage, parsedResponse.NumberOfPages, nil
}

func getBeerData() (beerOutput, beerOutput, error) {
	config, err := loadBreweryDBConfig()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	// names will contain the final list of unique beers
	var names beerOutput
	// styles will contain the final list of unique beer styles
	var styles beerOutput
	// results will contain just the results from any particular query
	var results []beerDetails
	var currPage int
	// maxPage is intentionally set here because we don't know how many pages the query will have until we return the first time
	maxPage := 1
	for currPage = 1; currPage <= maxPage; currPage++ {
		rawResults, err := queryBeersAPI(currPage, config.Endpoint, config.APIKey)
		if err != nil {
			return names, styles, err
		}

		// maxPage gets set with the correct number here, this also accounts for the number of pages changing
		results, currPage, maxPage, err = parseQueryBeersAPI(rawResults)
		if err != nil {
			return names, styles, err
		}

		for _, b := range results {
			names.BeerData = appendIfUnique(names.BeerData, b.Name)
			styles.BeerData = appendIfUnique(styles.BeerData, b.Style)
		}

		// BreweryDB's API has a limit of 10 requests per second, this is a dumb implementation to honour that
		// TODO: implement smarter version that takes into account time since query
		time.Sleep(120 * time.Millisecond)
	}

	return names, styles, nil
}
