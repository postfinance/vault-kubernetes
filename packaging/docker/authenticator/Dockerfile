FROM golang:1.15.7-alpine as builder
RUN apk --no-cache add git

ENV GO111MODULE=on
ENV CGO_ENABLED=0

WORKDIR /vgo/
COPY ./cmd/authenticator/ .

RUN go build -o /authenticator main.go

# Build runtime
FROM alpine:latest as runtime
MAINTAINER OpenSource PF <opensource@postfinance.ch>

RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /authenticator /authenticator

# Run as nobody:x:65534:65534:nobody:/:/sbin/nologin
USER 65534

CMD ["/authenticator"]
