package processors

// ============================================================================
// STRATEGY PATTERN - Define processing strategies
// ============================================================================

// Processor defines the interface for data processing steps
type Processor interface {
	Process(data []map[string]interface{}) ([]map[string]interface{}, error)
	Name() string
}
