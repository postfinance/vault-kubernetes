FROM gcr.io/distroless/static-debian12:latest as runtime
LABEL maintainer="OpenSource PF <opensource@postfinance.ch>"

COPY authenticator /authenticator

# Run as nobody:x:65534:65534:nobody:/:/sbin/nologin
USER 65534

CMD ["/authenticator"]
