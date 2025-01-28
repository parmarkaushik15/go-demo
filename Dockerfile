
FROM golang:1.23.5 AS builder
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
ENV REDIS_URL=localhost:6379
ENV RABBITMQ_URL=amqp://guest:guest@localhost:5672/
ENV DB_HOST=localhost
ENV DB_PORT=5432
ENV DB_USER=user
ENV DB_PASSWORD=password
ENV DB_NAME=mydb
ENV CSV_PATH="/root/go-demo-api/users.csv"
WORKDIR /go/src/go-demo-api
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
COPY .env .env
RUN go build -o /go-demo-api
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go-demo-api .
EXPOSE 8080
CMD ["/root/go-demo-api"]
