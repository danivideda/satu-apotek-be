package env

import (
	"log"
	"os"
	"strconv"
	"time"
)

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		return fallback
	}

	return val
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		return fallback
	}

	valInt, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("%s. Returning the fallback val", err.Error())
		return fallback
	}

	return valInt
}

func GetDuration(key string, fallback time.Duration) time.Duration {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(val)
	if err != nil {
		return fallback
	}

	return parsed
}

func GetBool(key string, fallback bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		return fallback
	}

	switch val {
	case "true":
		return true
	case "false":
		return false
	default:
		return fallback
	}
}
