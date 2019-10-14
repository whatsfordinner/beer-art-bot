package main

import (
	"log"
)

func main() {
	names, styles, err := getBeerData()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	err = writeOutputToDisk(names, "beer_names.json")
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	err = writeOutputToDisk(styles, "beer_styles.json")
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	log.Printf("data fetch complete")
}
