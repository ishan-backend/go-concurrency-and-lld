package logger

import (
	"context"

	"github.com/ishanpandey/lowleveldesign/multithreadedlogging/model"
)

// ApplicationAwareLogger interface defines the contract for application-level logging
type ApplicationAwareLogger interface {
	Info(ctx context.Context, message *model.ApplicationAwareMessage) error
	Debug(ctx context.Context, message *model.ApplicationAwareMessage) error
	Warn(ctx context.Context, message *model.ApplicationAwareMessage) error
	Error(ctx context.Context, message *model.ApplicationAwareMessage) error
}
