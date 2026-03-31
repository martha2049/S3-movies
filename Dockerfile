FROM golang:1.24 AS build

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o /bin/s3-server cmd/main.go

CMD ["/bin/s3-server"]