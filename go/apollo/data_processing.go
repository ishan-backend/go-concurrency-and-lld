package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ============================================================================
// STRATEGY PATTERN - Define processing strategies
// ============================================================================

// Processor defines the interface for data processing steps
type Processor interface {
	Process(data []map[string]interface{}) ([]map[string]interface{}, error)
	Name() string
}

// ============================================================================
// CONCRETE PROCESSOR IMPLEMENTATIONS
// ============================================================================

// FilterProcessor filters data based on conditions
type FilterProcessor struct {
	field     string
	operation string
	value     interface{}
}

func CreateFilterProcessor(config map[string]interface{}) (Processor, error) {
	field, ok := config["field"].(string)
	if !ok {
		return nil, fmt.Errorf("field is required for filter processor")
	}

	operation, ok := config["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation is required for filter processor")
	}

	value := config["value"]

	return &FilterProcessor{
		field:     field,
		operation: operation,
		value:     value,
	}, nil
}

func (fp *FilterProcessor) Process(data []map[string]interface{}) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)

	for _, row := range data {
		if fp.matchesCondition(row) {
			result = append(result, row)
		}
	}

	return result, nil
}

func (fp *FilterProcessor) matchesCondition(row map[string]interface{}) bool {
	value, exists := row[fp.field]
	if !exists {
		return false
	}

	switch fp.operation {
	case "equals":
		return fmt.Sprintf("%v", value) == fmt.Sprintf("%v", fp.value)
	case "contains":
		return strings.Contains(fmt.Sprintf("%v", value), fmt.Sprintf("%v", fp.value))
	case "not_empty":
		return value != nil && fmt.Sprintf("%v", value) != ""
	default:
		return false
	}
}

func (fp *FilterProcessor) Name() string {
	return fmt.Sprintf("Filter[%s %s %v]", fp.field, fp.operation, fp.value)
}

// TransformProcessor transforms field values
type TransformProcessor struct {
	field     string
	operation string
}

func CreateTransformProcessor(config map[string]interface{}) (Processor, error) {
	field, ok := config["field"].(string)
	if !ok {
		return nil, fmt.Errorf("field is required for transform processor")
	}

	operation, ok := config["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation is required for transform processor")
	}

	return &TransformProcessor{
		field:     field,
		operation: operation,
	}, nil
}

func (tp *TransformProcessor) Process(data []map[string]interface{}) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, len(data))

	for i, row := range data {
		newRow := make(map[string]interface{})
		for k, v := range row {
			newRow[k] = v
		}

		if value, exists := row[tp.field]; exists {
			newRow[tp.field] = tp.transform(value)
		}

		result[i] = newRow
	}

	return result, nil
}

func (tp *TransformProcessor) transform(value interface{}) interface{} {
	str := fmt.Sprintf("%v", value)

	switch tp.operation {
	case "uppercase":
		return strings.ToUpper(str)
	case "lowercase":
		return strings.ToLower(str)
	case "trim":
		return strings.TrimSpace(str)
	default:
		return value
	}
}

func (tp *TransformProcessor) Name() string {
	return fmt.Sprintf("Transform[%s: %s]", tp.field, tp.operation)
}

// AggregateProcessor groups and counts data
type AggregateProcessor struct {
	groupBy string
}

func CreateAggregateProcessor(config map[string]interface{}) (Processor, error) {
	groupBy, ok := config["group_by"].(string)
	if !ok {
		return nil, fmt.Errorf("group_by is required for aggregate processor")
	}

	return &AggregateProcessor{
		groupBy: groupBy,
	}, nil
}

func (ap *AggregateProcessor) Process(data []map[string]interface{}) ([]map[string]interface{}, error) {
	groups := make(map[string]int)

	for _, row := range data {
		if value, exists := row[ap.groupBy]; exists {
			key := fmt.Sprintf("%v", value)
			groups[key]++
		}
	}

	result := make([]map[string]interface{}, 0, len(groups))
	for key, count := range groups {
		result = append(result, map[string]interface{}{
			ap.groupBy: key,
			"count":    count,
		})
	}

	return result, nil
}

func (ap *AggregateProcessor) Name() string {
	return fmt.Sprintf("Aggregate[group_by: %s]", ap.groupBy)
}

// ============================================================================
// CHAIN OF RESPONSIBILITY PATTERN - Pipeline processing
// ============================================================================

// Pipeline chains multiple processors together
type Pipeline struct {
	processors []Processor
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		processors: make([]Processor, 0),
	}
}

func (p *Pipeline) AddProcessor(processor Processor) *Pipeline {
	p.processors = append(p.processors, processor)
	return p
}

func (p *Pipeline) Execute(data []map[string]interface{}) ([]map[string]interface{}, error) {
	result := data
	var err error

	for _, processor := range p.processors {
		fmt.Printf("Executing processor: %s\n", processor.Name())
		result, err = processor.Process(result)
		if err != nil {
			return nil, fmt.Errorf("error in %s: %w", processor.Name(), err)
		}
	}

	return result, nil
}

// ============================================================================
// FACTORY PATTERN - Create processors dynamically
// ============================================================================

type ProcessorFactory struct {
	creators map[string]func(config map[string]interface{}) (Processor, error)
}

func NewProcessorFactory() *ProcessorFactory {
	factory := &ProcessorFactory{
		creators: make(map[string]func(config map[string]interface{}) (Processor, error)),
	}

	// Register default processors
	factory.Register("filter", CreateFilterProcessor)
	factory.Register("transform", CreateTransformProcessor)
	factory.Register("aggregate", CreateAggregateProcessor)

	return factory
}

func (f *ProcessorFactory) Register(name string, creator func(config map[string]interface{}) (Processor, error)) {
	f.creators[name] = creator
}

func (f *ProcessorFactory) Create(name string, config map[string]interface{}) (Processor, error) {
	creator, exists := f.creators[name]
	if !exists {
		return nil, fmt.Errorf("unknown processor type: %s", name)
	}
	return creator(config)
}

// ============================================================================
// ADAPTER PATTERN - Support multiple data sources
// ============================================================================

// DataSource defines interface for reading data
type DataSource interface {
	Read() ([]map[string]interface{}, error)
	Close() error
}

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

// ============================================================================
// ADAPTER PATTERN - Support multiple output destinations
// ============================================================================

// DataSink defines interface for writing data
type DataSink interface {
	Write(data []map[string]interface{}) error
	Close() error
}

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

// ============================================================================
// BUILDER PATTERN - Fluent pipeline construction
// ============================================================================

type PipelineBuilder struct {
	pipeline *Pipeline
	factory  *ProcessorFactory
}

func NewPipelineBuilder() *PipelineBuilder {
	return &PipelineBuilder{
		pipeline: NewPipeline(),
		factory:  NewProcessorFactory(),
	}
}

func (pb *PipelineBuilder) Filter(field, operation string, value interface{}) *PipelineBuilder {
	config := map[string]interface{}{
		"field":     field,
		"operation": operation,
		"value":     value,
	}
	processor, err := pb.factory.Create("filter", config)
	if err != nil {
		fmt.Printf("Error creating filter: %v\n", err)
		return pb
	}
	pb.pipeline.AddProcessor(processor)
	return pb
}

func (pb *PipelineBuilder) Transform(field, operation string) *PipelineBuilder {
	config := map[string]interface{}{
		"field":     field,
		"operation": operation,
	}
	processor, err := pb.factory.Create("transform", config)
	if err != nil {
		fmt.Printf("Error creating transform: %v\n", err)
		return pb
	}
	pb.pipeline.AddProcessor(processor)
	return pb
}

func (pb *PipelineBuilder) Aggregate(groupBy string) *PipelineBuilder {
	config := map[string]interface{}{
		"group_by": groupBy,
	}
	processor, err := pb.factory.Create("aggregate", config)
	if err != nil {
		fmt.Printf("Error creating aggregate: %v\n", err)
		return pb
	}
	pb.pipeline.AddProcessor(processor)
	return pb
}

func (pb *PipelineBuilder) Build() *Pipeline {
	return pb.pipeline
}

// ============================================================================
// MAIN EXECUTION ENGINE
// ============================================================================

// type DataProcessor struct {
// 	source   DataSource
// 	pipeline *Pipeline
// 	sink     DataSink
// }

// func NewDataProcessor(source DataSource, pipeline *Pipeline, sink DataSink) *DataProcessor {
// 	return &DataProcessor{
// 		source:   source,
// 		pipeline: pipeline,
// 		sink:     sink,
// 	}
// }

// func (dp *DataProcessor) Execute() error {
// 	defer dp.source.Close()
// 	defer dp.sink.Close()

// 	// Read data
// 	fmt.Println("Reading data from source...")
// 	data, err := dp.source.Read()
// 	if err != nil {
// 		return fmt.Errorf("error reading data: %w", err)
// 	}
// 	fmt.Printf("Read %d records\n\n", len(data))

// 	// Process data
// 	fmt.Println("Processing data through pipeline...")
// 	result, err := dp.pipeline.Execute(data)
// 	if err != nil {
// 		return fmt.Errorf("error processing data: %w", err)
// 	}
// 	fmt.Printf("\nProcessed %d records\n\n", len(result))

// 	// Write data
// 	fmt.Println("Writing data to sink...")
// 	if err := dp.sink.Write(result); err != nil {
// 		return fmt.Errorf("error writing data: %w", err)
// 	}
// 	fmt.Println("Done!")

// 	return nil
// }

// ============================================================================
// EXAMPLE USAGE
// ============================================================================
