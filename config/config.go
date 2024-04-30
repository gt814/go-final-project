package config

import (
	"log"
	"os"
	"strconv"
)

const port = 7540
const dbFile = "./scheduler.db"
const taskLimit = 50

func GetPort() string {
	todoPort := os.Getenv("TODO_PORT")
	if len(todoPort) > 0 {
		intTodoPort, err := strconv.Atoi(todoPort)
		if err != nil {
			log.Fatalln(err)
		}

		return strconv.Itoa(intTodoPort)
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
		return intTodoTaskLimit
	}

	return taskLimit
}
