FROM golang:1.18

ENV GOPATH=

COPY go.mod go.sum ./
CMD go mod tidy

COPY ./app/main.go .
CMD go run main.go