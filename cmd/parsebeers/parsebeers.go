//TODO(whatsfordinner): comment your code, you monster
package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/whatsfordinner/beer-art-bot/pkg/awsutil"
	"github.com/whatsfordinner/beer-art-bot/pkg/brewerydb"
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

	//TODO(whatsfordinner): shoud potentially all go into one package for AWS interactions
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	beerBytes, err := awsutil.GetByteSliceFromS3(bucket, "beer_names.json", sess)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	var beers brewerydb.BeerOutput
	err = json.Unmarshal(beerBytes, &beers)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	numBeers := 20
	beersToParse := beers.BeerData[0:numBeers]
	startTime := time.Now()
	awsBeers := []*string{}
	for _, beer := range beersToParse {
		awsBeers = append(awsBeers, aws.String(strings.ToLower(beer)))
	}
	comprehendResults, err := parseBeersWithComprehend(awsBeers, sess)
	if err != nil {
		log.Fatalf("%s", err)
	}
	finishTime := time.Now()
	log.Printf("Parsing %d beers with comprehend took %v", numBeers, finishTime.Sub(startTime))
	log.Printf("%v", comprehendResults)
	log.Printf("done")
}

func parseBeersWithComprehend(beers []*string, sess *session.Session) ([]string, error) {
	log.Printf("starting to parse %d beers with comprehend", len(beers))
	svc := comprehend.New(sess)
	taggedBeers := []string{}
	result, err := svc.BatchDetectSyntax(&comprehend.BatchDetectSyntaxInput{
		LanguageCode: aws.String("en"),
		TextList:     beers,
	})

	if err != nil {
		return taggedBeers, err
	}

	for _, results := range result.ResultList {
		tags := []string{}
		for _, tokens := range results.SyntaxTokens {
			tags = append(tags, *tokens.PartOfSpeech.Tag)
		}
		taggedBeers = append(taggedBeers, strings.Join(tags, " "))
	}

	return taggedBeers, nil
}
