# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.13 as build-env
EXPOSE 3001
ENV GO111MODULE=on

RUN  mkdir -p /go/src \
  && mkdir -p /go/bin \
  && mkdir -p /go/pkg
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH   

# now copy your app to the proper build path
RUN mkdir -p $GOPATH/src/github.com/KyberNetwork/cache 
ADD . $GOPATH/src/github.com/KyberNetwork/cache

# should be able to build now
WORKDIR $GOPATH/src/github.com/KyberNetwork/cache 
# RUN go mod vendor
RUN make cache
CMD ["/go/src/github.com/KyberNetwork/cache/build/bin/cache"]



