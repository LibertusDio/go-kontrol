FROM golang:1.18-alpine AS builder
RUN apk add --no-cache make \
&& rm -vrf /var/cache/apk/*



WORKDIR /go/src/app
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o server ./server/


FROM alpine:3.14
RUN apk add --no-cache --update ca-certificates tzdata


COPY --from=builder /go/src/app/server /app/server
WORKDIR /app
CMD ["/app/server"]