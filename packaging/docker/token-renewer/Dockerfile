FROM alpine:3.8

LABEL authors="marc.sauter@postfinance.ch"

ARG BINARY
ADD $BINARY /token-renewer

ENTRYPOINT ["/token-renewer"]

