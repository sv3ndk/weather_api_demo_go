## How to run

```sh
sam build

# only the first time, then choose to same settings to file
sam deploy --guided

# all subsequent deployments
sam deploy
```

Obtain the URL of the public API:

```sh
    aws cloudformation describe-stacks --stack-name weather-api-demo --query 'Stacks[0].Outputs[?OutputKey==`WeatherAPI`].OutputValue' --output text
```

Obtain the API key id
```
# Customer1ApiKeyId or Customer2ApiKeyId here
aws cloudformation describe-stacks --stack-name weather-api-demo --query 'Stacks[0].Outputs[?OutputKey==`Customer1ApiKeyId`].OutputValue' --output text


aws apigateway get-api-key --include-value --api-key <some id> --query "value"
```

Tail the logs:

```sh
sam logs -n WeatherReadFrontendFunction --stack-name weather-api-demo --tail
```

Get some weather events (or use the rest client):

```sh
# query parameters:
#   device_id: int in [1000, 1009]
#   from, to: recent ISO8601 timestamp (random data are generated in real time every second), not that '%2B' is the '+' in the ISO8601 format
# url, api key: see above

curl \
    -X GET 'https://...url...v1/weather?device_id=1005&from=2023-02-17T20:13:25%2B0100&to=2025-02-17T20:13:55%2B0100' \
    -H 'X-API-Key: ...api key...'
```

