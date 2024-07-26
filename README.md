# Auth

![Tests](https://github.com/abyan-dev/auth/actions/workflows/ci.yaml/badge.svg) [![codecov](https://codecov.io/gh/abyan-dev/auth/graph/badge.svg?token=S679A5TSW7)](https://codecov.io/gh/abyan-dev/auth) [![Go Report](https://goreportcard.com/badge/abyan-dev/auth)](https://goreportcard.com/report/YanSystems/compiler) [![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/YanSystems/compiler/blob/main/LICENSE)

Authentication and authorization service that exposes a REST API for issuing,  invalidating, and revoking HS256-encrypted JSON Web Tokens (JWT). Tokens are stored in `Httponly`, `Secure`, and `Strictly Samesite` cookies, thereby minimizing XSS and CSRF vulnerabilities. 

**Table of Contents:**

- [Features](#features)
- [Quickstart](#quickstart)
- [API](#api)

## Features

The service currently provides the following features:

- Credentials-based registration, login, and logout
- Account verification through email
- Forgot password / password reset mechanism
- Two-factor authentication by emailing shortlived JWT

## Quickstart

To integrate this service to your application, build a linux executable and use it to create a lightweight docker image on alpine OS:

```
env GOOS=linux CGO_ENABLED=0 go build -o auth ./cmd/api && make image
```

The app is dependent on a postgres instance. Run the following to quickly spin up one: 

```
make db-up
```

Now serve the container:

```
docker run -p 8080:8080 --name auth auth
```

Your application can interact with the service at `localhost:8080`. Alternatively, if you are using an orchestrator like Compose or Kubernetes, you can leverage their DNS by using `auth:8080` instead.

## API

### `POST` /api/auth/register

This endpoint accepts the following request payload:

```json
{
  "name": "username",
  "email": "user@example.com",
  "password": "securePassword@123",
  "confirm_password": "securePassword@123",
}
```

It creates a new user in the database, notably with field `verified = false`. An email will be sent to the address specified containing a JWT-embedded URL that will call the `/api/auth/verify` route to set `verified = true`. 

The verification URL expires 10 minutes after it was created, and there is a scheduled cleanup of unverified users that triggers every 24 hours. 

### `POST` /api/auth/verify

This endpoint sets the `verified` field in created users to `true`. It accepts a `?token=<token>` query parameter, which would have been embedded to the URL sent by the `/api/auth/register` endpoint to the user's email address.

### `POST` /api/auth/login

This endpoint accepts the following request payload:

```json
{
  "email": "user@example.com",
  "password": "securePassword@123"
}
```

It checks for an existing record of the user in the database, and issues access and refresh tokens as `Httponly`, `Secure`, `Strictly samesite` cookies. 

### `POST` /api/auth/logout **(PROTECTED)**

**PROTECTED RESOURCE:** This endpoint is a protected resource. A middleware expects `access_token` and `refresh_token` cookies with every request. 

This endpoint is responsible for (1) invalidating these tokens by altering their expiration date to be in the past, and (2) revoking these tokens so they can no longer be used for authorization.

### `GET` /api/auth/decode **(PROTECTED)**

**PROTECTED RESOURCE:** This endpoint is a protected resource. A middleware expects `access_token` and `refresh_token` cookies with every request. 

This endpoint simply decodes the access token and returns the claims to be used by the frontend - typically for managing user state.

### `POST` /api/auth/2fa/email/request

This endpoint sends a verification URL to the user's email address, embedded with a JWT that expires in 10 minutes.

### `POST` /api/auth/2fa/email/verify

This endpoint expects a `?token=<token>` query parameter. If the token is valid, then access and refresh tokens will be issued as `Httponly`, `Secure`, and `Strictly Samesite` cookies.
