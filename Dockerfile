# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang
EXPOSE 3001

ENV GIN_MODE=release

RUN  mkdir -p /go/src \
  && mkdir -p /go/bin \
  && mkdir -p /go/pkg
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH   

# now copy your app to the proper build path
RUN mkdir -p $GOPATH/src/github.com/KyberNetwork/server-go 
ADD . $GOPATH/src/github.com/KyberNetwork/server-go

# should be able to build now
WORKDIR $GOPATH/src/github.com/KyberNetwork/server-go 
RUN go build -o server . 
CMD ["/go/src/github.com/KyberNetwork/server-go/server"]



