// Package env is responsible for getting environment variables from OS and verifying if it has been set.
package env

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// GetString gets string environment variables from OS.
func GetString(envVar string, defaultValue ...string) string {
	value := os.Getenv(envVar)
	if value == "" && len(defaultValue) > 0 {
		value = defaultValue[0]
	}

	return value
}

// GetInt gets int environment variables from OS.
func GetInt(envVar string, defaultValue int) int {
	if valueStr := os.Getenv(envVar); valueStr != "" {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

// CheckRequired verify if the environment variables passed in envVarArgs parameter are set in the OS.
func CheckRequired(envVarArgs ...string) {
	for _, envVar := range envVarArgs {
		if os.Getenv(envVar) == "" {
			fmt.Print("1")
			log.Fatalf("Environment variable '%s' is required.", envVar)
		}
		fmt.Print("2")
		log.Printf("Environment variable '%s' is ok.", envVar)
	}
}
