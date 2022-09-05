package util

import (
	"fmt"
	"os"
	"strconv"
)

func GetEnvRequired(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("environment variable '%s' must be set", value)
	}
	return value, nil
}

func GetEnvDefault(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}

func GetEnvIntDefault(key string, defaultValue int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	vi, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return vi
}
