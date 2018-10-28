package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/logconv/pkg/logconv"
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

	batchInterval, err := strconv.Atoi(batchTime)
	if err != nil {
		fmt.Printf("Invalid batch time: %s", batchTime)
		os.Exit(1)
	}

	config := logconv.LogConvConf{
		InputLogFilePath: inputLogFile,
		BatchTime:        batchInterval,
		Type:             logconv.LogConvBatchType,
		ServerType:       serverType,
	}
	lc, err := logconv.NewLogConv(config)
	if err != nil {
		fmt.Printf("Error creating LogConv (%v)", err)
	}
	err = lc.Start()
	if err != nil {
		fmt.Printf("Could not start logconv (%v)", err)
	}
}
