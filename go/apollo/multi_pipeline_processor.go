package main

import (
	"fmt"
	"sync"
)

// ============================================================================
// MULTI-PIPELINE DATA PROCESSOR
// ============================================================================
// Implements fan-out/fan-in pattern for parallel pipeline execution
// Follows SOLID principles:
// - Single Responsibility: Handles multi-pipeline orchestration
// - Open/Closed: Extends functionality without modifying existing code
// - Liskov Substitution: Can be used where DataProcessor is expected
// - Interface Segregation: Uses existing interfaces
// - Dependency Inversion: Depends on abstractions (DataSource, Pipeline, DataSink)
// ============================================================================

// MultiPipelineDataProcessor handles one source, multiple pipelines, one sink
// Reads from source once, fans out to multiple pipelines in parallel,
// and fans in results to a single sink in a concurrent-safe manner
type MultiPipelineDataProcessor struct {
	source    DataSource
	pipelines []*Pipeline
	sink      DataSink
}

// PipelineResult represents the result from a single pipeline execution
type PipelineResult struct {
	PipelineIndex int
	PipelineName  string
	Data          []map[string]interface{}
	Error         error
}

// NewMultiPipelineDataProcessor creates a new multi-pipeline processor
// Parameters:
//   - source: DataSource to read from (read once)
//   - pipelines: Slice of pipelines to execute in parallel
//   - sink: DataSink to write merged results to
func NewMultiPipelineDataProcessor(
	source DataSource,
	pipelines []*Pipeline,
	sink DataSink,
) *MultiPipelineDataProcessor {
	return &MultiPipelineDataProcessor{
		source:    source,
		pipelines: pipelines,
		sink:      sink,
	}
}

// Execute reads from source once, fans out to multiple pipelines in parallel,
// collects results safely, merges them, and writes to sink
func (mp *MultiPipelineDataProcessor) Execute() error {
	defer mp.source.Close()
	defer mp.sink.Close()

	// Step 1: Read data from source once
	fmt.Println("Reading data from source (once)...")
	data, err := mp.source.Read()
	if err != nil {
		return fmt.Errorf("error reading data: %w", err)
	}
	fmt.Printf("Read %d records\n\n", len(data))

	// Step 2: Fan-out - Launch goroutines for each pipeline
	fmt.Printf("Fanning out to %d pipelines in parallel...\n", len(mp.pipelines))
	resultChan := make(chan PipelineResult, len(mp.pipelines))
	var wg sync.WaitGroup

	// Launch a goroutine for each pipeline
	for i, pipeline := range mp.pipelines {
		wg.Add(1)
		go func(pipelineIndex int, p *Pipeline) {
			defer wg.Done()

			// Execute pipeline with a deep copy of the data
			// Each pipeline gets its own copy to avoid race conditions
			pipelineData := deepCopyData(data)

			result, err := p.Execute(pipelineData)

			// Send result to channel
			resultChan <- PipelineResult{
				PipelineIndex: pipelineIndex,
				PipelineName:  fmt.Sprintf("Pipeline-%d", pipelineIndex+1),
				Data:          result,
				Error:         err,
			}
		}(i, pipeline)
	}

	// Step 3: Close channel when all goroutines complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Step 4: Fan-in - Collect results safely from all pipelines
	fmt.Println("Fanning in results from all pipelines...")
	var allResults []map[string]interface{}
	var errors []error
	successCount := 0

	for result := range resultChan {
		if result.Error != nil {
			errors = append(errors, fmt.Errorf("%s: %w", result.PipelineName, result.Error))
			fmt.Printf("  ✗ %s failed: %v\n", result.PipelineName, result.Error)
		} else {
			allResults = append(allResults, result.Data...)
			successCount++
			fmt.Printf("  ✓ %s completed: %d records\n", result.PipelineName, len(result.Data))
		}
	}

	fmt.Printf("\nProcessed %d pipelines successfully, %d failed\n", successCount, len(errors))
	fmt.Printf("Total merged records: %d\n\n", len(allResults))

	// Step 5: Handle errors (continue on error strategy)
	if len(errors) > 0 {
		fmt.Printf("Warning: %d pipeline(s) encountered errors:\n", len(errors))
		for _, err := range errors {
			fmt.Printf("  - %v\n", err)
		}
		fmt.Println()
	}

	// Step 6: Write merged results to sink
	fmt.Println("Writing merged results to sink...")
	if err := mp.sink.Write(allResults); err != nil {
		return fmt.Errorf("error writing data to sink: %w", err)
	}
	fmt.Println("Done!")

	return nil
}

// GetPipelineCount returns the number of pipelines configured
func (mp *MultiPipelineDataProcessor) GetPipelineCount() int {
	return len(mp.pipelines)
}

// deepCopyData creates a deep copy of the data slice and maps
// This ensures each pipeline operates on its own copy, preventing race conditions
func deepCopyData(data []map[string]interface{}) []map[string]interface{} {
	if data == nil {
		return nil
	}

	copied := make([]map[string]interface{}, len(data))
	for i, row := range data {
		copiedRow := make(map[string]interface{})
		for k, v := range row {
			copiedRow[k] = v
		}
		copied[i] = copiedRow
	}

	return copied
}
