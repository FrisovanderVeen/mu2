FROM golang:alpine

RUN apk add --update ffmpeg
RUN apk add --update git
RUN go get -u github.com/golang/dep/cmd/dep

ADD . /go/src/github.com/FrisovanderVeen/mu2
WORKDIR /go/src/github.com/FrisovanderVeen/mu2
RUN dep ensure
RUN go install

CMD [ "mu2" ]
