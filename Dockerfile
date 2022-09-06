FROM zdzserver/runtime:alpine3.16

WORKDIR /user/local/webhook

RUN apk add --no-cache ca-certificates

ADD bin/webhook bin/webhook

ENTRYPOINT ["bin/webhook"]