# GO PARSER

Сервис для парсинга данных с использованием:
- Playwright для веб-скрапинга
- RabbitMQ для очереди задач
- MongoDB для хранения результатов

## Требования

- Go 1.23+
- MongoDB
- RabbitMQ
- Node.js (для Playwright)

## Install

```shell
go run github.com/playwright-community/playwright-go/cmd/playwright install
```

## Структура проекта

```bash
.
├── cmd/app/ # Точка входа в приложение
├── internal/ # Внутренний код приложения
│ ├── config/ # Конфигурация
│ ├── handlers/ # Обработчики
│ ├── models/ # Модели данных
│ ├── services/ # Сервисы для работы с внешними системами
│ └── utils/ # Вспомогательные функции
└── ...
```

## Запуск

1. Скопируйте `.env.example` в `.env` и настройте переменные окружения
2. `go run cmd/app/main.go`
