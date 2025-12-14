# Multi-Pipeline Data Processor - Implementation Plan

## Overview
Extend Apollo to support reading from one source once, fanning out to multiple pipelines in parallel, and fanning in results to one sink in a concurrent-safe environment.

## SOLID Principles Applied

### 1. Single Responsibility Principle (SRP)
- **Current**: `DataProcessor` handles single source → single pipeline → single sink
- **New**: Create `MultiPipelineDataProcessor` with separate responsibility for multi-pipeline orchestration
- **Separation**: Keep existing `DataProcessor` unchanged, add new processor type

### 2. Open/Closed Principle (OCP)
- **Open for Extension**: New `MultiPipelineDataProcessor` extends functionality without modifying existing code
- **Closed for Modification**: `DataProcessor` remains untouched
- **Interface-based**: Use existing interfaces (`DataSource`, `Pipeline`, `DataSink`)

### 3. Liskov Substitution Principle (LSP)
- Both processors can be used interchangeably where a common interface is expected
- Both respect the same contracts (error handling, resource cleanup)

### 4. Interface Segregation Principle (ISP)
- Create minimal interface `PipelineExecutor` if needed
- Use existing interfaces without modification
- No client forced to depend on unused methods

### 5. Dependency Inversion Principle (DIP)
- Depend on abstractions (`DataSource`, `DataSink`, `Pipeline` interfaces)
- Not on concrete implementations
- Use dependency injection in constructors

## Architecture Design

### Components

1. **MultiPipelineDataProcessor** (New File: `multi_pipeline_processor.go`)
   - Orchestrates: source read → fan-out → parallel execution → fan-in → sink write
   - Manages concurrency and error handling
   - Thread-safe result collection

2. **PipelineResult** (Internal Type)
   - Encapsulates result from a single pipeline execution
   - Includes pipeline identifier and error information

3. **ResultCollector** (Internal Component)
   - Thread-safe collection of results from multiple pipelines
   - Uses channels or mutex-protected slice
   - Merges results in order

### Design Patterns

1. **Fan-Out/Fan-In Pattern**
   - Fan-out: Distribute data to multiple goroutines
   - Fan-in: Collect results from multiple goroutines safely

2. **Worker Pool Pattern** (Optional)
   - Control concurrency level
   - Prevent resource exhaustion

3. **Strategy Pattern** (Already in use)
   - Different pipelines as different strategies
   - No modification needed

## Implementation Structure

```
apollo/
├── data_processor.go              (UNCHANGED)
├── multi_pipeline_processor.go    (NEW - main implementation)
└── main.go                        (UPDATED - add example usage)
```

## Detailed Implementation Plan

### File 1: `multi_pipeline_processor.go`

```go
package main

import (
    "fmt"
    "sync"
)

// MultiPipelineDataProcessor handles one source, multiple pipelines, one sink
type MultiPipelineDataProcessor struct {
    source    DataSource
    pipelines []*Pipeline
    sink      DataSink
}

// PipelineResult represents result from a single pipeline execution
type PipelineResult struct {
    PipelineIndex int
    Data          []map[string]interface{}
    Error         error
}

// NewMultiPipelineDataProcessor creates a new multi-pipeline processor
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

// Execute reads once, fans out to pipelines, fans in results, writes to sink
func (mp *MultiPipelineDataProcessor) Execute() error {
    // 1. Read from source once
    // 2. Fan-out: Launch goroutines for each pipeline
    // 3. Fan-in: Collect results safely
    // 4. Merge results
    // 5. Write to sink
}
```

### Key Features

1. **Single Source Read**
   - Read once at the beginning
   - Share data with all pipelines

2. **Parallel Execution**
   - Each pipeline runs in its own goroutine
   - Use `sync.WaitGroup` to wait for completion
   - Use channels for result collection

3. **Concurrent-Safe Result Collection**
   - Option A: Buffered channel for results
   - Option B: Mutex-protected slice
   - Recommended: Channel-based (more idiomatic Go)

4. **Error Handling**
   - Collect errors from each pipeline
   - Continue processing other pipelines even if one fails
   - Aggregate errors or fail-fast based on configuration

5. **Result Merging**
   - Combine results from all pipelines
   - Preserve order or merge based on strategy
   - Handle duplicates if needed

## Implementation Details

### Fan-Out Strategy
```go
// Launch goroutines for each pipeline
var wg sync.WaitGroup
resultChan := make(chan PipelineResult, len(mp.pipelines))

for i, pipeline := range mp.pipelines {
    wg.Add(1)
    go func(idx int, p *Pipeline) {
        defer wg.Done()
        result, err := p.Execute(data)
        resultChan <- PipelineResult{
            PipelineIndex: idx,
            Data:          result,
            Error:         err,
        }
    }(i, pipeline)
}

// Close channel when all goroutines complete
go func() {
    wg.Wait()
    close(resultChan)
}()
```

### Fan-In Strategy
```go
// Collect results safely
var allResults []map[string]interface{}
var errors []error

for result := range resultChan {
    if result.Error != nil {
        errors = append(errors, result.Error)
        continue
    }
    allResults = append(allResults, result.Data...)
}
```

### Error Handling Strategy
- **Option 1**: Continue on error, collect all errors
- **Option 2**: Fail-fast on first error
- **Option 3**: Configurable behavior

## Usage Example

```go
// In main.go
func main() {
    // ... existing code ...
    
    // Example 3: Multi-pipeline processing
    source := NewCSVDataSource("sample_data.csv")
    
    pipeline1 := NewPipelineBuilder().
        Filter("status", "equals", "active").
        Transform("name", "uppercase").
        Build()
    
    pipeline2 := NewPipelineBuilder().
        Filter("department", "equals", "engineering").
        Aggregate("status").
        Build()
    
    multiProcessor := NewMultiPipelineDataProcessor(
        source,
        []*Pipeline{pipeline1, pipeline2},
        NewConsoleDataSink(),
    )
    
    if err := multiProcessor.Execute(); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

## Benefits

1. **No Breaking Changes**: Existing `DataProcessor` remains unchanged
2. **Extensible**: Easy to add more pipelines
3. **Concurrent-Safe**: Proper synchronization mechanisms
4. **SOLID Compliant**: Follows all five principles
5. **Testable**: Each component can be tested independently
6. **Maintainable**: Clear separation of concerns

## Testing Considerations

1. Unit tests for result collection
2. Concurrent execution tests
3. Error handling tests
4. Integration tests with real pipelines
5. Performance benchmarks

## Future Enhancements

1. **Configurable Concurrency**: Limit number of parallel pipelines
2. **Result Merging Strategies**: Union, intersection, custom
3. **Pipeline Prioritization**: Execute critical pipelines first
4. **Result Filtering**: Filter results before merging
5. **Metrics Collection**: Track execution time per pipeline

