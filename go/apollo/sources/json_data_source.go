package sources

import (
	"encoding/json"
	"os"
)

// JSONDataSource reads from JSON files
type JSONDataSource struct {
	filepath string
	file     *os.File
}

func NewJSONDataSource(filepath string) *JSONDataSource {
	return &JSONDataSource{filepath: filepath}
}

func (ds *JSONDataSource) Read() ([]map[string]interface{}, error) {
	file, err := os.Open(ds.filepath)
	if err != nil {
		return nil, err
	}
	ds.file = file

	var data []map[string]interface{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

func (ds *JSONDataSource) Close() error {
	if ds.file != nil {
		return ds.file.Close()
	}
	return nil
}
