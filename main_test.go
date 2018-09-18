package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	t.Run("boom", func(t *testing.T) {
		inputJson := readJsonFromFile(t, "./testdata/s3-event.json")

		var s3Event events.S3Event
		if err := json.Unmarshal(inputJson, &s3Event); err != nil {
			t.Errorf("could not unmarshal event: %v", err)
		}

		_, err := handler(s3Event)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func readJsonFromFile(t *testing.T, inputFile string) []byte {
	inputJson, err := ioutil.ReadFile(inputFile)
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	return inputJson
}
