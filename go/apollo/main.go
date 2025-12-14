package main

import (
	"apollo_test/pipeline"
	"apollo_test/producers"
	"apollo_test/sinks"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"sync"
)

func main() {
	// Example 1: Simple CSV processing with builder pattern
	// fmt.Println("=== Example 1: CSV Processing ===\n")

	// Create sample CSV data
	createSampleCSV()

	pipeline := NewPipelineBuilder().
		Filter("status", "equals", "active").
		Transform("name", "uppercase").
		Build()

	processor := NewDataProcessor(
		NewCSVDataSource("sample_data.csv"),
		pipeline,
		NewConsoleDataSink(),
	)

	if err := processor.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// Example 2: Aggregation pipeline
	// fmt.Println("=== Example 2: Aggregation ===\n")

	aggregatePipeline := NewPipelineBuilder().
		Filter("status", "not_empty", nil).
		Aggregate("status").
		Build()

	aggregateProcessor := NewDataProcessor(
		NewCSVDataSource("sample_data.csv"),
		aggregatePipeline,
		NewConsoleDataSink(),
	)

	if err := aggregateProcessor.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// Example 3: Multi-pipeline processing (read once, fan-out to multiple pipelines)
	fmt.Println("=== Example 3: Multi-Pipeline Processing ===")
	fmt.Println()

	// Create multiple pipelines
	filterTransformPipeline := NewPipelineBuilder().
		Filter("status", "equals", "active").
		Transform("name", "uppercase").
		Build()

	engineeringAggregatePipeline := NewPipelineBuilder().
		Filter("department", "equals", "engineering").
		Aggregate("status").
		Build()

	departmentFilterPipeline := NewPipelineBuilder().
		Filter("department", "not_empty", nil).
		Build()

	// Create multi-pipeline processor
	// Reads from source ONCE, fans out to 3 pipelines in parallel, fans in results
	multiProcessor := NewMultiPipelineDataProcessor(
		NewCSVDataSource("sample_data.csv"), // Read once
		[]*Pipeline{
			filterTransformPipeline,
			engineeringAggregatePipeline,
			departmentFilterPipeline,
		},
		NewConsoleDataSink(), // Write merged results
	)

	if err := multiProcessor.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// Example 4: Channel-based pipeline with multiple producer workers
	fmt.Println("=== Example 4: Channel-Based Pipeline with Multiple Producers ===")
	fmt.Println()

	runChannelPipelineExample()
}

// runChannelPipelineExample demonstrates multiple producer workers producing to the same channel
// and a pipeline consuming, filtering, transforming, and aggregating the data
func runChannelPipelineExample() {
	// Create a channel for records
	recordChan := make(chan map[string]interface{}, 10) // Buffered channel

	// Create pipeline with filter, transform, and aggregate
	channelPipeline := pipeline.NewPipelineBuilder().
		Filter("status", "equals", "active").
		Transform("name", "uppercase").
		Aggregate("department").
		Build()

	// Create console sink
	consoleSink := sinks.NewConsoleDataSink()

	// Create channel pipeline processor
	channelProcessor := pipeline.NewChannelPipeline(channelPipeline, consoleSink)

	// Create sample data for multiple producers
	producer1Data := []map[string]interface{}{
		{"id": "1", "name": "john doe", "status": "active", "department": "engineering"},
		{"id": "2", "name": "jane smith", "status": "active", "department": "marketing"},
		{"id": "3", "name": "bob jones", "status": "inactive", "department": "engineering"},
	}

	producer2Data := []map[string]interface{}{
		{"id": "4", "name": "alice wilson", "status": "active", "department": "sales"},
		{"id": "5", "name": "charlie brown", "status": "active", "department": "engineering"},
		{"id": "6", "name": "diana prince", "status": "active", "department": "marketing"},
	}

	producer3Data := []map[string]interface{}{
		{"id": "7", "name": "eve adams", "status": "inactive", "department": "sales"},
		{"id": "8", "name": "frank miller", "status": "active", "department": "engineering"},
	}

	// Create wait group for producers
	var producerWg sync.WaitGroup

	// Start multiple producer workers
	producer1 := producers.NewProducerWorker(1, producer1Data, recordChan, &producerWg)
	producer2 := producers.NewProducerWorker(2, producer2Data, recordChan, &producerWg)
	producer3 := producers.NewProducerWorker(3, producer3Data, recordChan, &producerWg)

	producerWg.Add(3)
	producer1.Start()
	producer2.Start()
	producer3.Start()

	// Start a goroutine to close the channel when all producers are done
	go func() {
		producerWg.Wait()
		close(recordChan)
		fmt.Println("\nAll producers finished, channel closed")
	}()

	// Process from channel (this will block until channel is closed)
	doneChan := make(chan struct{})
	if err := channelProcessor.ProcessFromChannel(recordChan, doneChan); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// Helper function to create sample data
func createSampleCSV() {
	file, _ := os.Create("sample_data.csv")
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"id", "name", "status", "department"})
	writer.Write([]string{"1", "john doe", "active", "engineering"})
	writer.Write([]string{"2", "jane smith", "active", "marketing"})
	writer.Write([]string{"3", "bob jones", "inactive", "engineering"})
	writer.Write([]string{"4", "alice wilson", "active", "sales"})
}
