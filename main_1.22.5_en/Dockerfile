# syntax=docker/dockerfile:1

FROM golang:1.22.5-bookworm

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY *.go ./

RUN go build -o docker-knucklebones-go

EXPOSE 8080

CMD ["/app/docker-knucklebones-go"]
