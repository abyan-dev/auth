# Auth

![Tests](https://github.com/abyan-dev/auth/actions/workflows/ci.yaml/badge.svg) [![codecov](https://codecov.io/gh/abyan-dev/auth/graph/badge.svg?token=S679A5TSW7)](https://codecov.io/gh/abyan-dev/auth) [![Go Report](https://goreportcard.com/badge/abyan-dev/auth)](https://goreportcard.com/report/YanSystems/compiler) [![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/YanSystems/compiler/blob/main/LICENSE)

Authentication, authorization, and account management service that exposes a REST API for issuing access and refresh JSON Web Tokens. Tokens are stored in `httponly` and `strict` samesite cookies, thereby protected from XSS and CSRF attacks, and are revoked on logout. 

There is additionally a scheduling of the deletion of unverified accounts, which runs concurrently on a separate goroutine.

## Running the service locally

Make sure to populate the `JWT_SECRET` environment variable in `.env.default`. To generate a random key, use:

```
openssl rand -base64 64
```

Now, load the environment variables to your shell session:

```
cp .env.default .env && source .env
```

The service uses postgres with GORM to store user credentials. It can also use mailhog on port `8025` to simulate the sending of verification emails locally. To start both postgres and mailhog as containers, run:

```
make db-up
make mailhog-up
```

You can stop both services by running `make db-down` and `make mailhog-down` when you're done.

Finally, run the server:

```
make run
```

There is a `GET` health check to verify the server is up and running at `/api/health`, as well as a protected variant at `/api/health/protected` for testing the authorization middleware.

## Integration

To integrate this service to your application, build a linux executable and use it to create a lightweight docker image on alpine OS:

```
env GOOS=linux CGO_ENABLED=0 go build -o auth ./cmd/api && make image
```

Then, serve the container (make sure postgres is running):

```
docker run -p 8080:8080 --name auth auth
```

Your application can interact with the service at `http://localhost:8080`. If you use an orchestration tool like Docker Compose or Kubernetes, you can leverage their DNS to simplify routing by referring to the container's name i.e., "auth" at `http://auth:8080` and you can remove the port flag.

**PRODUCTION:** In production, remember to configure the environment variables for email sending with actual delivery services like Mailgun.
