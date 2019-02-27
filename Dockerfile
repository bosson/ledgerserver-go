FROM alpine

ARG PORT=8080

RUN apk --no-cache add sqlite-libs ca-certificates

ADD ledgerserver-go-linux-amd64 /ledgerserver

EXPOSE ${PORT}

ENTRYPOINT ["/ledgerserver"]
