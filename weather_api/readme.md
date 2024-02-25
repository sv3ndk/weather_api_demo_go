# Weather API

## How to run

Build and deploy the SAM application:

```sh
sam build

# only the first time, then choose to same settings to file
sam deploy --guided

# all subsequent deployments
sam deploy
```

Obtain the URL of the public REST API:

```sh
aws cloudformation describe-stacks --stack-name weather-api-demo --query 'Stacks[0].Outputs[?OutputKey==`WeatherAPIRestEndpoint`].OutputValue' --output text
```

Obtain one of the API key id:
```sh
# use Customer1ApiKeyId or Customer2ApiKeyId here
aws cloudformation describe-stacks --stack-name weather-api-demo --query 'Stacks[0].Outputs[?OutputKey==`Customer1ApiKeyId`].OutputValue' --output text

aws apigateway get-api-key --include-value --api-key <some id> --query "value"
```

Obtain the URL of the websocket endpoint:
```sh
aws cloudformation describe-stacks --stack-name weather-api-demo --query 'Stacks[0].Outputs[?OutputKey==`WeatherAPIWsEndpoint`].OutputValue' --output text
```

Tail the logs:

```sh
sam logs -n WeatherReadFrontendFunction --stack-name weather-api-demo --tail
sam logs -n WeatherEventWSPushFunction --stack-name weather-api-demo --tail
sam logs -n WeatherDataGeneratorFunction --stack-name weather-api-demo --tail
```

Get some weather events (or use the client projects to query the REST or websocket endpoint):

```sh
curl \
    -X GET 'https://<httpUrl>?device_id=1005&from=2023-02-17T20:13:25%2B0100&to=2025-02-17T20:13:55%2B0100' \
    -H 'X-API-Key: <api key>'
```
