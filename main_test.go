package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	t.Run("should parse and send an entry", func(t *testing.T) {
		firehoseMock := NewMockFirehoseWriter()
		firehoseWriter = firehoseMock
		s3Reader = NewMockS3Reader(gz(`http 2018-09-18T21:38:37.519183Z app/test1/7f050ffab5373730 95.90.211.80:4254 172.31.7.183:80 0.001 0.000 0.000 200 200 461 654 "GET http://test1-332132803.eu-west-1.elb.amazonaws.com:80/ HTTP/1.1" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36" - - arn:aws:elasticloadbalancing:eu-west-1:328305143014:targetgroup/target-group-1/76784dcf3348e51f "Root=1-5ba1705d-5617cf32286ecfb8172f73cc" "-" "-" 0 2018-09-18T21:38:37.518000Z "forward" "-" `))

		inputJson := readJsonFromFile(t, "./testdata/s3-event.json")
		var s3Event events.S3Event
		if err := json.Unmarshal(inputJson, &s3Event); err != nil {
			t.Errorf("could not unmarshal event: %v", err)
		}

		err := handler(s3Event)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(firehoseMock.records) == 0 {
			t.Errorf("expected entries")
		}

		var entry AlbLogEntry
		err = json.Unmarshal(firehoseMock.records[0], &entry)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if entry.Request != "GET http://test1-332132803.eu-west-1.elb.amazonaws.com:80/ HTTP/1.1" {
			t.Errorf("unexpected request: %v", entry.Request)
		}
	})
}

func gz(str string) io.Reader {
	var buf bytes.Buffer
	b := []byte(str)
	gz := gzip.NewWriter(&buf)
	gz.Write(b)
	gz.Flush()
	gz.Close()
	return &buf
}

func TestReadLogEntries(t *testing.T) {
	logEntries, err := readLogEntries(ioutil.NopCloser(gz(
		`http 2018-09-18T21:38:37.519183Z app/test1/7f050ffab5373730 95.90.211.80:4254 172.31.7.183:80 0.001 0.000 0.000 200 200 461 654 "GET http://test1-332132803.eu-west-1.elb.amazonaws.com:80/ HTTP/1.1" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36" - - arn:aws:elasticloadbalancing:eu-west-1:328305143014:targetgroup/target-group-1/76784dcf3348e51f "Root=1-5ba1705d-5617cf32286ecfb8172f73cc" "-" "-" 0 2018-09-18T21:38:37.518000Z "forward" "-" `)))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entry := logEntries[0]
	if entry.Type != "http" {
		t.Errorf("invalid type")
	}
	if entry.Timestamp != "2018-09-18T21:38:37.519183Z" {
		t.Errorf("invalid timestamp")
	}
	if entry.Client != "95.90.211.80" {
		t.Errorf("invalid client")
	}
	if entry.Target != "172.31.7.183" {
		t.Errorf("invalid target")
	}
	if entry.Request != "GET http://test1-332132803.eu-west-1.elb.amazonaws.com:80/ HTTP/1.1" {
		t.Errorf("invalid request, %v", entry.Request)
	}
}

func readJsonFromFile(t *testing.T, inputFile string) []byte {
	inputJson, err := ioutil.ReadFile(inputFile)
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	return inputJson
}
