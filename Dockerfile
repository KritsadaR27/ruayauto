FROM golang:1.24-alpine
WORKDIR /app
COPY backend/ .
RUN go mod tidy
RUN go build -o main ./cmd/server
CMD ["./main"]
