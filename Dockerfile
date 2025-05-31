FROM golang:1.24-alpine
WORKDIR /app
COPY . .
RUN go mod init ruayAutoMsg || true
RUN go mod tidy
RUN go build -o main .
CMD ["./main"]
