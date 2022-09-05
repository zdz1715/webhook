FROM alpine:3.16


RUN apk add --no-cache ca-certificates

ADD bin/webhook /user/local/webhook

ENTRYPOINT ["/user/local/webhook"]