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

Then simply send a request:

```sh
curl -X GET
```