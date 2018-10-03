# Application load bakancer (ALB) to Kinesis Firehose lambda

Serverless function to stream access logs of Application ELB from S3 to Amazon Kinesis Firehose.

This SAM template creates the Lambda function & associated policy + IAM role, and new S3 bucket
with enabled events notifications to this Lambda function.

Send your ALB access logs to this newly created S3 Bucket. To enable access logging for ALB:
http://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-access-logs.html#enable-access-logging

Built using AWS Serverless Application Model

Configuration:
 * `S3ALBLogsBucketName` - the name of the S3 bucket
 * `DeliveryStreamName` - the name of the Firehose delivery stream

Useful commands:

* Generate a test event: `sam local generate-event s3 put --bucket foo --key bar > event.json`
* Test locally: `make samtest`
* SAM validate: `make samvalidate`
* SAM package: `make sampackage bucket=sambucket-alb`
* SAM deploy: `make samdeploy stack_name=alb-firehose-lambda delivery_stream_name=test1 bucket=sambucket-alb` 
