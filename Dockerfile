FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server .
COPY uploads/script.txt uploads/script.txt
EXPOSE 3000
CMD ["./server"]
