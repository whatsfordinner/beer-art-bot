package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// writeOutputToDisk will marshal a BeerOutput into a JSON string and write it to the working directory.
func writeOutputToDisk(output beerOutput, filename string) error {
	outputBlob, err := json.Marshal(output)
	if err != nil {
		return err
	}
	workingDirectory, err := os.Getwd()
	if err != nil {
		return err
	}

	log.Printf("writing file: %s in directory %s", filename, workingDirectory)
	err = ioutil.WriteFile(filename, outputBlob, 0644)
	if err != nil {
		return err
	}

	return nil
}
