FROM golang:1.22.5-bookworm

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY *.go ./

RUN go build -o docker-knucklebones-go

RUN rm -f *.go && rm -f go.mod && rm -f go.sum

EXPOSE 8080
