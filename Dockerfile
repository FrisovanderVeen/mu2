FROM golang:alpine

ADD . /go/src/github.com/fvdveen/mu2
WORKDIR /go/src/github.com/fvdveen/mu2
RUN go build

FROM golang:alpine

RUN apk add --update ffmpeg
COPY --from=0 /go/src/github.com/fvdveen/mu2/mu2 /go/bin

CMD [ "mu2", "run" ]
