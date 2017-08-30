FROM resin/raspberrypi3-golang

WORKDIR /go/src/github.com/qsz13/ooxxbot

COPY . .

RUN go get golang.org/x/net/html &&  go install

CMD ["/go/bin/ooxxbot"]




