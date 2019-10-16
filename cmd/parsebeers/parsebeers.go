package main

import (
	"encoding/json"
	"log"
	"os"

	prose "gopkg.in/jdkato/prose.v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/whatsfordinner/beer-art-bot/pkg/brewerydb"
)

func main() {
	log.Printf("creating new AWS session")
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

	log.Printf("starting to parse beers")
	for _, beer := range beers.BeerData {
		processedBeer, err := prose.NewDocument(beer)
		if err != nil {
			log.Fatalf(err.Error())
		}

		for _, token := range processedBeer.Tokens() {
			log.Printf("%s %s", token.Text, token.Tag)
		}
	}

	log.Printf("done")
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
