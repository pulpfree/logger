FROM gliderlabs/alpine

RUN apk update && apk add ca-certificates && apk add curl

ADD logger /go/src/
COPY config/*.yaml /go/src/config/
# COPY ssl/* /go/src/ssl/

# ENV environment="dev"
# ENV environment="prod"

# HEALTHCHECK --interval=1m --timeout=3s \
# CMD curl -f http://localhost:3003/healthcheck || exit 1

WORKDIR /go/src

ENTRYPOINT /go/src/logger

EXPOSE 3021