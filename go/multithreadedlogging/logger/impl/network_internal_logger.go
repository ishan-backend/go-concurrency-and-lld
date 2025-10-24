package impl

import (
	"fmt"
	"sync"

	"github.com/ishanpandey/lowleveldesign/multithreadedlogging/logger"
	"github.com/ishanpandey/lowleveldesign/multithreadedlogging/model"
)

// networkInternalLogger implements InternalLogger for network logging
type networkInternalLogger struct {
	// In a real implementation, this would contain network connection details
}

var (
	networkLoggerInstance *networkInternalLogger
	networkLoggerOnce     sync.Once
)

// GetNetworkInternalLogger returns a singleton instance of NetworkInternalLogger
func GetNetworkInternalLogger() logger.InternalLogger {
	networkLoggerOnce.Do(func() {
		networkLoggerInstance = &networkInternalLogger{}
		// Initialize network connection here
		fmt.Println("Initializing Network logger connection")
	})
	return networkLoggerInstance
}

// Log writes the message to the network
func (n *networkInternalLogger) Log(message *model.WritableMessage) {
	// In a real implementation, this would send over the network
	fmt.Printf("[Network Logger] %s\n", message.String())
}
