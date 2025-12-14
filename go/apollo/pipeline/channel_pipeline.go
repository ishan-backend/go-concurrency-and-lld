package pipeline

import (
	"apollo_test/sinks"
	"fmt"
	"sync"
)

// ChannelPipeline consumes records from a channel, processes them through the pipeline,
// and sends results to a sink
type ChannelPipeline struct {
	pipeline *Pipeline
	sink     sinks.DataSink
}

// NewChannelPipeline creates a new channel-based pipeline processor
func NewChannelPipeline(pipeline *Pipeline, sink sinks.DataSink) *ChannelPipeline {
	return &ChannelPipeline{
		pipeline: pipeline,
		sink:     sink,
	}
}

// ProcessFromChannel consumes records from the channel, processes them through the pipeline,
// and writes results to the sink
func (cp *ChannelPipeline) ProcessFromChannel(inputChan <-chan map[string]interface{}, doneChan <-chan struct{}) error {
	defer cp.sink.Close()

	// Collect all records from channel
	var records []map[string]interface{}
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Start a goroutine to collect records from the channel
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case record, ok := <-inputChan:
				if !ok {
					// Channel closed, stop collecting
					return
				}
				mu.Lock()
				records = append(records, record)
				mu.Unlock()
				fmt.Printf("Pipeline: Received record: %v\n", record)
			case <-doneChan:
				// Done signal received (optional early termination)
				return
			}
		}
	}()

	// Wait for all records to be collected (channel will be closed by producers)
	wg.Wait()

	fmt.Printf("\nPipeline: Collected %d records from channel\n", len(records))

	// Process records through the pipeline
	fmt.Println("Pipeline: Processing records through pipeline...")
	result, err := cp.pipeline.Execute(records)
	if err != nil {
		return fmt.Errorf("error processing data: %w", err)
	}
	fmt.Printf("Pipeline: Processed %d records\n\n", len(result))

	// Write results to sink
	fmt.Println("Pipeline: Writing results to sink...")
	if err := cp.sink.Write(result); err != nil {
		return fmt.Errorf("error writing data to sink: %w", err)
	}
	fmt.Println("Pipeline: Done!")

	return nil
}

