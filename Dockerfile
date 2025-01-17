FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod tidy

COPY . .

RUN go build -o server cmd/main.go

EXPOSE 8080

CMD ["./server"]

