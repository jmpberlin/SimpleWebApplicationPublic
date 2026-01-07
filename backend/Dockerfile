FROM golang:1.25.0-alpine

WORKDIR /usr/app

COPY . .

RUN go build -o server main.go

CMD ["./server"]