# Weather API demo

## Status

- a [data generator lambda](weather_api/weather_data_generator/main.go), triggered every minute, adds random weather events to DynamoDB
- REST integration:
  * a [REST API](weather_api/weather_rest_frontend/main.go) exposed via the API Gateway allows to query those events. 
  * Its access and throttling is based on API keys
  * a [CLI client app](weather_rest_client/weather_client/client.go) queries the REST endpoint and obtain the last n minutes of events
- websocket integration: 
  * the [on-connect lambda](weather_api/weather_ws_on_connection_event/main.go) is notified when a websocket client connects or disconnect and keeps track of its connection id
  * the [ws-push lambda](weather_api/weather_event_ws_push/main.go) is notified when events are added to Dynamodb and forwards them in JSON to all connected clients
  * a [CLI websocket client](weather_ws_client/main.go) allows to print weather events as they arrive

## TODO (maybe)

* add Webocket security: API key? Or first request a temp token through REST, then pass it in the `connect` wss phase
* add custom domain name
* setup mutual TLS for the REST endoipn (cf https://docs.aws.amazon.com/apigateway/latest/developerguide/rest-api-mutual-tls.html
   and https://venilnoronha.io/a-step-by-step-guide-to-mtls-in-go )
* improve data generator: use a step function to parallelize per battery (useless, but I want to...)
* lot's of copy/pasted code should be cleaned up :)
* add OpenAPI spec to REST endpoint
* review duplicated makefiles and code folder structure

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
  