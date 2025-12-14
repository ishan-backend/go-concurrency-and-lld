package sinks

// ============================================================================
// ADAPTER PATTERN - Support multiple output destinations
// ============================================================================

// DataSink defines interface for writing data
type DataSink interface {
	Write(data []map[string]interface{}) error
	Close() error
}
