# Stage 1: Build
FROM golang:1.26-alpine AS builder

# Устанавливаем необходимые инструменты для сборки (если нужны)
RUN apk add --no-cache git

WORKDIR /src

# Сначала копируем файлы зависимостей для эффективного кэширования слоев Docker
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Сборка:
# -s -w убирают отладочную информацию (минус пару мегабайт)
# CGO_ENABLED=0 делает бинарник статическим (не требует библиотек в Alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /fleetmaster ./cmd/app/main.go

# Stage 2: Final
FROM alpine:3.23

# Добавляем сертификаты для работы с HTTPS и создаем пользователя
RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -g '' appuser

WORKDIR /app

# Копируем ТОЛЬКО бинарный файл из образа builder
COPY --from=builder --chown=appuser:appuser /fleetmaster .

# Переключаемся на обычного пользователя
USER appuser

EXPOSE ${APP_PORT}

# Запускаем через ./ для надежности
CMD ["./fleetmaster"]
