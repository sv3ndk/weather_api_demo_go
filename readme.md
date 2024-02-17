# Weather API demo

Demo of a AWS API Gateway app with a REST endpoint and a websocket connection

* Weather API deployed to AWS [weather_api](weather_api/readme.md)
* Data sample creator: [populate_samples](populate_samples/readme.md)

## Status

- DynamoDB for storing the events + script for populating it
- Basic REST API returning hard-coded value with one GET and no security

### TODO

* let API return value from Dynamo
* setup mutual TLS https://docs.aws.amazon.com/apigateway/latest/developerguide/rest-api-mutual-tls.html
* event bridge cron based processor to add events to Dynamo 
* improve data generator: use a step function to parallelize per battery (useless, but I want to...)
* add websocket API to be notified of any new event + demo that prints stuff to stdout (should be startable several times)
