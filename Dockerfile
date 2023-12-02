FROM golang:1.21-alpine as builder

WORKDIR /app
COPY app ./
RUN go mod download

RUN go build -o main .
CMD ["./main"]
