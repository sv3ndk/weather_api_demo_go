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

Send a request:

```sh
curl -X GET 'https://...amazonaws.com/Prod/weather?device_id=123&from=2023-02-17T20:13:25%2B0100&to=2025-02-17T20:13:55%2B0100'
```

## References

### Golang AWS SDK

* High level discussion:
  https://aws.github.io/aws-sdk-go-v2/docs/

* AWS GO SDK v2:
  https://aws.github.io/aws-sdk-go-v2/docs/making-requests/ 

### Golang AWS DynamoDB SDK  

* API client, operations, and parameter types for Amazon DynamoDB. 
  https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/dynamodb#Client

* Types and functions to create Amazon DynamoDB Expression strings, ExpressionAttributeNames maps, and ExpressionAttributeValues maps
  https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression#hdr-Using_the_Package

* Marshal slices, maps, structs, and scalar values to and from the AttributeValue type
  https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue

* AWS GO SDK examples
  https://github.com/awsdocs/aws-doc-sdk-examples/tree/main/go

### AWS Lambda implementation in go

* Go lambda handler
  https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html