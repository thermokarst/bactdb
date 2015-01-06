FROM golang:1.4.0
ADD . /go/src/github.com/thermokarst/bactdb
RUN go get -d -v github.com/thermokarst/bactdb/cmd/bactdb
RUN go install github.com/thermokarst/bactdb/cmd/bactdb
CMD /go/bin/bactdb serve --keys=/bactdb/keys/
EXPOSE 8901

