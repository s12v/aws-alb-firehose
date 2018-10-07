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

func handler(s3Reader S3Reader, writer FirehoseWriter) func(s3Event events.S3Event) error {
	return func(s3Event events.S3Event) error {
		log.Print("handler()")
		for _, rec := range s3Event.Records {
			albLogEntries, err := readAlbLogsFromS3(s3Reader, rec.S3.Bucket.Name, rec.S3.Object.Key)
			if err != nil {
				log.Printf("skip %v/%v, error: %v", rec.S3.Bucket.Name, rec.S3.Object.Key, err)
			}

			sendAlbLogsToFirehose(albLogEntries, writer)
		}

		return nil
	}
}

func readAlbLogsFromS3(s3Reader S3Reader, bucket string, key string) ([]*AlbLogEntry, error) {
	s3File, err := s3Reader.readS3File(bucket, key)
	if err != nil {
		return nil, err
	}
	defer func() {
		s3File.Close()
	}()

	entries, err := readLogEntries(s3File)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func sendAlbLogsToFirehose(entries []*AlbLogEntry, writer FirehoseWriter) {
	var records [][]byte
	for _, l := range entries {
		record, err := json.Marshal(l)
		if err != nil {
			log.Printf("drop entry %v, error: %v", l, err)
			continue
		}
		records = append(records, record)
		if len(records) == 500 {
			writer.write(records)
			records = nil
		}
	}

	if records != nil {
		writer.write(records)
	}
}

func readLogEntries(reader io.Reader) (logEntries []*AlbLogEntry, err error) {
	gzipReader, _ := gzip.NewReader(reader)
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
	lambda.Start(handler(
		NewDefaultS3Reader(awsConfig),
		NewDefaultFirehoseWriter(awsConfig, deliveryStreamName)))
}
