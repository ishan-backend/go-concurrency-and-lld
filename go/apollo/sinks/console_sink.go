package sinks

import (
	"encoding/json"
	"os"
)

// ConsoleDataSink writes to stdout
type ConsoleDataSink struct{}

func NewConsoleDataSink() *ConsoleDataSink {
	return &ConsoleDataSink{}
}

func (ds *ConsoleDataSink) Write(data []map[string]interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (ds *ConsoleDataSink) Close() error {
	return nil
}
