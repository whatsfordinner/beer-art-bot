package awsutil

import (
	"bytes"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func writeByteSliceToS3(bucket string, key string, blob []byte, sess *session.Session, region string) error {
	uploader := s3manager.NewUploader(sess)
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

func getByteSliceFromS3(bucket string, key string, sess *session.Session, region string) ([]byte, error) {
	downloader := s3manager.NewDownloader(sess)
	awsBuff := aws.NewWriteAtBuffer(make([]byte, 0))
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
