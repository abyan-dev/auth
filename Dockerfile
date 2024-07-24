FROM golang:1.22.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/app ./cmd/api

FROM alpine:3.20

COPY --from=builder /app/app /app/app

WORKDIR /app

CMD ["/app/app"]

