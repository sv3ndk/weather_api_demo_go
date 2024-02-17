# weather API demo

Demo of a AWS API Gateway usage with a REST endpoint and a websocket connection

## How to use

```sh
sam build

# only the first time, then choose to same settings to file
sam deploy --guided

# all subsequent deployments
sam deploy
```

## Status

### Done

* scafold and deploy fist empty dynamo table

### TODO

* go side project ran on the laptop to populate table with dummy weather events
* follow SAM/GW tuto to create basic REST read interface

