.PHONY: deps clean build

deps:
	go get -u ./...

clean: 
	rm -rf ./alb_log_processor/aws-alb-firehose
	
build:
	GOOS=linux GOARCH=amd64 go build -o ./aws-alb-firehose ./alb_log_processor
