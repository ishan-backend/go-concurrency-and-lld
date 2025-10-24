package logger

import "github.com/ishanpandey/lowleveldesign/multithreadedlogging/model"

// InternalLogger interface defines the contract for internal logging operations
type InternalLogger interface {
	Log(message *model.WritableMessage)
}
