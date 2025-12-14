package processors

import "fmt"

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
