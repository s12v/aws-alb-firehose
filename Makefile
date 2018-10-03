.PHONY: all
all: deps clean test build

deps:
	dep ensure -v

clean: 
	rm ./aws-alb-firehose || true

test:
	go test -v -coverprofile=coverage.txt -covermode=atomic

build:
	GOOS=linux GOARCH=amd64 go build -o ./aws-alb-firehose

samtest:
	sam validate

samrun: clean test build
	sam local invoke AlbToKinesisFunction -e testdata/s3-event.json --env-vars testdata/env.json

sampackage: clean test build samtest
	@[ "${bucket}" ] || ( echo ">> bucket is not defined"; exit 1 )
	sam package --template-file template.yaml --s3-bucket $(bucket) --output-template-file packaged.yaml

samdeploy:
	@[ "${bucket}" ] || ( echo ">> bucket is not defined"; exit 1 )
	@[ "${stack_name}" ] || ( echo ">> stack_name is not defined"; exit 1 )
	@[ "${delivery_stream_name}" ] || ( echo ">> delivery_stream_name is not defined"; exit 1 )
	aws cloudformation deploy --template-file packaged.yaml --capabilities CAPABILITY_IAM --stack-name $(stack_name) --parameter-overrides DeliveryStreamName=$(delivery_stream_name) S3ALBLogsBucketName=$(bucket)
