# Этап сборки
FROM golang:1.25.1 AS builder


# Рабочая директория
WORKDIR /app

# Копируем go.mod и go.sum отдельно для кэширования зависимостей
COPY ../../go.mod ../../go.sum ./
RUN go mod download



# копируем исходный код
COPY ../../ .

# собираем после копирования
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Финальный минимальный образ
FROM alpine:3.20

WORKDIR /app

# Копируем бинарник
COPY --from=builder /app/server .

# порт
EXPOSE 9091

# Запускаем миграции и сервер
CMD ["./server"]