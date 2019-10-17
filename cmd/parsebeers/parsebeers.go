//TODO(whatsfordinner): comment your code, you monster
package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	prose "gopkg.in/jdkato/prose.v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

	s3Downloader := s3manager.NewDownloader(sess)
	beerBytes, err := getByteSliceFromS3(bucket, "beer_names.json", s3Downloader)
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
	proseResults, err := parseBeersWithProse(beersToParse)
	if err != nil {
		log.Fatalf("%s", err)
	}
	finishTime := time.Now()
	log.Printf("Parsing %d beers with prose took %v", numBeers, finishTime.Sub(startTime))

	awsBeers := []*string{}
	for _, beer := range beersToParse {
		awsBeers = append(awsBeers, aws.String(strings.ToLower(beer)))
	}

	startTime = time.Now()
	comprehendResults, err := parseBeersWithComprehend(awsBeers, sess)
	if err != nil {
		log.Fatalf("%s", err)
	}
	finishTime = time.Now()
	log.Printf("Parsing %d beers with comprehend took %v", numBeers, finishTime.Sub(startTime))
	for i, beer := range proseResults {
		log.Printf("%s:\n\tprose:\t\t%s\n\tcomprehend:\t%s", beers.BeerData[i], beer, comprehendResults[i])
	}

	log.Printf("done")
}

func parseBeersWithProse(beers []string) ([]string, error) {
	log.Printf("starting to parse %d beers with prose", len(beers))
	taggedBeers := []string{}
	for _, beer := range beers {
		processedBeer, err := prose.NewDocument(beer)
		if err != nil {
			return taggedBeers, err
		}

		// join the tokens to form a string
		tags := []string{}
		for _, token := range processedBeer.Tokens() {
			tags = append(tags, token.Tag)
		}

		taggedBeers = append(taggedBeers, strings.Join(tags, " "))
	}

	return taggedBeers, nil
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

func getByteSliceFromS3(bucket string, key string, downloader *s3manager.Downloader) ([]byte, error) {
	log.Printf("downloading %s from %s", key, bucket)
	var buf []byte
	awsBuff := aws.NewWriteAtBuffer(buf)
	_, err := downloader.Download(awsBuff, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	returnBytes := awsBuff.Bytes()
	log.Printf("successfully downloaded %d bytes", len(returnBytes))
	return returnBytes, nil
}
