package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"log"
	"strings"
)

type LogEntry struct {
	Type      string `json:"type"`
	Timestamp string `json:"@timestamp"`
	Elb       string `json:"type"`
}

func handler(s3Event events.S3Event) (string, error) {

	for i, rec := range s3Event.Records {
		fmt.Print(i)
		bucket := rec.S3.Bucket.Name
		key := rec.S3.Object.Key
		log.Printf("processing %v/%v", bucket, key)
		readS3File(bucket, key)
	}

	return "boom", nil
}

func readS3File(bucket string, key string) {
	cfg, _ := external.LoadDefaultAWSConfig()
	s3svc := s3.New(cfg)
	req := s3svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	resp, err := req.Send()
	if err != nil {
		fmt.Println(err)
	} else {
		gzReader, _ := gzip.NewReader(resp.Body)
		br := bufio.NewReader(gzReader)
		for {
			line, err := br.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				// boom
			}
			line = strings.Trim(line, "\n")
			cols := strings.Split(line, " ")
			logEntry := &LogEntry{
				Type: cols[0],
				Timestamp: cols[1],
				Elb: cols[2],
			}
			fmt.Println(logEntry)
		}
	}
}

//func processLogFile(bucket string, key string) {
//	req := s3svc.GetObjectRequest(&s3.GetObjectInput{
//		Bucket: &bucket,
//		Key:    &key,
//	})
//	resp, err := req.Send()
//	if err == nil {
//		fmt.Println(resp)
//	}
//}

func main() {
	//cfg, _ := external.LoadDefaultAWSConfig()
	//s3svc = s3.New(cfg)
	lambda.Start(handler)
}
