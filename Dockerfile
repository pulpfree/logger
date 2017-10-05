FROM gliderlabs/alpine

RUN apk update && apk add ca-certificates && apk add curl

ADD logger /go/src/

COPY config/*.yaml /go/src/config/

WORKDIR /go/src

ENTRYPOINT /go/src/logger

# NOTE: Dont't really need this if we're using docker network
EXPOSE 3021