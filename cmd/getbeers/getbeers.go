//TODO(whatsfordinner): comment your code, you monster
package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/whatsfordinner/beer-art-bot/pkg/awsutil"
	"github.com/whatsfordinner/beer-art-bot/pkg/brewerydb"
)

func main() {
	log.Printf("fetching beers and styles from BreweryDB")
	names, styles, err := brewerydb.GetBeerData()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	//TOOD(whatsfordinner): DRY. separate into function or make method in brewerydb package
	log.Printf("converting %d beer names to JSON blob", len(names.BeerData))
	namesBlob, err := json.Marshal(names)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	log.Printf("converting %d beer styles to JSON blob", len(styles.BeerData))
	stylesBlob, err := json.Marshal(styles)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	//TODO(whatsfordinner): separate config into struct and own function/init()
	bucket, exist := os.LookupEnv("UPLOAD_BUCKET")
	if !exist {
		log.Fatalf("UPLOAD_BUCKET not set, cannot continue")
	}
	region, exist := os.LookupEnv("UPLOAD_REGION")
	if !exist {
		log.Fatalf("UPLOAD_REGION not set, cannot continue")
	}

	log.Printf("creating new AWS session")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	err = awsutil.WriteByteSliceToS3(bucket, "beer_names.json", namesBlob, sess)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	err = awsutil.WriteByteSliceToS3(bucket, "beer_styles.json", stylesBlob, sess)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	log.Printf("finished")
}
