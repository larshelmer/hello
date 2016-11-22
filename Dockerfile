FROM golang
ADD . /go/src/github.com/larshelmer/hello
RUN go install github.com/larshelmer/hello
ENTRYPOINT /go/bin/hello
EXPOSE 8080
