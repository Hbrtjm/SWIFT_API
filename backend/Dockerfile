FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

# DEVELOPMENT 

# CMD ["go", "run", "./cmd/api/main.go"]

# DEPLOYMENT

RUN go build -o main ./cmd/api/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs

RUN mkdir /app/logs
RUN chmod +x ./main

CMD ["./main"]