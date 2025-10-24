package model

import (
	"fmt"

	"github.com/ishanpandey/lowleveldesign/multithreadedlogging/constant"
)

// WritableMessage wraps an ApplicationAwareMessage with additional metadata for logging
type WritableMessage struct {
	ApplicationAwareMessage *ApplicationAwareMessage
	LogLevel                constant.LogLevel
	Medium                  constant.SupportedLogMedium
}

// NewWritableMessage creates a new WritableMessage
func NewWritableMessage(appMsg *ApplicationAwareMessage, logLevel constant.LogLevel, medium constant.SupportedLogMedium) *WritableMessage {
	return &WritableMessage{
		ApplicationAwareMessage: appMsg,
		LogLevel:                logLevel,
		Medium:                  medium,
	}
}

// String returns the string representation of WritableMessage
func (w *WritableMessage) String() string {
	return fmt.Sprintf("Log Level = %s Goroutine ID = %d timestamp = %s message = %s",
		w.LogLevel.String(),
		w.ApplicationAwareMessage.GoroutineID,
		w.ApplicationAwareMessage.Time.Format("2006-01-02 15:04:05"),
		w.ApplicationAwareMessage.Msg)
}
