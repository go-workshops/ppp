package logging

import (
	"log"
	"os"
	"time"
)

// CustomLogger wraps the standard log.Logger to simulate slow logging.
type CustomLogger struct {
	logger *log.Logger
}

// NewCustomLogger creates a new instance of CustomLogger.
func NewCustomLogger() *CustomLogger {
	return &CustomLogger{
		logger: log.New(os.Stdout, "CUSTOM: ", log.LstdFlags),
	}
}

// Println simulates a slow logging operation by introducing a delay.
func (c *CustomLogger) Println(v ...interface{}) {
	// Simulate slow log writing
	time.Sleep(2 * time.Second) // Introduce a delay of 2 seconds
	c.logger.Println(v...)
}

func main() {
	// Create an instance of CustomLogger
	customLogger := NewCustomLogger()

	// Simulate logging
	startTime := time.Now()
	customLogger.Println("This is a synchronous log message with delay.")
	elapsedTime := time.Since(startTime)

	// Print elapsed time
	log.Printf("Logging took %s\n", elapsedTime)
}
