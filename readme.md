# Weather API demo

Demo of a AWS API Gateway app with a REST endpoint and a websocket connection

* Weather API deployed to AWS [weather_api](weather_api/readme.md)
* Data sample creator: [populate_samples](populate_samples/readme.md)

## Status

### Done

* scafold SAM project and deploy fist empty dynamo table
* create CLI data sample creator

### TODO

* follow SAM/GW tuto to create basic REST read interface
* event bridge cron based processor to add events to Dynamo 
* improve data generator: use a step function to parallelize per battery (useless, but I want to...)
* add websocket API to be notified of any new event + demo that prints stuff to stdout (should be startable several times)
