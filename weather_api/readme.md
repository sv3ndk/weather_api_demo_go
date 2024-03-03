# Weather API

## Deployment

### Custom domain pre-requisites

Both the REST and Websocek endpoints are associated with a public DNS subdomain name, which requires some manual additional setup
(or just to comment out the domain and mapping config in the SAM template...)

* register a domain name
* create a [certificate for `*.weather-api-demo.poc.domain-name` in ACM](https://docs.aws.amazon.com/acm/latest/userguide/gs-acm-request-public.html)

=> use the ARN of that certificate as `DomainCertificateArn` input parameter of the SAM template.

### mutual TLS pre-requisite

The SAM stack activate mutual TLS on the REST API, which requires a trustore to be manually created and uploaded.

* execute [gen-certs.sh](../weather_rest_client/certificates/gen-certs.sh) to create a fake root CA, a client private key 
  and certificate and a trustore to be used by the REST server
* upload the trustore to S3

=> use the path of that trustore on S3 as `RestApiMtlsTruststore` input parameter of the SAM template.

### Stack deployment

Build and deploy the SAM application:

```sh
sam build

# only the first time, then choose to same settings to file
sam deploy --guided

# all subsequent deployments
sam deploy
```

### DNS registration

The stack contains 2 API Gateway custom domain mappings that needs to be associated with the desired DNS name of the services.

```sh
# AWS mapping domain name for the rest endpoint:
aws cloudformation describe-stacks --stack-name weather-api-demo --query 'Stacks[0].Outputs[?OutputKey==`WeatherRestRegionalAwsDomain`].OutputValue' --output text

# public DNS name of the rest endpoint
aws cloudformation describe-stacks --stack-name weather-api-demo --query 'Stacks[0].Outputs[?OutputKey==`WeatherRestPublicDomain`].OutputValue' --output text
```

Add a DNS CNAME that let the public DNS name resolve to the AWS mapping name

```sh
# AWS mapping domain name for the websocket endpoint:
aws cloudformation describe-stacks --stack-name weather-api-demo --query 'Stacks[0].Outputs[?OutputKey==`WeatherWsRegionalAwsDomain`].OutputValue' --output text

# public DNS name of the websocket endpoint
aws cloudformation describe-stacks --stack-name weather-api-demo --query 'Stacks[0].Outputs[?OutputKey==`WeatherWsPublicDomain`].OutputValue' --output text
```

Add a DNS CNAME that let the public DNS name resolve to the AWS mapping name

```
# mapping domain name for the websocket endpoint:
aws cloudformation describe-stacks --stack-name weather-api-demo --query 'Stacks[0].Outputs[?OutputKey==`WeatherWsRegionalAwsDomain`].OutputValue' --output text
```

### Logs

```sh
sam logs -n WeatherReadFrontendFunction --stack-name weather-api-demo --tail
sam logs -n WeatherEventWSPushFunction --stack-name weather-api-demo --tail
sam logs -n WeatherDataGeneratorFunction --stack-name weather-api-demo --tail
```

## How to invoke

### Obtain the REST endpoints and credentials

Use the puclic DNS name of the REST service (see above).

Obtain one of the API key id:
```sh
# use Customer1ApiKeyId or Customer2ApiKeyId here
aws cloudformation describe-stacks --stack-name weather-api-demo --query 'Stacks[0].Outputs[?OutputKey==`Customer1ApiKeyId`].OutputValue' --output text

aws apigateway get-api-key --include-value --api-key <some id> --query "value"
```

### Obtain the Web endpoints

Use the puclic DNS name of the websocket service (see above).

### Invoke the service

See the [REST client](../weather_rest_client/readme.md) 
and the [websocket client](../weather_ws_client/readme.md) 
to query each service.

For the REST service, we can also simply use `curl` :

```sh
curl GET \
    'https://rest.weather-api-demo.poc.svend.xyz/weather?device_id=1005&from=2023-02-17T20:13:25%2B0100&to=2025-02-17T20:13:55%2B0100' \
    -H 'X-API-Key: <api key>' \
    --key ../weather_rest_client/certificates/clientKey.pem \
    --cert ../weather_rest_client/certificates/clientCert.pem
```
