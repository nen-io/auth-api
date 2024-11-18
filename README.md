# GOLANG Auth API

This is a simple API that uses JWT to authenticate users. It is written in Go and uses the Fiber framework.
Also uses a postgres database to store user information.

## Installation

1. Clone the repository
2. install deps `go mod download`
3. run `make run`
4. hot reload for dev `make dev`
5. build `make build`

## Endpoints:

### POST /register
Registers a new user. Requires a JSON body with the following fields:
```
{
	"firstName": "",
	"lastName": "",
	"email": "",
	"password": ""
}
```

### POST /login
Logs in a user. Requires a JSON body with the following fields:
```
{
  "email": "",
  "password": ""
}
```

### GET /verify/:token
Verifies users email

### POST /request-verify
Request a new verification email. Requires a JSON body with the following fields:
```
{
  "email": ""
}
```

### GET /health
Returns a 200 status if the server is running
