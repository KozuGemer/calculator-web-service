# Calculator Web Service

Этот проект представляет собой веб-сервис для вычислений, который принимает математические выражения через API и возвращает результат. Также доступен обычный калькулятор через веб-интерфейс.

## Установка

Библиотека jq позволяет лучше читать json код.

### Windows

```
winget install jqlang.jq
```

Работает, если у вас выше Windows 10

### Linux (Debian/Ubuntu)

```sh
sudo apt update && sudo apt install -y golang jq
```

### macOS

```sh
brew install go jq
```

## Запуск сервера

```sh
go run main.go
```

После запуска сервер работает на `http://localhost:8080`.

## Использование

### Обычный калькулятор

Перейдите в браузере по адресу `http://localhost:8080` и введите выражение в поле ввода.

### API

#### 1. Создание задачи

```sh
curl -X POST http://localhost:8080/api/v1/tasks -H "Content-Type: application/json" -d '{"expression": "2+2"}' | jq
```

Ответ:

```json
{
  "id": "task-123456",
  "expression": "2+2",
  "status": "201 - Accepted for Processing",
  "message": "Task has been created and is being processed"
}
```

#### 2. Получение статуса задачи

```sh
curl -X GET "http://localhost:8080/api/v1/tasks/status?id=task-123456" | jq
```

Ответ:

```json
{
  "id": "task-123456",
  "expression": "2+2",
  "result": 4,
  "status": "done"
}
```

#### 3. Получение всех задач

```sh
curl -X GET http://localhost:8080/api/v1/expressions | jq
```

## Дополнительная документация

Полное руководство можно найти в [Wiki проекта](https://github.com/ТВОЙ-АККАУНТ/РЕПОЗИТОРИЙ/wiki).

