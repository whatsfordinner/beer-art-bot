//TODO(whatsfordinner): comment your code, you monster
package main

import (
	"encoding/json"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/whatsfordinner/beer-art-bot/pkg/awsutil"
	"github.com/whatsfordinner/beer-art-bot/pkg/brewerydb"
	"github.com/whatsfordinner/beer-art-bot/pkg/sliceutil"
)

func main() {
	log.Printf("creating new AWS session")
	//TODO(whatsfordinner): separate config into struct and own function/init
	bucket, exist := os.LookupEnv("CORPUS_BUCKET")
	if !exist {
		log.Fatalf("CORPUS_BUCKET not set, cannot continue")
	}

	region, exist := os.LookupEnv("CORPUS_REGION")
	if !exist {
		log.Fatalf("CORPUS_REGION not set, cannot continue")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	// get the list of all beers from S3
	beerBytes, err := awsutil.GetByteSliceFromS3(bucket, "beer_names.json", sess)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	var beers brewerydb.BeerOutput
	err = json.Unmarshal(beerBytes, &beers)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	// get the list of all beers that have already been processed from S3
	// if the list doesn't exist, assume that none have been processed yet

	// make a slice of beer names that contain only beers that haven't been processed
	numBeers := 40
	beersToParse := beers.BeerData[0:numBeers]
	//TODO(whatsfordinner): comprehend costs money so we should only be comprehending beers that
	// haven't been read yet
	startTime := time.Now()
	// get the beers tagged for parts of speech by comprehend
	comprehendResults, err := parseBeersWithComprehend(beersToParse, sess)
	if err != nil {
		log.Fatalf("%s", err)
	}
	finishTime := time.Now()
	log.Printf("Parsing %d beers with comprehend took %v", numBeers, finishTime.Sub(startTime))

	log.Printf("sorted syntax:")
	for _, beer := range comprehendResults {
		log.Printf("\t%s", beer)
	}

	// get the existing parts of speech from S3 and append the new ones to it

	// write the new parts of speech to S3
	var beerTags brewerydb.BeerOutput
	beerTags.BeerData = comprehendResults
	tagsBlob, err := json.Marshal(beerTags)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	err = awsutil.WriteByteSliceToS3(bucket, "beer_tags.json", tagsBlob, sess)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	// append the list of newly processed beers to the existing list and write to S3
	var processedBeers brewerydb.BeerOutput
	processedBeers.BeerData = beersToParse
	processedBlob, err := json.Marshal(processedBeers)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	err = awsutil.WriteByteSliceToS3(bucket, "processed_beers.json", processedBlob, sess)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	log.Printf("done")
}

func parseBeersWithComprehend(beers []string, sess *session.Session) ([]string, error) {
	log.Printf("starting to parse %d beers with comprehend", len(beers))
	batchSize := 25
	svc := comprehend.New(sess)
	taggedBeers := []string{}
	beersToTag, beersRemaining := sliceutil.SplitSliceAt(beers, batchSize)
	for len(beersToTag) > 0 {
		// converting the beer strings into string pointers for use with the AWS SDK
		awsBeers := []*string{}
		for _, beer := range beersToTag {
			awsBeers = append(awsBeers, aws.String(strings.ToLower(beer)))
		}

		// submitting the batch of beers to be tagged by comprehend
		result, err := svc.BatchDetectSyntax(&comprehend.BatchDetectSyntaxInput{
			LanguageCode: aws.String("en"),
			TextList:     awsBeers,
		})
		if err != nil {
			return taggedBeers, err
		}

		// the actual part of speech tags are nested quite deeply so we extract them
		// and then pull them into a single string
		for _, results := range result.ResultList {
			tags := []string{}
			for _, tokens := range results.SyntaxTokens {
				tags = append(tags, *tokens.PartOfSpeech.Tag)
			}
			taggedBeers = append(taggedBeers, strings.Join(tags, " "))
		}
		beersToTag, beersRemaining = sliceutil.SplitSliceAt(beersRemaining, batchSize)
	}
	sort.Strings(taggedBeers)

	return taggedBeers, nil
}
