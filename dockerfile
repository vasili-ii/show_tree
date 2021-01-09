# docker build -t mailgo_hw1 .
FROM golang:latest
COPY . .
RUN go test -v