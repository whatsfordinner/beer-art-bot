package main

import (
	"errors"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/whatsfordinner/beer-art-bot/pkg/brewerydb"
)

type tweetBeerConfig struct {
	Bucket        string
	Region        string
	WordnikAPIKey string
}

func main() {
	config, err := loadTweetBeerConfig()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	log.Printf("creating new AWS session")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.Region),
	})
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	beerTags, err := brewerydb.GetBeerOutputFromS3(config.Bucket, "beer_tags.json", true, sess)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	log.Printf("retrieved %d sets of parts of speech", len(beerTags.BeerData))

	beerStyles, err := brewerydb.GetBeerOutputFromS3(config.Bucket, "beer_styles.json", true, sess)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	log.Printf("retrieved %d beer styles", len(beerStyles.BeerData))

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	selectedBeerTag := beerTags.BeerData[r.Intn(len(beerTags.BeerData))]
	log.Printf("selected %s for parts of speech", selectedBeerTag)
	selectedBeerStyle := beerStyles.BeerData[r.Intn(len(beerStyles.BeerData))]
	log.Printf("selected %s for beer style", selectedBeerStyle)

}

func loadTweetBeerConfig() (tweetBeerConfig, error) {
	bucket, exist := os.LookupEnv("CORPUS_BUCKET")
	if !exist {
		return tweetBeerConfig{}, errors.New("CORPUS_BUCKET not set, cannot continue")
	}

	region, exist := os.LookupEnv("CORPUS_REGION")
	if !exist {
		return tweetBeerConfig{}, errors.New("CORPUS_REGION not set, cannot continue")
	}

	wordnikKey, exist := os.LookupEnv("WORDNIK_APIKEY")
	if !exist {
		return tweetBeerConfig{}, errors.New("WORDNIK_APIKEY not set, cannot continue")
	}

	return tweetBeerConfig{
		Bucket:        bucket,
		Region:        region,
		WordnikAPIKey: wordnikKey,
	}, nil
}
