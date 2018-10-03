.PHONY: all
all: deps clean test build

deps:
	dep ensure -v

clean: 
	rm ./aws-alb-firehose || true

test:
	go test -v

build:
	GOOS=linux GOARCH=amd64 go build -o ./aws-alb-firehose

samtest:
	sam validate

samrun: clean test build
	sam local invoke AlbToKinesisFunction -e testdata/s3-event.json --env-vars testdata/env.json

sampackage: clean test build samtest
    ifndef bucket
        $(error bucket is not defined)
    endif
	sam package --template-file template.yaml --s3-bucket $(bucket) --output-template-file packaged.yaml

samdeploy:
    ifndef stack_name
        $(error stack_name is not defined)
    endif
    ifndef delivery_stream_name
        $(error delivery_stream_name is not defined)
    endif
    ifndef bucket
        $(error bucket is not defined)
    endif
	aws cloudformation deploy --template-file packaged.yaml --capabilities CAPABILITY_IAM --stack-name $(stack_name) --parameter-overrides DeliveryStreamName=$(delivery_stream_name) S3ALBLogsBucketName=$(bucket)
