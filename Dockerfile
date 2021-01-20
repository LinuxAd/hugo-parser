FROM golang:1.15.7-alpine3.13

WORKDIR /src
COPY . .
RUN go build -o /hugo-parser

ENTRYPOINT /hugo-parser