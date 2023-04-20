FROM golang:1.20.3-alpine3.17 AS builder

ARG package

WORKDIR /src
COPY . .
#COPY --from=base /src/srv/cmd/addressbook/config/config.json ./config/

RUN apk add --update --no-cache git
RUN go build -o srv ./cmd/addressbook

FROM alpine:latest

COPY --from=builder /src/srv /usr/local/bin/
COPY --from=builder /src/cmd/addressbook/config/config.json /usr/local/bin/


EXPOSE 8080

CMD ["srv"]