FROM gcr.io/distroless/static-debian11:latest as runtime
MAINTAINER OpenSource PF <opensource@postfinance.ch>

COPY token-renewer /token-renewer

# Run as nobody:x:65534:65534:nobody:/:/sbin/nologin
USER 65534

CMD ["/token-renewer"]
