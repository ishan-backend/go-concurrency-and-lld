package processors

import (
	"fmt"
	"strings"
)

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
