# ---------- build stage ----------
FROM golang:1.25.1 AS builder

WORKDIR /app

# копируем зависимости
COPY ../../go.mod ../../go.sum ./
RUN go mod download

# копируем исходный код
COPY ../../ .

# Устанавливаем goose и собираем приложение
RUN go install github.com/pressly/goose/v3/cmd/goose@latest && \
    CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# ---------- runtime stage ----------
FROM alpine:3.20

WORKDIR /app

# копируем бинарник
# Копируем ТОЛЬКО нужные файлы
COPY --from=builder /app/internal/repository/migrations ./migrations
COPY --from=builder /app/server .
COPY --from=builder /go/bin/goose /usr/local/bin/goose

#  В Alpine Linux goose нужно устанавливать отдельно
RUN wget -O /usr/local/bin/goose https://github.com/pressly/goose/releases/download/v3.20.0/goose_linux_x86_64 && \
    chmod +x /usr/local/bin/goose

# твой сервер слушает 9091
EXPOSE 9091

# Запускаем миграции и сервер
CMD ["sh", "-c", "goose -dir /app/migrations postgres \"$DB_URL\" up && ./server"]
