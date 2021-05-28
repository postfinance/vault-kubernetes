FROM alpine:latest as runtime
MAINTAINER OpenSource PF <opensource@postfinance.ch>

RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates
COPY token-renewer /token-renewer

# Run as nobody:x:65534:65534:nobody:/:/sbin/nologin
USER 65534

CMD ["/token-renewer"]
