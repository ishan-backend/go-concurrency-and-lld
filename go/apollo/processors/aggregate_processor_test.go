package processors

import (
	"fmt"
	"testing"
)

func TestCreateAggregateProcessor(t *testing.T) {
	tests := []struct {
		name    string
		config  map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid config",
			config: map[string]interface{}{
				"group_by": "status",
			},
			wantErr: false,
		},
		{
			name: "missing group_by",
			config: map[string]interface{}{
				"field": "status",
			},
			wantErr: true,
		},
		{
			name: "group_by is not string",
			config: map[string]interface{}{
				"group_by": 123,
			},
			wantErr: true,
		},
		{
			name:    "empty config",
			config:  map[string]interface{}{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor, err := CreateAggregateProcessor(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAggregateProcessor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && processor == nil {
				t.Error("CreateAggregateProcessor() returned nil processor without error")
			}
		})
	}
}

func TestAggregateProcessor_Process(t *testing.T) {
	tests := []struct {
		name      string
		processor *AggregateProcessor
		data      []map[string]interface{}
		want      map[string]int // map of group value to count
		wantErr   bool
	}{
		{
			name: "simple aggregation",
			processor: &AggregateProcessor{
				groupBy: "status",
			},
			data: []map[string]interface{}{
				{"id": 1, "status": "active"},
				{"id": 2, "status": "inactive"},
				{"id": 3, "status": "active"},
				{"id": 4, "status": "active"},
			},
			want: map[string]int{
				"active":   3,
				"inactive": 1,
			},
			wantErr: false,
		},
		{
			name: "single group",
			processor: &AggregateProcessor{
				groupBy: "status",
			},
			data: []map[string]interface{}{
				{"id": 1, "status": "active"},
				{"id": 2, "status": "active"},
				{"id": 3, "status": "active"},
			},
			want: map[string]int{
				"active": 3,
			},
			wantErr: false,
		},
		{
			name: "field does not exist - skips row",
			processor: &AggregateProcessor{
				groupBy: "missing",
			},
			data: []map[string]interface{}{
				{"id": 1, "status": "active"},
				{"id": 2, "status": "inactive"},
			},
			want:    map[string]int{},
			wantErr: false,
		},
		{
			name: "empty data",
			processor: &AggregateProcessor{
				groupBy: "status",
			},
			data:    []map[string]interface{}{},
			want:    map[string]int{},
			wantErr: false,
		},
		{
			name: "numeric grouping",
			processor: &AggregateProcessor{
				groupBy: "age",
			},
			data: []map[string]interface{}{
				{"id": 1, "age": 25},
				{"id": 2, "age": 30},
				{"id": 3, "age": 25},
				{"id": 4, "age": 30},
				{"id": 5, "age": 30},
			},
			want: map[string]int{
				"25": 2,
				"30": 3,
			},
			wantErr: false,
		},
		{
			name: "mixed types in group_by field",
			processor: &AggregateProcessor{
				groupBy: "category",
			},
			data: []map[string]interface{}{
				{"id": 1, "category": "A"},
				{"id": 2, "category": 123},
				{"id": 3, "category": "A"},
				{"id": 4, "category": "B"},
			},
			want: map[string]int{
				"A":   2,
				"123": 1,
				"B":   1,
			},
			wantErr: false,
		},
		{
			name: "nil values in group_by field - skips",
			processor: &AggregateProcessor{
				groupBy: "status",
			},
			data: []map[string]interface{}{
				{"id": 1, "status": "active"},
				{"id": 2, "status": nil},
				{"id": 3, "status": "active"},
			},
			want: map[string]int{
				"active": 2,
				"<nil>":  1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.processor.Process(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("AggregateProcessor.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Convert result to map for easier comparison
			resultMap := make(map[string]int)
			for _, item := range got {
				groupValue := fmt.Sprintf("%v", item[tt.processor.groupBy])
				count, ok := item["count"].(int)
				if !ok {
					t.Errorf("AggregateProcessor.Process() result item missing 'count' field or wrong type")
					continue
				}
				resultMap[groupValue] = count
			}

			// Compare counts
			if len(resultMap) != len(tt.want) {
				t.Errorf("AggregateProcessor.Process() returned %d groups, want %d", len(resultMap), len(tt.want))
				t.Errorf("Got: %v, Want: %v", resultMap, tt.want)
				return
			}

			for key, wantCount := range tt.want {
				gotCount, exists := resultMap[key]
				if !exists {
					t.Errorf("AggregateProcessor.Process() missing group: %s", key)
					continue
				}
				if gotCount != wantCount {
					t.Errorf("AggregateProcessor.Process() group[%s] count = %d, want %d", key, gotCount, wantCount)
				}
			}
		})
	}
}

func TestAggregateProcessor_Process_ResultStructure(t *testing.T) {
	tests := []struct {
		name      string
		processor *AggregateProcessor
		data      []map[string]interface{}
		wantErr   bool
	}{
		{
			name: "verify result structure with status field",
			processor: &AggregateProcessor{
				groupBy: "status",
			},
			data: []map[string]interface{}{
				{"id": 1, "status": "active"},
				{"id": 2, "status": "inactive"},
				{"id": 3, "status": "active"},
			},
			wantErr: false,
		},
		{
			name: "verify result structure with category field",
			processor: &AggregateProcessor{
				groupBy: "category",
			},
			data: []map[string]interface{}{
				{"id": 1, "category": "A"},
				{"id": 2, "category": "B"},
				{"id": 3, "category": "A"},
			},
			wantErr: false,
		},
		{
			name: "verify result structure with empty data",
			processor: &AggregateProcessor{
				groupBy: "status",
			},
			data:    []map[string]interface{}{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.processor.Process(tt.data)
			if (err != nil) != tt.wantErr {
				t.Fatalf("AggregateProcessor.Process() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify result structure
			for _, item := range result {
				// Check that each item has group_by field
				if _, exists := item[tt.processor.groupBy]; !exists {
					t.Errorf("AggregateProcessor.Process() result item missing group_by field: %s", tt.processor.groupBy)
				}
				// Check that each item has count field
				if _, exists := item["count"]; !exists {
					t.Errorf("AggregateProcessor.Process() result item missing 'count' field")
				}
				// Check that count is an integer
				if count, ok := item["count"].(int); !ok {
					t.Errorf("AggregateProcessor.Process() result item 'count' is not int, got %T", count)
				}
			}
		})
	}
}

func TestAggregateProcessor_Name(t *testing.T) {
	tests := []struct {
		name      string
		processor *AggregateProcessor
		want      string
	}{
		{
			name: "basic name",
			processor: &AggregateProcessor{
				groupBy: "status",
			},
			want: "Aggregate[group_by: status]",
		},
		{
			name: "different field",
			processor: &AggregateProcessor{
				groupBy: "category",
			},
			want: "Aggregate[group_by: category]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.processor.Name()
			if got != tt.want {
				t.Errorf("AggregateProcessor.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}
