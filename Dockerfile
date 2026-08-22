# ======= DEVELOPMENT STAGE =======
# Hot-reload için air yüklenir, local'de kullanılır
FROM golang:1.25-alpine AS dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/air-verse/air@latest

COPY . .

CMD ["air"]

# ======= BUILD STAGE =======
# Go kodu derlenir, tek binary üretilir
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /city-pulse .

# ======= PRODUCTION STAGE =======
# Sadece binary ve config kopyalanır, çok küçük bir imaj olur
FROM alpine:3.21 AS production

WORKDIR /app

RUN apk add --no-cache ca-certificates wget tzdata

ENV TZ=Europe/Istanbul

COPY --from=builder /city-pulse /app/city-pulse
COPY config/prod-config.json /app/config/prod-config.json
COPY ui/ /app/ui/


EXPOSE 3649

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO- http://127.0.0.1:3649/api/1.0/health || exit 1

CMD ["/app/city-pulse"]
