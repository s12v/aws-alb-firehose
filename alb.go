package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type AlbLogEntry struct {
	Type                   string  `json:"type"`
	Timestamp              string  `json:"@timestamp"`
	Elb                    string  `json:"elb"`
	Client                 string  `json:"client"`
	ClientPort             int64   `json:"client_port"`
	Target                 string  `json:"target"`
	TargetPort             int64   `json:"target_port"`
	RequestProcessingTime  float64 `json:"request_processing_time"`
	TargetProcessingTime   float64 `json:"target_processing_time"`
	ResponseProcessingTime float64 `json:"response_processing_time"`
	ElbStatusCode          int64   `json:"elb_status_code"`
	TargetStatusCode       int64   `json:"target_status_code"`
	ReceivedBytes          int64   `json:"received_bytes"`
	SentBytes              int64   `json:"sent_bytes"`
	Request                string  `json:"request"`
	UserAgent              string  `json:"user_agent"`
	SslCipher              string  `json:"ssl_cipher"`
	SslProtocol            string  `json:"ssl_protocol"`
	TargetGroupArn         string  `json:"target_group_arn"`
	TraceId                string  `json:"trace_id"`
	DomainName             string  `json:"domain_name"`
	ChosenCertArn          string  `json:"chosen_cert_arn"`
	MatchedRulePriority    string  `json:"matched_rule_priority"`
	RequestCreationTime    string  `json:"request_creation_time"`
	ActionsExecuted        string  `json:"actions_executed"`
	RedirectUrl            string  `json:"redirect_url"`
}

func CreateLogEntry(cols []string) (entry *AlbLogEntry, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Unable to parse '%v', error: %v", cols, r))
		}
	}()

	logEntry := AlbLogEntry{
		Type:                   cols[0],
		Timestamp:              cols[1],
		Elb:                    cols[2],
		Client:                 strings.Split(cols[3], ":")[0],
		ClientPort:             -1,
		Target:                 strings.Split(cols[4], ":")[0],
		TargetPort:             -1,
		RequestProcessingTime:  -1,
		TargetProcessingTime:   -1,
		ResponseProcessingTime: -1,
		ElbStatusCode:          -1,
		TargetStatusCode:       -1,
		ReceivedBytes:          -1,
		SentBytes:              -1,
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
	}

	logEntry.ClientPort, err = strconv.ParseInt(strings.Split(cols[3], ":")[1], 10, 64)
	if err != nil {
		log.Printf("error while parsing log entry %v: %v", cols, err)
	}

	logEntry.TargetPort, err = strconv.ParseInt(strings.Split(cols[4], ":")[1], 10, 64)
	if err != nil {
		log.Printf("error while parsing log entry %v: %v", cols, err)
	}

	logEntry.RequestProcessingTime, err = strconv.ParseFloat(cols[5], 32)
	if err != nil {
		log.Printf("error while parsing log entry %v: %v", cols, err)
	}

	logEntry.TargetProcessingTime, err = strconv.ParseFloat(cols[6], 32)
	if err != nil {
		log.Printf("error while parsing log entry %v: %v", cols, err)
	}

	logEntry.ResponseProcessingTime, err = strconv.ParseFloat(cols[7], 32)
	if err != nil {
		log.Printf("error while parsing log entry %v: %v", cols, err)
	}

	logEntry.ElbStatusCode, err = strconv.ParseInt(cols[8], 10, 32)
	if err != nil {
		log.Printf("error while parsing log entry %v: %v", cols, err)
	}

	logEntry.TargetStatusCode, err = strconv.ParseInt(cols[9], 10, 32)
	if err != nil {
		log.Printf("error while parsing log entry %v: %v", cols, err)
	}

	logEntry.ReceivedBytes, err = strconv.ParseInt(cols[10], 10, 32)
	if err != nil {
		log.Printf("error while parsing log entry %v: %v", cols, err)
	}

	logEntry.SentBytes, err = strconv.ParseInt(cols[11], 10, 32)
	if err != nil {
		log.Printf("error while parsing log entry %v: %v", cols, err)
	}

	return &logEntry, nil
}
