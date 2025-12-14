package sources

import (
	"encoding/csv"
	"os"
)

// CSVDataSource reads from CSV files
type CSVDataSource struct {
	filepath string
	file     *os.File
}

func NewCSVDataSource(filepath string) *CSVDataSource {
	return &CSVDataSource{filepath: filepath}
}

func (ds *CSVDataSource) Read() ([]map[string]interface{}, error) {
	file, err := os.Open(ds.filepath)
	if err != nil {
		return nil, err
	}
	ds.file = file

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return []map[string]interface{}{}, nil
	}

	headers := records[0]
	data := make([]map[string]interface{}, 0, len(records)-1)

	for _, record := range records[1:] {
		row := make(map[string]interface{})
		for i, value := range record {
			if i < len(headers) {
				row[headers[i]] = value
			}
		}
		data = append(data, row)
	}

	return data, nil
}

func (ds *CSVDataSource) Close() error {
	if ds.file != nil {
		return ds.file.Close()
	}
	return nil
}
