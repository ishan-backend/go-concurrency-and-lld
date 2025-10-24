package factory

import (
	"github.com/ishanpandey/lowleveldesign/multithreadedlogging/constant"
	"github.com/ishanpandey/lowleveldesign/multithreadedlogging/logger"
	"github.com/ishanpandey/lowleveldesign/multithreadedlogging/logger/impl"
)

// InternalLoggerFactory provides factory methods for creating InternalLogger instances
type InternalLoggerFactory struct{}

// NewInternalLoggerFactory creates a new InternalLoggerFactory
func NewInternalLoggerFactory() *InternalLoggerFactory {
	return &InternalLoggerFactory{}
}

// GetLoggerInstance returns an InternalLogger instance based on the medium
func (f *InternalLoggerFactory) GetLoggerInstance(medium constant.SupportedLogMedium) logger.InternalLogger {
	switch medium {
	case constant.DB:
		return impl.GetDBInternalLogger()
	case constant.FILE:
		return impl.GetFileInternalLogger()
	case constant.NETWORK:
		return impl.GetNetworkInternalLogger()
	default:
		return nil
	}
}
