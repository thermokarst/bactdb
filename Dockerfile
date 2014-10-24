FROM golang:1.3.1
ADD . /go/src/github.com/thermokarst/bactdb
RUN go get -d -v github.com/thermokarst/bactdb/cmd/bactdb
RUN go install github.com/thermokarst/bactdb/cmd/bactdb
ENTRYPOINT /go/bin/bactdb serve
EXPOSE 8901

