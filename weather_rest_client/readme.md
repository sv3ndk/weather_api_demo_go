# Weather REST client

Simple REST client for the weather API.

Usage:

```sh
# see [readme of the SAM stack](../weather_api/readme.md) for details on obtaining URL and API key
go run . \
    -url <api-gateway-endpoint>  \
    -apiKey <api-key> \
    -deviceId <device-id> \
    -timeDelta <some-duration-in-minutes>
```
