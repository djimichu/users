FROM golang:1.24-alpine

# Установим goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Добавим сертификаты и bash (если нужно)
RUN apk add --no-cache ca-certificates bash

# Создаем рабочую директорию и копируем миграции
WORKDIR /app
COPY internal/repository/migrations ./migrations

# Устанавливаем переменные по умолчанию
ENV PATH="/go/bin:$PATH"

ENTRYPOINT ["goose"]