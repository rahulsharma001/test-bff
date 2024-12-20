package utils

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/spf13/viper"
)

var (
	loggerInstance *fluent.Fluent
	once           sync.Once
)

// InitializeFluent initializes the Fluent logger once
func InitializeFluent() *fluent.Fluent {
	once.Do(func() {
		var err error
		loggerInstance, err = fluent.New(fluent.Config{
			FluentHost: viper.GetString("FLUENT_HOST"),
			FluentPort: viper.GetInt("FLUENT_PORT"),
		})
		if err != nil {
			log.Fatalf("Failed to connect to Fluentd: %v", err)
		}
	})
	return loggerInstance
}

// LogToFluent logs a message with structured data to Fluentd
func Log(message interface{}) {
	fmt.Println(message)
	if viper.GetString("ENVIRONMENT") == "local" {
		return
	}

	// Ensure Fluentd logger is initialized
	logger := InitializeFluent()

	// Send the log message with the given tag

	tag := viper.GetString("FLUENT_TAG")

	pc, _, line, _ := runtime.Caller(1)

	funcName := runtime.FuncForPC(pc).Name()

	// Convert the string message into a structured log map
	logData := map[string]interface{}{
		"function":  funcName,                        // Log the function name
		"line":      line,                            // Log the line number
		"message":   message,                         // Log the message
		"timestamp": time.Now().Format(time.RFC3339), // Add a timestamp
		"level":     "info",                          // You can add a log level if desired
	}

	err := logger.Post(tag, logData)
	if err != nil {
		log.Printf("Failed to send log to Fluentd: %v", err)
	}
}

// CloseFluent closes the Fluent logger connection
func CloseFluent() {
	if loggerInstance != nil {
		loggerInstance.Close()
	}
}
