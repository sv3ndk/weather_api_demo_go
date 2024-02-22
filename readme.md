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
* setup mutual TLS https://docs.aws.amazon.com/apigateway/latest/developerguide/rest-api-mutual-tls.html
