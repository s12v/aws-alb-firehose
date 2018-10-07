.PHONY: all
all: deps clean test build

releases_bucket="aws-alb-firehose-releases"
sam_application_id="arn:aws:serverlessrepo:us-east-1:328305143014:applications"

deps:
	dep ensure -v

clean: 
	rm ./aws-alb-firehose || true

test:
	go test -v -coverprofile=coverage.txt -covermode=atomic

build:
	GOOS=linux GOARCH=amd64 go build -o main

samtest:
	sam validate

samrun: clean test build
	sam local invoke AlbToKinesisFunction -e testdata/s3-event.json --env-vars testdata/env.json

# bucket - S3 bucket for lambda binaries
sampackage: clean test build samtest
	sam package --template-file template.yaml --s3-bucket $(releases_bucket) --output-template-file packaged.yaml

release: sampackage
	@[ "${version}" ] || ( echo ">> version is not defined"; exit 1 )
	CODE_URI=`cat packaged.yaml | grep CodeUri | sed -n '/CodeUri/s/  */ /p' | cut -d ' ' -f 3`
	@[ "${CODE_URI}" ] || ( echo ">> CODE_URI is not defined"; exit 1 )
	@echo "code uri: $CODE_URI"
	VERSIONED="s3://$(releases_bucket)/alb-firehose-$TRAVIS_TAG"
	aws s3 mv "$CODE_URI" "$VERSIONED" --acl public-read

	cat packaged.yaml | sed "s~$CODE_URI~$VERSIONED~" > packaged.yaml

	aws serverlessrepo create-application-version \
	--application-id $(sam_application_id) \
	--semantic-version "$TRAVIS_TAG" \
	--template-body packaged.yaml \
	--source-code-url "https://github.com/s12v/aws-alb-firehose/releases/tag/$TRAVIS_TAG"

# Deploy Cloudformation stack
# stack-name - the name of the Cloudformation stack
# delivery_stream_name - the name of Kinesis firehose delivery stream
# bucket - S3 bucket for ALB logs (will be created)
samdeploy:
	@[ "${bucket}" ] || ( echo ">> bucket is not defined"; exit 1 )
	@[ "${stack_name}" ] || ( echo ">> stack_name is not defined"; exit 1 )
	@[ "${delivery_stream_name}" ] || ( echo ">> delivery_stream_name is not defined"; exit 1 )
	aws cloudformation deploy --template-file packaged.yaml --capabilities CAPABILITY_IAM --stack-name $(stack_name) --parameter-overrides DeliveryStreamName=$(delivery_stream_name) S3ALBLogsBucketName=$(bucket)

