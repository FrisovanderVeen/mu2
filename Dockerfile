FROM golang:alpine

RUN apk add --update git
RUN go get -u github.com/golang/dep/cmd/dep

ADD . /go/src/github.com/FrisovanderVeen/mu2
WORKDIR /go/src/github.com/FrisovanderVeen/mu2
RUN dep ensure
RUN go build

FROM golang:alpine

RUN apk add --update ffmpeg
COPY --from=0 /go/src/github.com/FrisovanderVeen/mu2/mu2 /bot/mu2

CMD [ "/bot/mu2" ]
