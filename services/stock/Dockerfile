FROM golang:1.15-alpine

WORKDIR /build

COPY . .

RUN go build -o main .

WORKDIR /app

RUN cp /build/main .

EXPOSE 8080

CMD ["/app/main"]