FROM golang:1.12-alpine

RUN set -ex; \
    apk update; \
    apk add --no-cache git

RUN go get -u "cloud.google.com/go/firestore"
RUN go get -u "firebase.google.com/go"