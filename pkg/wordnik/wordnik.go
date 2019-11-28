package wordnik

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Client is the struct used for retrieving random words from the Wordnik API
type Client struct {
	APIKey string
}

type randomWordReturn struct {
	ID   int    `json:"id"`
	Word string `josn:"word"`
}

type messageReturn struct {
	Message string `json:"message"`
}

// NewClient takes in your Wordnik API key and returns a Client
func NewClient(apiKey string) *Client {
	return &Client{
		APIKey: apiKey,
	}
}

// GetRandomWord gets a random word from the Wordnik API with the type specified
func (c *Client) GetRandomWord(typeOfWord string) (string, error) {
	wordObject, err := queryWordnikAPI(c.APIKey, typeOfWord)
	if err != nil {
		return "", err
	}
	return wordObject.Word, nil
}

func queryWordnikAPI(apiKey string, typeOfWord string) (randomWordReturn, error) {
	endpoint := "https://api.wordnik.com/v4/words.json/randomWord?"
	params := fmt.Sprintf("hasDictionaryDef=true&includePartOfSpeech=%s", typeOfWord)
	auth := fmt.Sprintf("api_key=%s", apiKey)
	queryString := fmt.Sprintf("%s%s&%s", endpoint, params, auth)

	response, err := http.Get(queryString)
	if err != nil {
		return randomWordReturn{}, err
	}

	defer response.Body.Close()
	wordBlob, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return randomWordReturn{}, err
	}

	if response.StatusCode != 200 {
		var returnObject messageReturn
		err = json.Unmarshal(wordBlob, &returnObject)
		if err != nil {
			return randomWordReturn{}, err
		}

		errorString := fmt.Sprintf("wordnik API error: Status Code: %d. Message: %s", response.StatusCode, returnObject.Message)
		return randomWordReturn{}, errors.New(errorString)
	}

	var wordObject randomWordReturn
	err = json.Unmarshal(wordBlob, &wordObject)
	if err != nil {
		return randomWordReturn{}, err
	}

	return wordObject, nil
}
