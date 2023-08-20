# Build stage
FROM golang:1.20.7-alpine3.17 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz


# Run stage
FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/main .
# COPY --from=builder /app/migrate.linux-amd64 ./migrate
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY db/migrations ./migrations
COPY start.sh .
COPY wait-for.sh .

EXPOSE 8080

CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]