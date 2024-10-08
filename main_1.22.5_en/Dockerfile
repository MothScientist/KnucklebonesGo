# syntax=docker/dockerfile:1

FROM golang:1.22.5

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-knucklebones-go

EXPOSE 8080

CMD ["/docker-knucklebones-go"]