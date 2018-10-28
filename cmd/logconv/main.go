package main

import (
	"fmt"
	"os"
)

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}
	return value
}

func main() {
	inputLogFile := getEnv("INPUT_LOG_FILE", "/var/log/nginx/access.log")
	batchTime := getEnv("BATCH_TIME", "5")
	serverType := getEnv("SERVER_TYPE", "Nginx")

	fmt.Println(batchTime)
	fmt.Println(serverType)
	fmt.Println(inputLogFile)
}
