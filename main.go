package main

import (
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"errors"
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
	Type                   string `json:"type"`
	Timestamp              string `json:"@timestamp"`
	Elb                    string `json:"type"`
	Client                 string `json:"client"`
	Target                 string `json:"target"`
	RequestProcessingTime  string `json:"request_processing_time"`
	TargetProcessingTime   string `json:"target_processing_time"`
	ResponseProcessingTime string `json:"response_processing_time"`
	ElbStatusCode          string `json:"elb_status_code"`
	TargetStatusCode       string `json:"target_status_code"`
	ReceivedBytes          string `json:"received_bytes"`
	SentBytes              string `json:"sent_bytes"`
	Request                string `json:"request"`
	UserAgent              string `json:"user_agent"`
	SslCipher              string `json:"ssl_cipher"`
	SslProtocol            string `json:"ssl_protocol"`
	TargetGroupArn         string `json:"target_group_arn"`
	TraceId                string `json:"trace_id"`
	DomainName             string `json:"domain_name"`
	ChosenCertArn          string `json:"chosen_cert_arn"`
	MatchedRulePriority    string `json:"matched_rule_priority"`
	RequestCreationTime    string `json:"request_creation_time"`
	ActionsExecuted        string `json:"actions_executed"`
	RedirectUrl            string `json:"redirect_url"`
}

func handler(s3Event events.S3Event) error {
	for _, rec := range s3Event.Records {
		closer, _ := readS3File(rec.S3.Bucket.Name, rec.S3.Object.Key)
		logEntries := read(closer)
		for _, l := range logEntries {
			fmt.Println(*l)
		}
	}

	return nil
}

func readS3File(bucket string, key string) (io.ReadCloser, error) {
	log.Printf("reading %v/%v", bucket, key)
	cfg, _ := external.LoadDefaultAWSConfig()
	s3svc := s3.New(cfg)
	req := s3svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func read(closer io.ReadCloser) []*LogEntry {
	defer func() {
		closer.Close()
		if r := recover(); r != nil {
			log.Printf("recovered from '%v', skipping entry", r)
		}
	}()

	var logEntries []*LogEntry
	gzReader, _ := gzip.NewReader(closer)
	br := bufio.NewReader(gzReader)
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			break
		}
		logEntry, err := createLogEntryFromLine(line)
		if err == nil {
			logEntries = append(logEntries, logEntry)
		} else {
			log.Println(err)
		}
	}
	return logEntries
}

func readLogEntries(line string) (entry *LogEntry, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Unable to parse line '%v', error: %v", line, r))
		}
	}()

	line = strings.Trim(line, "\n")
	r := csv.NewReader(strings.NewReader(line))
	r.Comma = ' '

	records, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	cols := records[0]
	return &LogEntry{
		Type:                   cols[0],
		Timestamp:              cols[1],
		Elb:                    cols[2],
		Client:                 cols[3],
		Target:                 cols[4],
		RequestProcessingTime:  cols[5],
		TargetProcessingTime:   cols[6],
		ResponseProcessingTime: cols[7],
		ElbStatusCode:          cols[8],
		TargetStatusCode:       cols[9],
		ReceivedBytes:          cols[10],
		SentBytes:              cols[11],
		Request:                strings.Trim(cols[12], "\""),
		UserAgent:              strings.Trim(cols[13], "\""),
		SslCipher:              cols[14],
		SslProtocol:            cols[15],
		TargetGroupArn:         cols[16],
		TraceId:                strings.Trim(cols[17], "\""),
		DomainName:             strings.Trim(cols[18], "\""),
		ChosenCertArn:          strings.Trim(cols[19], "\""),
		MatchedRulePriority:    cols[20],
		RequestCreationTime:    cols[21],
		ActionsExecuted:        strings.Trim(cols[22], "\""),
		RedirectUrl:            strings.Trim(cols[23], "\""),
	}, nil
}

func createLogEntryFromLine(line string) (entry *LogEntry, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Unable to parse line '%v', error: %v", line, r))
		}
	}()

	line = strings.Trim(line, "\n")
	r := csv.NewReader(strings.NewReader(line))
	r.Comma = ' '

	records, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	cols := records[0]
	return &LogEntry{
		Type:                   cols[0],
		Timestamp:              cols[1],
		Elb:                    cols[2],
		Client:                 cols[3],
		Target:                 cols[4],
		RequestProcessingTime:  cols[5],
		TargetProcessingTime:   cols[6],
		ResponseProcessingTime: cols[7],
		ElbStatusCode:          cols[8],
		TargetStatusCode:       cols[9],
		ReceivedBytes:          cols[10],
		SentBytes:              cols[11],
		Request:                strings.Trim(cols[12], "\""),
		UserAgent:              strings.Trim(cols[13], "\""),
		SslCipher:              cols[14],
		SslProtocol:            cols[15],
		TargetGroupArn:         cols[16],
		TraceId:                strings.Trim(cols[17], "\""),
		DomainName:             strings.Trim(cols[18], "\""),
		ChosenCertArn:          strings.Trim(cols[19], "\""),
		MatchedRulePriority:    cols[20],
		RequestCreationTime:    cols[21],
		ActionsExecuted:        strings.Trim(cols[22], "\""),
		RedirectUrl:            strings.Trim(cols[23], "\""),
	}, nil
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
