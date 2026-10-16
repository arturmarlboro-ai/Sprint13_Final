FROM golang:1.26 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM ubuntu:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 7540
CMD ["./main"]
