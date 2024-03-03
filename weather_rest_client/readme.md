# Weather REST client

Simple REST client for the weather API.

See [readme of the SAM stack](../weather_api/readme.md) for details regarding:
* the REST URL 
* the API key
* the creation of the client certificates

Usage:

```sh
go run . \
    -url https://rest.weather-api-demo.poc.svend.xyz/weather  \
    -deviceId <device-id> \
    -timeDelta <some-duration-in-minutes> \
    -apiKey <api-key> \
	-certFile certificates/clientCert.pem \
	-keyFile certificates/clientKey.pem
```
