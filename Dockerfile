FROM golang:1.8
WORKDIR /go/src/sts
COPY . .

RUN go install -v ./...

ENV MONGO_ADDRESS=172.17.0.2

EXPOSE 8081