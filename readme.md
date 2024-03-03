# Weather API demo

## Status

- REST integration:
  * a [REST API](weather_api/weather_rest_frontend/main.go) exposed via the API Gateway allows to query weather events.
  * a [CLI client app](weather_rest_client/readme.md) queries this REST endpoint
  * API keys are configured to limit traffic (usage/quotas)
  * authentication is based on mutual TLS 

- websocket integration: 
  * a websocket endpoint is exposed on the API Gateway
  * the [on-connect lambda](weather_api/weather_ws_on_connection_event/main.go) keeps track of the currently connected websocket clients
  * the [ws-push lambda](weather_api/weather_event_ws_push/main.go) is notified when events are added to DynamoDB and forwards them to all currently connected websocket clients
  * a [CLI websocket client](weather_ws_client/readme.md) streams weather events from the websocket endpoint and prints them

- both the REST and websocket endpoints are exposed on a custom DNS domain

- a [data generator lambda](weather_api/weather_data_generator/main.go), triggered every minute, adds random weather events to DynamoDB

## TODO (maybe)

* handle SIGINT correcty in ws socket client
* add Webocket security: API key? Or first request a temp token through REST, then pass it in the `connect` ws phase
* add OpenAPI spec to REST endpoint

## References

### Go AWS SDK

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
  
* Blog on mutual TLS using the API Gateway
  https://aws.amazon.com/blogs/compute/introducing-mutual-tls-authentication-for-amazon-api-gateway/