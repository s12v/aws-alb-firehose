package main

import (
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"io"
	"log"
	"os"
)

var s3Reader S3Reader
var firehoseWriter FirehoseWriter

func handler(s3Event events.S3Event) error {
	log.Print("handler()")
	for _, rec := range s3Event.Records {
		closer, err := s3Reader.readS3File(rec.S3.Bucket.Name, rec.S3.Object.Key)
		if err != nil {
			log.Printf("skip %v/%v, error: %v", rec.S3.Bucket.Name, rec.S3.Object.Key, err)
			continue
		}
		logEntries, err := readLogEntries(closer)
		if err != nil {
			log.Printf("skip %v/%v, error: %v", rec.S3.Bucket.Name, rec.S3.Object.Key, err)
			continue
		}

		var records [][]byte
		for _, l := range logEntries {
			record, err := json.Marshal(l)
			if err != nil {
				log.Printf("drop entry %v, error: %v", l, err)
				continue
			}
			records = append(records, record)
			if len(records) == 500 {
				firehoseWriter.write(records)
				records = nil
			}
		}

		if records != nil {
			firehoseWriter.write(records)
		}
	}

	return nil
}

func readLogEntries(closer io.ReadCloser) (logEntries []*AlbLogEntry, err error) {
	defer func() {
		closer.Close()
	}()

	gzipReader, _ := gzip.NewReader(closer)
	csvReader := csv.NewReader(gzipReader)
	csvReader.Comma = ' '
	records, err := csvReader.ReadAll()
	if err != nil {
		return
	}

	for _, r := range records {
		logEntry, err := CreateLogEntry(r)
		if err == nil {
			logEntries = append(logEntries, logEntry)
		} else {
			log.Printf("drop entry %v, error: %v", r, err)
		}
	}
	return
}

func main() {
	log.Print("main()")
	deliveryStreamName := os.Getenv("DELIVERY_STREAM_NAME")
	log.Printf("stream name: %v", deliveryStreamName)
	awsConfig, _ := external.LoadDefaultAWSConfig()
	s3Reader = NewDefaultS3Reader(awsConfig)
	firehoseWriter = NewDefaultFirehoseWriter(awsConfig, deliveryStreamName)
	lambda.Start(handler)
}
