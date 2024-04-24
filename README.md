# Файлы для итогового задания

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.

# Локальный запуск приложения
 export TODO_PORT=7540
 go run main.go

# Запуск тестов
`go test ./tests`

## Запуск конкретных методов тестов
`go test -run ^TestName$ ./tests`, где `TestName` - имя теста.
Например, `go test -run ^TestApp$ ./tests`.
