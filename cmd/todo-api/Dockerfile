##############
# Base image #
##############
FROM golang:1.26.3-alpine3.23 AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN go mod tidy
RUN go build -o todo-api ./cmd/todo-api/main.go

#################
# Runtime image #
#################
FROM alpine:3.23
WORKDIR /app
COPY --from=builder /app/todo-api .
CMD ["./todo-api"]
