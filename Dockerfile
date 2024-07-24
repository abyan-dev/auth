FROM alpine:3.20

WORKDIR /app

COPY ./auth /app

CMD ["/app/auth"]
