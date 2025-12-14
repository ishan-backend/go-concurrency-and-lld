package sinks

import (
	"encoding/json"
	"os"
)

// JSONDataSink writes to JSON files
type JSONDataSink struct {
	filepath string
	file     *os.File
}

func NewJSONDataSink(filepath string) *JSONDataSink {
	return &JSONDataSink{filepath: filepath}
}

func (ds *JSONDataSink) Write(data []map[string]interface{}) error {
	file, err := os.Create(ds.filepath)
	if err != nil {
		return err
	}
	ds.file = file

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (ds *JSONDataSink) Close() error {
	if ds.file != nil {
		return ds.file.Close()
	}
	return nil
}
