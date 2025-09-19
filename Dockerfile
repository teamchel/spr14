FROM golang:1.25-alpine

WORKDIR /app

# Устанавливаем git
RUN apk add --no-cache git

# Копируем модули
COPY go.mod .
COPY go.sum .

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o main .

# Запуск
CMD ["./main"]