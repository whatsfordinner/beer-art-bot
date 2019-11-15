package main

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/whatsfordinner/beer-art-bot/pkg/brewerydb"
	"github.com/whatsfordinner/beer-art-bot/pkg/sliceutil"
)

func main() {
	config, err := loadParseBeersConfig()
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

	err = parseBeers(config.Bucket, sess)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}

func parseBeers(bucket string, sess *session.Session) error {
	// get the list of all beers from S3
	beers, err := brewerydb.GetBeerOutputFromS3(bucket, "beer_names.json", true, sess)
	if err != nil {
		return err
	}

	// get the list of all beers that have already been processed from S3
	previouslyProcessedBeers, err := brewerydb.GetBeerOutputFromS3(bucket, "processed_beers.json", false, sess)
	if err != nil {
		return err
	}

	// make a slice of beer names that contain only beers that haven't been processed
	beersToProcess := sliceutil.GetMutuallyExclusiveElements(beers.BeerData, previouslyProcessedBeers.BeerData)

	numBeers := 3
	beersToParse := beersToProcess[0:numBeers]
	startTime := time.Now()
	// get the beers tagged for parts of speech by comprehend
	comprehendResults, err := parseBeersWithComprehend(beersToParse, sess)
	if err != nil {
		return err
	}
	finishTime := time.Now()
	log.Printf("Parsing %d beers with comprehend took %v", numBeers, finishTime.Sub(startTime))

	// get the list of already processed tags from S3
	previouslyProcessedTags, err := brewerydb.GetBeerOutputFromS3(bucket, "beer_tags.json", false, sess)
	if err != nil {
		return err
	}

	// write all beer tags to S3
	var beerTags brewerydb.BeerOutput
	beerTags.BeerData = append(previouslyProcessedTags.BeerData, comprehendResults...)
	err = brewerydb.WriteBeerOutputToS3(bucket, "beer_tags.json", beerTags, sess)
	if err != nil {
		return err
	}

	// append the list of newly processed beers to the existing list and write to S3
	var processedBeers brewerydb.BeerOutput
	processedBeers.BeerData = append(previouslyProcessedBeers.BeerData, beersToParse...)
	err = brewerydb.WriteBeerOutputToS3(bucket, "processed_beers.json", processedBeers, sess)
	if err != nil {
		return err
	}

	return nil
}
