package processors

import (
	"testing"
)

func TestCreateFilterProcessor(t *testing.T) {
	tests := []struct {
		name    string
		config  map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid config with equals operation",
			config: map[string]interface{}{
				"field":     "status",
				"operation": "equals",
				"value":     "active",
			},
			wantErr: false,
		},
		{
			name: "valid config with contains operation",
			config: map[string]interface{}{
				"field":     "name",
				"operation": "contains",
				"value":     "test",
			},
			wantErr: false,
		},
		{
			name: "valid config with not_empty operation",
			config: map[string]interface{}{
				"field":     "email",
				"operation": "not_empty",
				"value":     nil,
			},
			wantErr: false,
		},
		{
			name: "missing field",
			config: map[string]interface{}{
				"operation": "equals",
				"value":     "test",
			},
			wantErr: true,
		},
		{
			name: "missing operation",
			config: map[string]interface{}{
				"field": "status",
				"value": "active",
			},
			wantErr: true,
		},
		{
			name: "field is not string",
			config: map[string]interface{}{
				"field":     123,
				"operation": "equals",
				"value":     "test",
			},
			wantErr: true,
		},
		{
			name: "operation is not string",
			config: map[string]interface{}{
				"field":     "status",
				"operation": 123,
				"value":     "test",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor, err := CreateFilterProcessor(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateFilterProcessor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && processor == nil {
				t.Error("CreateFilterProcessor() returned nil processor without error")
			}
		})
	}
}

func TestFilterProcessor_Process(t *testing.T) {
	tests := []struct {
		name      string
		processor *FilterProcessor
		data      []map[string]interface{}
		want      []map[string]interface{}
		wantErr   bool
	}{
		{
			name: "filter equals - single match",
			processor: &FilterProcessor{
				field:     "status",
				operation: "equals",
				value:     "active",
			},
			data: []map[string]interface{}{
				{"id": 1, "status": "active", "name": "test1"},
				{"id": 2, "status": "inactive", "name": "test2"},
				{"id": 3, "status": "active", "name": "test3"},
			},
			want: []map[string]interface{}{
				{"id": 1, "status": "active", "name": "test1"},
				{"id": 3, "status": "active", "name": "test3"},
			},
			wantErr: false,
		},
		{
			name: "filter equals - no matches",
			processor: &FilterProcessor{
				field:     "status",
				operation: "equals",
				value:     "pending",
			},
			data: []map[string]interface{}{
				{"id": 1, "status": "active", "name": "test1"},
				{"id": 2, "status": "inactive", "name": "test2"},
			},
			want:    []map[string]interface{}{},
			wantErr: false,
		},
		{
			name: "filter contains - multiple matches",
			processor: &FilterProcessor{
				field:     "name",
				operation: "contains",
				value:     "test",
			},
			data: []map[string]interface{}{
				{"id": 1, "name": "test1"},
				{"id": 2, "name": "example"},
				{"id": 3, "name": "test2"},
			},
			want: []map[string]interface{}{
				{"id": 1, "name": "test1"},
				{"id": 3, "name": "test2"},
			},
			wantErr: false,
		},
		{
			name: "filter not_empty - filters out empty and nil",
			processor: &FilterProcessor{
				field:     "email",
				operation: "not_empty",
				value:     nil,
			},
			data: []map[string]interface{}{
				{"id": 1, "email": "test@example.com"},
				{"id": 2, "email": ""},
				{"id": 3, "email": nil},
				{"id": 4, "email": "another@example.com"},
			},
			want: []map[string]interface{}{
				{"id": 1, "email": "test@example.com"},
				{"id": 4, "email": "another@example.com"},
			},
			wantErr: false,
		},
		{
			name: "filter - field does not exist",
			processor: &FilterProcessor{
				field:     "missing",
				operation: "equals",
				value:     "test",
			},
			data: []map[string]interface{}{
				{"id": 1, "name": "test1"},
				{"id": 2, "name": "test2"},
			},
			want:    []map[string]interface{}{},
			wantErr: false,
		},
		{
			name: "filter - empty data",
			processor: &FilterProcessor{
				field:     "status",
				operation: "equals",
				value:     "active",
			},
			data:    []map[string]interface{}{},
			want:    []map[string]interface{}{},
			wantErr: false,
		},
		{
			name: "filter - unknown operation",
			processor: &FilterProcessor{
				field:     "status",
				operation: "unknown",
				value:     "active",
			},
			data: []map[string]interface{}{
				{"id": 1, "status": "active"},
			},
			want:    []map[string]interface{}{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.processor.Process(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterProcessor.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("FilterProcessor.Process() returned %d items, want %d", len(got), len(tt.want))
				return
			}
			// Compare results (simplified comparison)
			for i, item := range got {
				if i >= len(tt.want) {
					t.Errorf("FilterProcessor.Process() returned more items than expected")
					break
				}
				// Compare key fields
				for key, val := range item {
					if wantVal, exists := tt.want[i][key]; exists {
						if val != wantVal {
							t.Errorf("FilterProcessor.Process() item[%d][%s] = %v, want %v", i, key, val, wantVal)
						}
					}
				}
			}
		})
	}
}

func TestFilterProcessor_matchesCondition(t *testing.T) {
	tests := []struct {
		name      string
		processor *FilterProcessor
		row       map[string]interface{}
		want      bool
	}{
		{
			name: "equals - match",
			processor: &FilterProcessor{
				field:     "status",
				operation: "equals",
				value:     "active",
			},
			row:  map[string]interface{}{"status": "active"},
			want: true,
		},
		{
			name: "equals - no match",
			processor: &FilterProcessor{
				field:     "status",
				operation: "equals",
				value:     "active",
			},
			row:  map[string]interface{}{"status": "inactive"},
			want: false,
		},
		{
			name: "equals - numeric comparison",
			processor: &FilterProcessor{
				field:     "id",
				operation: "equals",
				value:     123,
			},
			row:  map[string]interface{}{"id": 123},
			want: true,
		},
		{
			name: "contains - match",
			processor: &FilterProcessor{
				field:     "name",
				operation: "contains",
				value:     "test",
			},
			row:  map[string]interface{}{"name": "test123"},
			want: true,
		},
		{
			name: "contains - no match",
			processor: &FilterProcessor{
				field:     "name",
				operation: "contains",
				value:     "test",
			},
			row:  map[string]interface{}{"name": "example"},
			want: false,
		},
		{
			name: "not_empty - has value",
			processor: &FilterProcessor{
				field:     "email",
				operation: "not_empty",
				value:     nil,
			},
			row:  map[string]interface{}{"email": "test@example.com"},
			want: true,
		},
		{
			name: "not_empty - empty string",
			processor: &FilterProcessor{
				field:     "email",
				operation: "not_empty",
				value:     nil,
			},
			row:  map[string]interface{}{"email": ""},
			want: false,
		},
		{
			name: "not_empty - nil value",
			processor: &FilterProcessor{
				field:     "email",
				operation: "not_empty",
				value:     nil,
			},
			row:  map[string]interface{}{"email": nil},
			want: false,
		},
		{
			name: "field does not exist",
			processor: &FilterProcessor{
				field:     "missing",
				operation: "equals",
				value:     "test",
			},
			row:  map[string]interface{}{"other": "value"},
			want: false,
		},
		{
			name: "unknown operation",
			processor: &FilterProcessor{
				field:     "status",
				operation: "unknown",
				value:     "active",
			},
			row:  map[string]interface{}{"status": "active"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.processor.matchesCondition(tt.row)
			if got != tt.want {
				t.Errorf("FilterProcessor.matchesCondition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterProcessor_Name(t *testing.T) {
	tests := []struct {
		name      string
		processor *FilterProcessor
		want      string
	}{
		{
			name: "basic name",
			processor: &FilterProcessor{
				field:     "status",
				operation: "equals",
				value:     "active",
			},
			want: "Filter[status equals active]",
		},
		{
			name: "contains operation",
			processor: &FilterProcessor{
				field:     "name",
				operation: "contains",
				value:     "test",
			},
			want: "Filter[name contains test]",
		},
		{
			name: "numeric value",
			processor: &FilterProcessor{
				field:     "id",
				operation: "equals",
				value:     123,
			},
			want: "Filter[id equals 123]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.processor.Name()
			if got != tt.want {
				t.Errorf("FilterProcessor.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}
