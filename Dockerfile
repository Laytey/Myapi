# 1 Базовый образ для сборки, builder - дать имя этому этапу
FROM golang:1.26 AS builder
# 2 WORKDIR	Установить рабочую директорию
# WORKDIR /app - Рабочая папка внутри контейнера
# 
WORKDIR /app
# 3 Переменные окружения
ENV CGO_ENABLED=0
#Отключить CGO
ENV GOOS=linux
#Собрать бинарник под Linux
# 4 Копируем файлы зависимостей
COPY go.mod go.sum ./
# 5 Скачиваем зависимости
RUN go mod download
# 6 Копируем весь код
COPY . .
# 7 Собираем бинарник
RUN go build -o myapi cmd/main.go
# -o myapi	Назвать бинарник myapi
# cmd/main.go	Точка входа
# в /app/myapi лежит готовый бинарник

# 8 Финальный образ
FROM alpine:latest
#FROM alpine:latest	Начать новый этап с лёгкого образа Alpine Linux
#alpine	Минимальный Linux (~5 МБ)
# 9 Рабочая папка
WORKDIR /root
# 10 Копируем бинарник из builder
COPY --from=builder /app/myapi .
# 11 Открываем порт
EXPOSE 8080
# 12 Команда запуска
ENTRYPOINT ["./myapi"]