package impl

import (
	"fmt"
	"sync"

	"github.com/ishanpandey/lowleveldesign/multithreadedlogging/logger"
	"github.com/ishanpandey/lowleveldesign/multithreadedlogging/model"
)

// dbInternalLogger implements InternalLogger for database logging
type dbInternalLogger struct {
	// In a real implementation, this would contain database connection details
}

var (
	dbLoggerInstance *dbInternalLogger
	dbLoggerOnce     sync.Once
)

// GetDBInternalLogger returns a singleton instance of DBInternalLogger
func GetDBInternalLogger() logger.InternalLogger {
	dbLoggerOnce.Do(func() {
		dbLoggerInstance = &dbInternalLogger{}
		// Initialize database connection here
		fmt.Println("Initializing DB logger connection")
	})
	return dbLoggerInstance
}

// Log writes the message to the database
func (d *dbInternalLogger) Log(message *model.WritableMessage) {
	// In a real implementation, this would write to the database
	fmt.Printf("[DB Logger] %s\n", message.String())
}
