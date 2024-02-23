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

Tail the logs:

```sh
sam logs -n WeatherReadFrontendFunction --stack-name weather-api-demo --tail
```

Get some weather events (or use the rest client):

```sh
# any device_id between 1000 and 1009
# '%2B' is the '+' in the ISO8601 timestamp
curl -X GET 'https://...amazonaws.com/Prod/weather?device_id=1005&from=2023-02-17T20:13:25%2B0100&to=2025-02-17T20:13:55%2B0100'
```

