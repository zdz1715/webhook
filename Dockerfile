FROM alpine:3.16

WORKDIR /user/local/webhook

RUN apk add --no-cache ca-certificates

ADD bin/webhook bin/webhook

ENTRYPOINT ["bin/webhook"]