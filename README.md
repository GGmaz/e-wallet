# wallet-arringo

## Prerequisites
- Go 1.22
- Docker

## Startup
`docker-compose up` 

## Notes
Swagger init: `swag init -d internal/server/api/v1 -g ../../server.go --parseDependency --parseInternal`

Swagger is on http://localhost:8082/swagger/index.html

### HMAC-SHA512
When you want to send a request with a body, you need to hash the body with the HMAC-SHA512 algorithm and send the hash in the `X-Authorization-Sign` header (secretKey is in .env file).

Hashing algorithm: HMAC-SHA512 - https://devglan.com/online-tools/hmac-sha256-online (do not forget to select SHA512 algorithm)