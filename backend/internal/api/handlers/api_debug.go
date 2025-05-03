package handlers

import (
	"os"
	"strconv"
)

func IsAPIDebugActive() bool {
	debugOn, err := strconv.ParseBool(os.Getenv("API_DEBUG"))
	if err != nil {
		return false
	}
	return debugOn
}
