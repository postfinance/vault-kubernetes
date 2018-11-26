FROM alpine:3.8

LABEL authors="marc.sauter@postfinance.ch"

ARG BINARY
COPY $BINARY /authenticator

ENTRYPOINT ["/authenticator"]

