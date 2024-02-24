# Weather API demo

Demo of a AWS API Gateway app with a REST endpoint and a websocket connection

## Status

- Data generator lambda triggered every minute to add random events to DynamoDB
- REST API exposed via API Gateway to query events from DynamoDB
- The REST API authorization and throttling are based on API keys
- CLI client app to query the REST endpoint and obtain the last n minutes of events

## TODO

* add websocket API to be notified of any new event + demo that prints stuff to stdout (should be startable several times)
* add Webocket security: first request a temp token through REST, then pass it in the `connect` wss phase
* add custom domain name
* setup mutual TLS for the REST endoipn (cf https://docs.aws.amazon.com/apigateway/latest/developerguide/rest-api-mutual-tls.html
   and https://venilnoronha.io/a-step-by-step-guide-to-mtls-in-go )
* improve data generator: use a step function to parallelize per battery (useless, but I want to...)
* re-use code across packages (data model)
* add OpenAPI spec to REST endpoint
* review duplicated makefiles and code folder structure

## References

### Go AWS SDK

* High level discussion
  https://aws.github.io/aws-sdk-go-v2/docs/

* AWS GO SDK v2
  https://aws.github.io/aws-sdk-go-v2/docs/making-requests/ 

* Go lambda handler
  https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
  
* Lambda events for each kind of integration
  https://github.com/aws/aws-lambda-go/tree/main/events

### Go DynamoDB SDK  

* API client, operations, and parameter types for Amazon DynamoDB
  https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/dynamodb#Client

* Types and functions to create Amazon DynamoDB Expression strings, ExpressionAttributeNames maps, and ExpressionAttributeValues maps
  https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression#hdr-Using_the_Package

* Marshal slices, maps, structs, and scalar values to and from the AttributeValue type
  https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue

* AWS GO SDK examples
  https://github.com/awsdocs/aws-doc-sdk-examples/tree/main/go

### Security

* API Gateways usage plans and API keys
  https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-api-usage-plans.html
  
