package main

import (
	"log"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/whatsfordinner/beer-art-bot/pkg/sliceutil"
)

func parseBeersWithComprehend(beers []string, sess *session.Session) ([]string, error) {
	log.Printf("starting to parse %d beers with comprehend", len(beers))
	// this is enforced by the SDK by returning an error
	// https://docs.aws.amazon.com/comprehend/latest/dg/API_BatchDetectSyntax.html
	batchSize := 25
	svc := comprehend.New(sess)
	taggedBeers := []string{}
	beersToTag, beersRemaining := sliceutil.SplitSliceAt(beers, batchSize)

	// we'll be constnatly splitting beersRemaining and once beersToTag is empty then there'll be
	// nothing left in it
	for len(beersToTag) > 0 {
		// converting the beer strings into string pointers for use with the AWS SDK
		awsBeers := []*string{}
		for _, beer := range beersToTag {
			awsBeers = append(awsBeers, aws.String(strings.ToLower(beer)))
		}

		// submitting the batch of beers to be tagged by comprehend
		result, err := svc.BatchDetectSyntax(&comprehend.BatchDetectSyntaxInput{
			LanguageCode: aws.String("en"),
			TextList:     awsBeers,
		})
		if err != nil {
			return taggedBeers, err
		}

		// the actual part of speech tags are nested quite deeply so we extract them
		// and then pull them into a single string
		for _, results := range result.ResultList {
			tags := []string{}
			for _, tokens := range results.SyntaxTokens {
				tags = append(tags, *tokens.PartOfSpeech.Tag)
			}
			taggedBeers = append(taggedBeers, strings.Join(tags, " "))
		}
		beersToTag, beersRemaining = sliceutil.SplitSliceAt(beersRemaining, batchSize)
	}

	// sorting isn't strictly necessary but makes debugging and later processing easier
	sort.Strings(taggedBeers)

	return taggedBeers, nil
}
