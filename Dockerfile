# build stage
FROM golang:alpine AS builder
# RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/app/our-expenses-be -v ./cmd/server
COPY config/config.yaml /go/bin/app/config/config.yaml
RUN mkdir -p /go/bin/app/storage/logs

# final stage
FROM alpine:latest
# RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /go/bin/app/ .
LABEL Name=our-expenses-be Version=0.0.1
EXPOSE 8080
ENTRYPOINT /app/our-expenses-be
