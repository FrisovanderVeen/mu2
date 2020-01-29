FROM golang:latest AS build

WORKDIR /mu2

ADD go.mod go.sum ./
RUN go mod download

ADD ./bot ./bot
ADD ./cmd ./cmd
ADD ./commands ./commands
ADD ./common ./common
ADD ./config ./config
ADD ./voice ./voice

RUN go build -ldflags "-linkmode external -extldflags -static" -a -o mu2 cmd/mu2/main.go

FROM alpine:latest AS RUN

WORKDIR /app/
COPY --from=build /mu2/mu2 mu2
RUN apk add --no-cache ca-certificates ffmpeg

ADD config.json config.json

CMD ["./mu2"]
