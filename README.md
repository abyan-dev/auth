# Auth

![Tests](https://github.com/abyan-dev/auth/actions/workflows/ci.yaml/badge.svg) [![codecov](https://codecov.io/gh/abyan-dev/auth/graph/badge.svg?token=S679A5TSW7)](https://codecov.io/gh/abyan-dev/auth) [![Go Report](https://goreportcard.com/badge/abyan-dev/auth)](https://goreportcard.com/report/YanSystems/compiler) [![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/YanSystems/compiler/blob/main/LICENSE)

Authentication, authorization, and account management service that exposes a REST API for issuing, invalidating, and revoking JSON Web Tokens (JWT). Tokens are stored in `Httponly`, `Secure`, and `strict` samesite cookies, thereby minimizing XSS and CSRF vulnerabilities. 

There is additionally a scheduling of the deletion of unverified accounts, which runs concurrently on a separate goroutine.

## Integrating this service

To integrate this service to your application, build a linux executable and use it to create a lightweight docker image on alpine OS:

```
env GOOS=linux CGO_ENABLED=0 go build -o auth ./cmd/api && make image
```

Then, serve the container (make sure postgres is running):

```
docker run -p 8080:8080 --name auth auth
```

Your application can interact with the service at `localhost:8080`. Alternatively, if you are using an orchestrator like Compose or Kubernetes, you can leverage their DNS by using `auth:8080` instead.

**Mails in Production:** In production, remember to configure the environment variables for email sending with actual delivery services like Mailgun.