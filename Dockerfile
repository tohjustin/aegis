FROM golang:alpine
ADD . /go/src/github.com/tohjustin/badger
RUN go install github.com/tohjustin/badger
CMD ["/go/bin/badger"]
EXPOSE 8080
