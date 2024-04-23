# Файлы для итогового задания

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.

# Локальный запуск приложения
 export TODO_PORT=7540
 go run main.go

 # Запуск тестов
 1. go test -run ^TestApp$ ./tests
 2. go test -run ^TestDB$ ./tests
 3. go test -run ^TestNextDate$ ./tests
 4. go test -run ^TestAddTask$ ./tests
 5. go test -run ^TestTasks$ ./tests
 6. go test -run ^TestTask$ ./tests
 7. go test -run ^TestEditTask$ ./tests
 8. go test -run ^TestDone$ ./tests
 9. go test -run ^TestDelTask$ ./tests
 10. go test ./tests