package main

import (
	"fmt"
)

// ============================================================================
// MAIN EXECUTION ENGINE
// ============================================================================

type DataProcessor struct {
	source   DataSource
	pipeline *Pipeline
	sink     DataSink
}

func NewDataProcessor(source DataSource, pipeline *Pipeline, sink DataSink) *DataProcessor {
	return &DataProcessor{
		source:   source,
		pipeline: pipeline,
		sink:     sink,
	}
}

func (dp *DataProcessor) Execute() error {
	defer dp.source.Close()
	defer dp.sink.Close()

	// Read data
	fmt.Println("Reading data from source...")
	data, err := dp.source.Read()
	if err != nil {
		return fmt.Errorf("error reading data: %w", err)
	}
	fmt.Printf("Read %d records\n\n", len(data))

	// Process data
	fmt.Println("Processing data through pipeline...")
	result, err := dp.pipeline.Execute(data)
	if err != nil {
		return fmt.Errorf("error processing data: %w", err)
	}
	fmt.Printf("\nProcessed %d records\n\n", len(result))

	// Write data
	fmt.Println("Writing data to sink...")
	if err := dp.sink.Write(result); err != nil {
		return fmt.Errorf("error writing data: %w", err)
	}
	fmt.Println("Done!")

	return nil
}
