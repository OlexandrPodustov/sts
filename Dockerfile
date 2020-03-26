FROM golang:1.14.1
WORKDIR /go/src/github.com/OlexandrPodustov/sts
COPY . .

RUN go install ./cmd

ENV MONGO_ADDRESS=172.17.0.2

EXPOSE 8080