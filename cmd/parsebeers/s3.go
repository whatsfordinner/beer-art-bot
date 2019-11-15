package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/whatsfordinner/beer-art-bot/pkg/awsutil"
	"github.com/whatsfordinner/beer-art-bot/pkg/brewerydb"
)

func getBeerOutputFromS3(bucket string, key string, mustExist bool, sess *session.Session) (brewerydb.BeerOutput, error) {
	var newBeerOutput brewerydb.BeerOutput
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

func writeBeerOutputToS3(bucket string, key string, data brewerydb.BeerOutput, sess *session.Session) error {
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
