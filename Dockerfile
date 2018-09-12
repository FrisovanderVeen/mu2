FROM golang:1.11rc2 AS build

WORKDIR /mu2

ADD go.mod go.sum ./
RUN go mod download

ADD . .

RUN CGO_ENABLED=0 go build -o mu2 main.go

FROM alpine:latest AS RUN

WORKDIR /app/
COPY --from=build /mu2/mu2 mu2
RUN apk add --no-cache ca-certificates

CMD ["./mu2"]
