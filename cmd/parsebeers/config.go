package main

import (
	"errors"
	"os"
)

type parseBeersConfig struct {
	Bucket string
	Region string
}

func loadParseBeersConfig() (parseBeersConfig, error) {
	var newConfig parseBeersConfig
	bucket, exist := os.LookupEnv("CORPUS_BUCKET")
	if !exist {
		return newConfig, errors.New("CORPUS_BUCKET not set, cannot continue")
	}

	region, exist := os.LookupEnv("CORPUS_REGION")
	if !exist {
		return newConfig, errors.New("CORPUS_REGION not set, cannot continue")
	}

	newConfig.Bucket = bucket
	newConfig.Region = region

	return newConfig, nil
}
