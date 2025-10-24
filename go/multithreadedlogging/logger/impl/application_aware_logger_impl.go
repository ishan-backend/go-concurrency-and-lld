package impl

import (
	"context"

	"github.com/ishanpandey/lowleveldesign/multithreadedlogging/constant"
	"github.com/ishanpandey/lowleveldesign/multithreadedlogging/logger"
	"github.com/ishanpandey/lowleveldesign/multithreadedlogging/model"
	"github.com/ishanpandey/lowleveldesign/multithreadedlogging/service"
)

// applicationAwareLoggerImpl implements ApplicationAwareLogger
// Applications will talk to this logger instance to publish messages to the system
type applicationAwareLoggerImpl struct {
	queueService service.QueueService
}

// NewApplicationAwareLogger creates a new ApplicationAwareLogger instance
func NewApplicationAwareLogger(queueService service.QueueService) logger.ApplicationAwareLogger {
	return &applicationAwareLoggerImpl{
		queueService: queueService,
	}
}

// Info logs an info level message
func (a *applicationAwareLoggerImpl) Info(ctx context.Context, message *model.ApplicationAwareMessage) error {
	return a.queueService.Publish(ctx, a.getWritableMessages(message, constant.INFO))
}

// Debug logs a debug level message
func (a *applicationAwareLoggerImpl) Debug(ctx context.Context, message *model.ApplicationAwareMessage) error {
	return a.queueService.Publish(ctx, a.getWritableMessages(message, constant.DEBUG))
}

// Warn logs a warning level message
func (a *applicationAwareLoggerImpl) Warn(ctx context.Context, message *model.ApplicationAwareMessage) error {
	return a.queueService.Publish(ctx, a.getWritableMessages(message, constant.WARN))
}

// Error logs an error level message
func (a *applicationAwareLoggerImpl) Error(ctx context.Context, message *model.ApplicationAwareMessage) error {
	return a.queueService.Publish(ctx, a.getWritableMessages(message, constant.ERROR))
}

// getWritableMessages converts ApplicationAwareMessage to WritableMessage for each medium
func (a *applicationAwareLoggerImpl) getWritableMessages(appMsg *model.ApplicationAwareMessage, logLevel constant.LogLevel) []*model.WritableMessage {
	var writableMessages []*model.WritableMessage

	for _, medium := range appMsg.Mediums {
		writableMessages = append(writableMessages, model.NewWritableMessage(appMsg, logLevel, medium))
	}

	return writableMessages
}
