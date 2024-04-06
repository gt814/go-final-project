package config

import (
	"os"
	"strconv"
)

// GetPort returns the port on which the application will listen for requests.
func GetPort() int {
	port, _ := strconv.Atoi(os.Getenv("TODO_PORT"))
	if port != 0 {
		return port
	}
	return Port
}
