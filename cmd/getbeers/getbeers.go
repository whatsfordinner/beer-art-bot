//TODO(whatsfordinner): comment your code, you monster
package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

	//TODO(whatsfordinner): should potentially all go into one package for AWS interactions
	log.Printf("creating new AWS session")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	s3Uploader := s3manager.NewUploader(sess)
	err = writeByteSliceToS3(bucket, "beer_names.json", namesBlob, s3Uploader)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	err = writeByteSliceToS3(bucket, "beer_styles.json", stylesBlob, s3Uploader)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	log.Printf("finished")
}

func writeByteSliceToS3(bucket string, key string, blob []byte, uploader *s3manager.Uploader) error {
	log.Printf("uploading %s to %s", key, bucket)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(blob),
	})
	if err != nil {
		return err
	}
	log.Printf("successfully uploaded %s", result.Location)
	return nil
}
