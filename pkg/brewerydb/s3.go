package brewerydb

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/whatsfordinner/beer-art-bot/pkg/awsutil"
)

// GetBeerOutputFromS3 fetches a byte slice from S3 and unmarshals into a BeerObject. If mustExist is true
// then this function will return an error if the object with the given key doesn't exist, otherwise it returns
// an empty BeerOutput
func GetBeerOutputFromS3(bucket string, key string, mustExist bool, sess *session.Session) (BeerOutput, error) {
	var newBeerOutput BeerOutput
	// check object exists in bucket
	exists, err := awsutil.ObjectExistsInBucket(bucket, key, sess)
	if err != nil {
		return newBeerOutput, err
	}

	if mustExist && !exists {
		errorStr := fmt.Sprintf("ObjectNotFound: object with key: %s not found in bucket: %s", key, bucket)
		return newBeerOutput, errors.New(errorStr)
	}

	if !exists {
		newBeerOutput.BeerData = []string{}
	} else {
		newByteSlice, err := awsutil.GetByteSliceFromS3(bucket, key, sess)
		if err != nil {
			return newBeerOutput, err
		}

		err = json.Unmarshal(newByteSlice, &newBeerOutput)
		if err != nil {
			return newBeerOutput, err
		}
	}

	return newBeerOutput, nil
}

// WriteBeerOutputToS3 marshals a BeerOutput object into a byteslice and writes it to S3
func WriteBeerOutputToS3(bucket string, key string, data BeerOutput, sess *session.Session) error {
	blob, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = awsutil.WriteByteSliceToS3(bucket, key, blob, sess)
	if err != nil {
		return err
	}

	return nil
}
