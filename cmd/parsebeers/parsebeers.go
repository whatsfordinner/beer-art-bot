//TODO(whatsfordinner): comment your code, you monster
package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	prose "gopkg.in/jdkato/prose.v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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

	log.Printf("starting to parse beers")
	for _, beer := range beers.BeerData {
		//TODO(whatsfordinner): this needs to be handled concurrently, way too slow
		processedBeer, err := prose.NewDocument(strings.ToLower(beer))
		if err != nil {
			log.Fatalf(err.Error())
		}

		// join the tokens to form a string
		tags := []string{}
		for _, token := range processedBeer.Tokens() {
			tags = append(tags, token.Tag)
		}

		tagString := strings.Join(tags, " ")
		log.Printf("%s: %s", beer, tagString)
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
