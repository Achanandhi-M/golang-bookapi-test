FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o book-tracker ./cmd/server

EXPOSE 8080

CMD ["./book-tracker"]