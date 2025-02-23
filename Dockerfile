# Используем базовый образ для Go
FROM golang:1.20 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем все файлы в контейнер
COPY . .

# Компилируем Go-приложение
RUN go mod tidy && go build -o server .

# Финальный образ для запуска
FROM ubuntu:latest

# Устанавливаем рабочую папку
WORKDIR /root/

# Копируем собранный бинарник из предыдущего шага
COPY --from=builder /app/server .

# Запускаем сервер
CMD ["./server"]