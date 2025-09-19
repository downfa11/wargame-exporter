# 1 Phase
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o metrics-exporter main.go

# 2 Phase
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/metrics-exporter .

VOLUME ["/mnt/metrics"]

EXPOSE 9100

ENV EXPORTER_PORT=9100
ENV METRICS_FILE=/mnt/metrics/metrics_snapshot.json

ENTRYPOINT ["./metrics-exporter"]