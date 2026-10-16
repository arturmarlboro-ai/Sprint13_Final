FROM golang:1.26 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM ubuntu:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/web ./web
EXPOSE 7540
CMD ["./main"]
