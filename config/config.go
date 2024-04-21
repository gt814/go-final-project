package config

import (
	"os"
	"strconv"
)

func GetPort() int {
	todoPort := os.Getenv("TODO_PORT")
	if len(todoPort) > 0 {
		if eport, err := strconv.ParseInt(todoPort, 10, 32); err == nil {
			port = int(eport)
		}
	}

	return port
}

func GetDBFileAppPath() string {
	todoDbFile := os.Getenv("TODO_DBFILE")

	if len(todoDbFile) > 0 {
		return todoDbFile
	}
	return dbFile
}

func GetDBFileTestPath() string {
	return dBFileTestPath + GetDBFileAppPath()
}

func GetFullNextDate() bool {
	return fullNextDate
}

func GetSearch() bool {
	return search
}

func GetToken() string {
	return token
}
