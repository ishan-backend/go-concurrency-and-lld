package constants

// Error messages for data sources
const (
	ErrFileOpenFailed     = "failed to open file"
	ErrFileReadFailed     = "failed to read file"
	ErrFileCloseFailed    = "failed to close file"
	ErrCSVReadFailed      = "failed to read CSV data"
	ErrJSONDecodeFailed   = "failed to decode JSON data"
	ErrInvalidDataSource  = "invalid data source"
	ErrDataSourceNotFound = "data source not found"
)

// Error messages for data sinks
const (
	ErrFileCreateFailed = "failed to create file"
	ErrFileWriteFailed  = "failed to write file"
	ErrJSONEncodeFailed = "failed to encode JSON data"
	ErrInvalidDataSink  = "invalid data sink"
	ErrDataSinkNotFound = "data sink not found"
)

// Error messages for processors
const (
	ErrUnknownProcessorType       = "unknown processor type"
	ErrProcessorCreationFailed    = "failed to create processor"
	ErrProcessorExecutionFailed   = "processor execution failed"
	ErrFilterFieldRequired        = "field is required for filter processor"
	ErrFilterOperationRequired    = "operation is required for filter processor"
	ErrTransformFieldRequired     = "field is required for transform processor"
	ErrTransformOperationRequired = "operation is required for transform processor"
	ErrAggregateGroupByRequired   = "group_by is required for aggregate processor"
	ErrInvalidProcessorConfig     = "invalid processor configuration"
	ErrProcessorNotFound          = "processor not found"
)

// Error messages for pipeline
const (
	ErrPipelineExecutionFailed = "pipeline execution failed"
	ErrPipelineEmpty           = "pipeline is empty"
	ErrPipelineProcessorError  = "error in processor"
)

// Error messages for data processor
const (
	ErrDataReadFailed    = "error reading data"
	ErrDataProcessFailed = "error processing data"
	ErrDataWriteFailed   = "error writing data"
	ErrDataProcessorInit = "failed to initialize data processor"
)

// Error messages for validation
const (
	ErrInvalidField      = "invalid field"
	ErrInvalidOperation  = "invalid operation"
	ErrInvalidValue      = "invalid value"
	ErrEmptyData         = "data is empty"
	ErrInvalidDataFormat = "invalid data format"
)

// Error messages for configuration
const (
	ErrConfigNotFound    = "configuration not found"
	ErrInvalidConfig     = "invalid configuration"
	ErrConfigParseFailed = "failed to parse configuration"
	ErrConfigRequired    = "configuration is required"
)
