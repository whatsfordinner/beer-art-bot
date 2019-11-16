//TODO(whatsfordinner): comment your code, you monster
package main

import (
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/whatsfordinner/beer-art-bot/pkg/brewerydb"
)

type getBeersConfig struct {
	Bucket string
	Region string
}

func main() {
	config, err := loadGetBeersConfig()
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

	log.Printf("fetching beers and styles from BreweryDB")
	names, styles, err := brewerydb.GetBeerData()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	err = brewerydb.WriteBeerOutputToS3(config.Bucket, "beer_names.json", names, sess)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	err = brewerydb.WriteBeerOutputToS3(config.Bucket, "beer_styles.json", styles, sess)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	log.Printf("finished")
}

func loadGetBeersConfig() (getBeersConfig, error) {
	bucket, exist := os.LookupEnv("UPLOAD_BUCKET")
	if !exist {
		return getBeersConfig{}, errors.New("UPLOAD_BUCKET not set, cannot continue")
	}
	region, exist := os.LookupEnv("UPLOAD_REGION")
	if !exist {
		return getBeersConfig{}, errors.New("UPLOAD_REGION not set, cannot continue")
	}

	return getBeersConfig{
		Bucket: bucket,
		Region: region,
	}, nil
}
