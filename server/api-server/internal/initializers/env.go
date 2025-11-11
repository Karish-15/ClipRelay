package initializers

import (
	"os"
	"strconv"
)

func GetEnvInt(key string, def int) int {
	val := os.Getenv(key)
	if val == "" {
		return def
	}

	n, err := strconv.Atoi(val)
	if err != nil {
		return def
	}

	return n
}
