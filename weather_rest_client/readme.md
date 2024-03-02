# Weather REST client

Simple REST client for the weather API.

See [readme of the SAM stack](../weather_api/readme.md) for details on obtaining the REST URL and API key

Usage:

```sh
go run . \
    -url https://rest.weather-api-demo.poc.svend.xyz/weather  \
    -apiKey <api-key> \
    -deviceId <device-id> \
    -timeDelta <some-duration-in-minutes>
```
