FROM golang:1.25-alpine AS builder
ENV GO111MODULE=on
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o scheduler .
FROM alpine:latest
LABEL maintainer="Your Name <your.email@example.com>"
LABEL description="Планировщик задач TODO на Go"
WORKDIR /app
COPY --from=builder /app/scheduler .
COPY --from=builder /app/web ./web
EXPOSE 7540
ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/scheduler.db
CMD ["./scheduler"]