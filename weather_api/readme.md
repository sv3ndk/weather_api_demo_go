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

Then simply send a request:

```sh
curl -X GET https://../weather?device_id=1&from=2024-02-17T20:13:25Z&to=2024-02-17T20:13:55Z'
```
