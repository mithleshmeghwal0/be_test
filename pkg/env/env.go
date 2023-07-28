package env

import (
	"fmt"
	"os"
	"strconv"
)

func MustGet(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprint("Required environment variable not set: ", key))
	}

	return val
}

func GetBool(key string, def ...bool) bool {
	if val, err := strconv.ParseBool(os.Getenv(key)); err == nil {
		return val
	}

	if len(def) != 0 {
		return def[0]
	}

	return false
}
