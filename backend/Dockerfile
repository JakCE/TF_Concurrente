# backend/Dockerfile

FROM golang:1.24.2

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o main

CMD ["./main"]