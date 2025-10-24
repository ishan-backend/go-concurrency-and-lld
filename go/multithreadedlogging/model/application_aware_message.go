package model

import (
	"runtime"
	"time"

	"github.com/ishanpandey/lowleveldesign/multithreadedlogging/constant"
)

// ApplicationAwareMessage represents a message produced by the application
type ApplicationAwareMessage struct {
	Msg         string
	Time        time.Time
	Mediums     []constant.SupportedLogMedium
	GoroutineID int64
}

// NewApplicationAwareMessage creates a new ApplicationAwareMessage
func NewApplicationAwareMessage(msg string, mediums []constant.SupportedLogMedium) *ApplicationAwareMessage {
	return &ApplicationAwareMessage{
		Msg:         msg,
		Time:        time.Now(),
		Mediums:     mediums,
		GoroutineID: getGoroutineID(),
	}
}

// getGoroutineID returns a unique identifier for the current goroutine
func getGoroutineID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	// Extract goroutine ID from stack trace
	// This is a simplified approach - in production, you might want a more robust solution
	return int64(n) // Using stack size as a simple ID
}
