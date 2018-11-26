FROM alpine:3.8

LABEL authors="marc.sauter@postfinance.ch"

ARG BINARY
ADD $BINARY /synchronizer

ENTRYPOINT ["/synchronizer"]

