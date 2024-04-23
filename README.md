# Файлы для итогового задания

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.

# Локальный запуск приложения
 export TODO_PORT=7540
 go run main.go

 # Запуск тестов
 go test -run ^TestApp$ ./tests
 go test -run ^TestDB$ ./tests
 go test -run ^TestNextDate$ ./tests
 go test -run ^TestAddTask$ ./tests
 go test -run ^TestTasks$ ./tests
 go test -run ^TestTask$ ./tests
 go test -run ^TestEditTask$ ./tests
 go test -run ^TestDone$ ./tests