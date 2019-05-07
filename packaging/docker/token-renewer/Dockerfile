FROM golang:1.12-alpine as builder
RUN apk --no-cache add git

ENV GO111MODULE=on
ENV CGO_ENABLED=0

WORKDIR /vgo/
COPY ./cmd/token-renewer/ .

RUN go build -o /token-renewer main.go

# Build runtime
FROM alpine:3.8 as runtime
MAINTAINER OpenSource PF <opensource@postfinance.ch>

RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /token-renewer /token-renewer

# Run as nobody:x:65534:65534:nobody:/:/sbin/nologin
USER 65534

CMD ["/token-renewer"]
