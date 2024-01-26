FROM golang:1.21.4

WORKDIR /app

COPY . .

RUN go build -o main .

EXPOSE 8008

CMD ["./main"]