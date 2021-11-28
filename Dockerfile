# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /bot

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /tg-bot

CMD [ "/tg-bot" ]
