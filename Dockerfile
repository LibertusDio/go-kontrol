FROM golang:1.18-alpine AS builder
RUN apk add --no-cache --update bash git
WORKDIR /go/src/app
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
COPY . .
RUN go build -o execute ./server/

FROM alpine:3.14
RUN apk add --no-cache --update ca-certificates tzdata curl
COPY --from=builder /go/src/app/execute /execute
WORKDIR /
ENTRYPOINT ["/execute"]