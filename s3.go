package main

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"io/ioutil"
	"log"
)

type S3Reader interface {
	readS3File(bucket string, key string) (io.ReadCloser, error)
}

type defaultS3Reader struct {
	S3Reader
	s3svc *s3.S3
}

func NewDefaultS3Reader(cfg aws.Config) S3Reader {
	return &defaultS3Reader{s3svc: s3.New(cfg)}
}

func (reader *defaultS3Reader) readS3File(bucket string, key string) (io.ReadCloser, error) {
	log.Printf("reading %v/%v", bucket, key)
	req := reader.s3svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}


type mockS3Reader struct {
	reader io.Reader
}

func (mock *mockS3Reader) readS3File(bucket string, key string) (io.ReadCloser, error) {
	return ioutil.NopCloser(mock.reader), nil
}

func NewMockS3Reader(response io.Reader) S3Reader {
	return &mockS3Reader{reader: response}
}
