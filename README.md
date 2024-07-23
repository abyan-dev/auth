# Auth v1

Authentication, authorization, and account management service that exposes a REST API for issuing access and refresh as JSON Web Tokens (JWT) stored in `httponly` and `strict` samesite cookies.

## Running the service locally

Initialize `.env` by running:

```
cp .env.default .env
```

Make sure to populate the `JWT_SECRET` environment variable. To generate a random key, use:

```
openssl rand -base64 64
```

The service uses postgres with GORM to store user credentials. It can also use mailhog on port `8025` to simulate the sending of verification emails locally. To start both postgres and mailhog as containers, run:

```
make db-up
mail mailhog-up
```

You can stop both services by running `make db-down` and `make mailhog-down` when you're done.

Finally, run the server:

```
make run
```

There is a `GET` health check to verify the server is up and running at `/api/health`, as well as a protected variant at `/api/health/protected` for testing the authorization middleware.
