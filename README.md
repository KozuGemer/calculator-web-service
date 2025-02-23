# Calculator Web Service

Этот проект представляет собой веб-сервис для вычислений, который принимает математические выражения через API и возвращает результат. Также доступен обычный калькулятор через веб-интерфейс.

## Установка



### Windows
Скачайте .exe из Releases. Запустите .exe програаму, когда сервер запустится, на экране на пишетсятся Server running :8080. Перейдите по [localhost:8080](http://localhost:8080) И начинайте возится обычным калькулятором.

### Другие системы
Для других систем вам нужно установить через Docker или из исходников. Как установить рассказано на [Wiki:Starting](https://github.com//KozuGemer/calculator-web-service/wiki/Starting)

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
Дополнительная документация, на английском языке.

Полное руководство можно найти в [Wiki проекта](https://github.com//KozuGemer/calculator-web-service/wiki).

