FROM golang:latest AS build

RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/fvdveen/mu2
ADD . .

RUN dep init
RUN CGO_ENABLED=0 go build -o mu2 main.go

FROM alpine:latest AS RUN

WORKDIR /app/
COPY --from=build /go/src/github.com/fvdveen/mu2/mu2 mu2
RUN apk add --no-cache ca-certificates

CMD ["./mu2"]
