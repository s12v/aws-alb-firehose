package main

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/firehose"
	"log"
)


type FirehoseWriter interface {
	write(entries [][]byte) (failedPutCount int64, err error)
}


type defaultFirehoseWriter struct {
	FirehoseWriter
	deliveryStreamName string
	firehoseSvc *firehose.Firehose
}

func NewDefaultFirehoseWriter(cfg aws.Config, deliveryStreamName string) FirehoseWriter {
	return &defaultFirehoseWriter{
		deliveryStreamName: deliveryStreamName,
		firehoseSvc: firehose.New(cfg),
	}
}

func (writer *defaultFirehoseWriter) write(entries [][]byte) (failedPutCount int64, err error) {
	var records []firehose.Record
	for _, entry := range entries {
		records = append(records, firehose.Record{Data: entry})
	}

	input := firehose.PutRecordBatchInput{
		DeliveryStreamName: &writer.deliveryStreamName,
		Records: records,
	}
	request := writer.firehoseSvc.PutRecordBatchRequest(&input)
	output, err := request.Send()
	if output != nil && output.FailedPutCount != nil {
		failedPutCount = *output.FailedPutCount
	}
	log.Printf("sent %v records to %v, %v failed", len(entries), writer.deliveryStreamName, failedPutCount)
	return
}

type mockFirehoseWriter struct {
	records [][]byte
}

func (writer *mockFirehoseWriter) write(records [][]byte) (int64, error) {
	writer.records = append(writer.records, records...)
	return 0, nil
}

func NewMockFirehoseWriter() *mockFirehoseWriter {
	return &mockFirehoseWriter{}
}
