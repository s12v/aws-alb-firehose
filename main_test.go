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

		err := handler(s3Event)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestCreateLogEntryFromLine(t *testing.T) {
	logEntry, err := createLogEntryFromLine(`http 2018-09-18T21:38:37.519183Z app/test1/7f050ffab5373730 95.90.211.80:4254 172.31.7.183:80 0.001 0.000 0.000 200 200 461 654 "GET http://test1-332132803.eu-west-1.elb.amazonaws.com:80/ HTTP/1.1" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36" - - arn:aws:elasticloadbalancing:eu-west-1:328305143014:targetgroup/target-group-1/76784dcf3348e51f "Root=1-5ba1705d-5617cf32286ecfb8172f73cc" "-" "-" 0 2018-09-18T21:38:37.518000Z "forward" "-" `)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if logEntry.Type != "http" {
		t.Errorf("invalid type")
	}
	if logEntry.Timestamp != "2018-09-18T21:38:37.519183Z" {
		t.Errorf("invalid timestamp")
	}
	if logEntry.Client != "95.90.211.80:4254" {
		t.Errorf("invalid client")
	}
	if logEntry.Target != "172.31.7.183:80" {
		t.Errorf("invalid target")
	}
	if logEntry.Request != "GET http://test1-332132803.eu-west-1.elb.amazonaws.com:80/ HTTP/1.1" {
		t.Errorf("invalid request, %v", logEntry.Request)
	}
}

func readJsonFromFile(t *testing.T, inputFile string) []byte {
	inputJson, err := ioutil.ReadFile(inputFile)
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	return inputJson
}
