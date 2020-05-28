# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.13 as build-env

ENV GO111MODULE=on

RUN  mkdir -p /go/src \
  && mkdir -p /go/bin \
  && mkdir -p /go/pkg
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH   

# now copy your app to the proper build path
RUN mkdir -p $GOPATH/src/github.com/KyberNetwork/cache 
ADD . $GOPATH/src/github.com/KyberNetwork/cache

WORKDIR $GOPATH/src/github.com/KyberNetwork/cache 
RUN make cache

FROM debian:stretch
RUN apt-get update && \
    apt install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*
COPY --from=build-env /go/src/github.com/KyberNetwork/cache/build/bin/cache /cache
COPY --from=build-env /go/src/github.com/KyberNetwork/cache/env/ /env/

EXPOSE 3001
ENV GIN_MODE release
ENV KYBER_ENV production
ENV LOG_TO_STDOUT true

CMD ["/cache"]
