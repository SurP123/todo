FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o main .

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates


RUN chmod -R 755 /app/templates

EXPOSE 8181

CMD ["./main"]