FROM golang:alpine3.19 as build

RUN apk update; apk add git

WORKDIR /workdir
COPY build/compile.sh /workdir
RUN /workdir/compile.sh

FROM alpine:3.19.1

ENV PATH=/usr/local/bin:$PATH
COPY --from=build /dist/* /usr/local/bin/
