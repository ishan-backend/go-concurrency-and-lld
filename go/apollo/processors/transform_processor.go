package processors

import (
	"fmt"
	"strings"
)

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
