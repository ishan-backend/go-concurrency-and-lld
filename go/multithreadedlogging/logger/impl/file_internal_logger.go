package impl

import (
	"fmt"
	"sync"

	"github.com/ishanpandey/lowleveldesign/multithreadedlogging/logger"
	"github.com/ishanpandey/lowleveldesign/multithreadedlogging/model"
)

// fileInternalLogger implements InternalLogger for file logging
type fileInternalLogger struct {
	// In a real implementation, this would contain file handle details
}

var (
	fileLoggerInstance *fileInternalLogger
	fileLoggerOnce     sync.Once
)

// GetFileInternalLogger returns a singleton instance of FileInternalLogger
func GetFileInternalLogger() logger.InternalLogger {
	fileLoggerOnce.Do(func() {
		fileLoggerInstance = &fileInternalLogger{}
		// Initialize file handle here
		fmt.Println("Initializing File logger connection")
	})
	return fileLoggerInstance
}

// Log writes the message to the file
func (f *fileInternalLogger) Log(message *model.WritableMessage) {
	// In a real implementation, this would write to a file
	fmt.Printf("[File Logger] %s\n", message.String())
}
