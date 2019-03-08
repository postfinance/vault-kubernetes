FROM golang:1.12-alpine as builder
RUN apk --no-cache add git

ENV GO111MODULE=on
ENV CGO_ENABLED=0

WORKDIR /vgo/
COPY . .

RUN ls -al
RUN go build -o /synchronizer ./cmd/synchronizer/main.go

# Build runtime
FROM alpine:3.8 as runtime
MAINTAINER OpenSource PF <opensource@postfinance.ch>

COPY --from=builder /synchronizer /synchronizer

# Run as nobody:x:65534:65534:nobody:/:/sbin/nologin
USER 65534

CMD ["/synchronizer"]

