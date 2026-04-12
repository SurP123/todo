FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o main .

FROM fedora:39

WORKDIR /app

RUN dnf install -y sqlite && dnf clean all

COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates

RUN mkdir -p /app/data

EXPOSE 8181

CMD ["./main"]