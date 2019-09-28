FROM golang:alpine AS build

RUN apk add --no-cache bash git

WORKDIR /go/src/app

ADD . ./

RUN DEV=true scripts/build.sh

FROM alpine 

COPY --from=build /go/src/app/build/bin/letitgo /letitgo
ENTRYPOINT [ "/letitgo" ]
