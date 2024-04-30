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

# Работа с docker

## Сборка docker образа
`docker build -t gt814/go-final-project:v1.0.0 .`

## Запуск docker контейнера
`docker run -p 7540:7540 gt814/go-final-project:v1.0.0`

## Добавление тега версии для активации процесса пуша образа в docker registry
`git tag v1.0.0 git push --tags`

## Запуск docker конетейнера из DockerHub
```
docker pull --platform linux/x86_64 gt814/go-final-project:v1.0.0 
docker run -p 7540:7540 --platform linux/x86_64 gt814/go-final-project:v1.0.0
```