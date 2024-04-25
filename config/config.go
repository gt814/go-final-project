package config

import (
	"log"
	"os"
	"strconv"
)

func GetPort() string {
	todoPort := os.Getenv("TODO_PORT")
	if len(todoPort) > 0 {
		intTodoPort, err := strconv.Atoi(todoPort)
		if err != nil {
			log.Fatalln(err)
		}

		port = intTodoPort
	}
	return strconv.Itoa(port)
}

func GetDBFileAppPath() string {
	todoDbFile := os.Getenv("TODO_DBFILE")

	if len(todoDbFile) > 0 {
		return todoDbFile
	}
	return dbFile
}

func GetTaskLimit() int {
	todoTaskLimit := os.Getenv("TODO_TASKLIMIT")

	if len(todoTaskLimit) > 0 {
		intTodoTaskLimit, err := strconv.Atoi(todoTaskLimit)
		if err != nil {
			log.Fatalln(err)
		}
		taskLimit = intTodoTaskLimit
	}

	return taskLimit
}
