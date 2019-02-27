FROM alpine

ARG PORT=8080

RUN apk --no-cache add sqlite-libs ca-certificates

ADD ledgerserver-linux-amd64 /ledger-server

EXPOSE ${PORT}

ENTRYPOINT ["/ledger-server"]
