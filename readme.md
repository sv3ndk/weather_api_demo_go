# Weather API demo

Demo of a AWS API Gateway app with a REST endpoint and a websocket connection

## Status

- Basic REST API returning value read from DynamoDB (no security)
- Lambda function triggered every minute adding some events to Dynamo
- Rest client sending basic GET request

## TODO

* add websocket API to be notified of any new event + demo that prints stuff to stdout (should be startable several times)
* improve data generator: use a step function to parallelize per battery (useless, but I want to...)
* re-use code across packages (data model)
* add custom domain name
* setup mutual TLS https://docs.aws.amazon.com/apigateway/latest/developerguide/rest-api-mutual-tls.html

## References

### Golang AWS SDK

* High level discussion:
  https://aws.github.io/aws-sdk-go-v2/docs/

* AWS GO SDK v2:
  https://aws.github.io/aws-sdk-go-v2/docs/making-requests/ 

* Lambda events for each kind of integration:
  https://github.com/aws/aws-lambda-go/tree/main/events

### Golang AWS DynamoDB SDK  

* API client, operations, and parameter types for Amazon DynamoDB. 
  https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/dynamodb#Client

* Types and functions to create Amazon DynamoDB Expression strings, ExpressionAttributeNames maps, and ExpressionAttributeValues maps
  https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression#hdr-Using_the_Package

* Marshal slices, maps, structs, and scalar values to and from the AttributeValue type
  https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue

* AWS GO SDK examples
  https://github.com/awsdocs/aws-doc-sdk-examples/tree/main/go

### AWS Lambda implementation in go

* Go lambda handler
  https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
