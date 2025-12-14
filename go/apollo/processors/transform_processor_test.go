package processors

import (
	"testing"
)

func TestCreateTransformProcessor(t *testing.T) {
	tests := []struct {
		name    string
		config  map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid config with uppercase operation",
			config: map[string]interface{}{
				"field":     "name",
				"operation": "uppercase",
			},
			wantErr: false,
		},
		{
			name: "valid config with lowercase operation",
			config: map[string]interface{}{
				"field":     "name",
				"operation": "lowercase",
			},
			wantErr: false,
		},
		{
			name: "valid config with trim operation",
			config: map[string]interface{}{
				"field":     "name",
				"operation": "trim",
			},
			wantErr: false,
		},
		{
			name: "missing field",
			config: map[string]interface{}{
				"operation": "uppercase",
			},
			wantErr: true,
		},
		{
			name: "missing operation",
			config: map[string]interface{}{
				"field": "name",
			},
			wantErr: true,
		},
		{
			name: "field is not string",
			config: map[string]interface{}{
				"field":     123,
				"operation": "uppercase",
			},
			wantErr: true,
		},
		{
			name: "operation is not string",
			config: map[string]interface{}{
				"field":     "name",
				"operation": 123,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor, err := CreateTransformProcessor(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTransformProcessor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && processor == nil {
				t.Error("CreateTransformProcessor() returned nil processor without error")
			}
		})
	}
}

func TestTransformProcessor_Process(t *testing.T) {
	tests := []struct {
		name      string
		processor *TransformProcessor
		data      []map[string]interface{}
		want      []map[string]interface{}
		wantErr   bool
	}{
		{
			name: "uppercase transformation",
			processor: &TransformProcessor{
				field:     "name",
				operation: "uppercase",
			},
			data: []map[string]interface{}{
				{"id": 1, "name": "john"},
				{"id": 2, "name": "jane"},
			},
			want: []map[string]interface{}{
				{"id": 1, "name": "JOHN"},
				{"id": 2, "name": "JANE"},
			},
			wantErr: false,
		},
		{
			name: "lowercase transformation",
			processor: &TransformProcessor{
				field:     "name",
				operation: "lowercase",
			},
			data: []map[string]interface{}{
				{"id": 1, "name": "JOHN"},
				{"id": 2, "name": "JANE"},
			},
			want: []map[string]interface{}{
				{"id": 1, "name": "john"},
				{"id": 2, "name": "jane"},
			},
			wantErr: false,
		},
		{
			name: "trim transformation",
			processor: &TransformProcessor{
				field:     "name",
				operation: "trim",
			},
			data: []map[string]interface{}{
				{"id": 1, "name": "  john  "},
				{"id": 2, "name": "\t\njane\t\n"},
			},
			want: []map[string]interface{}{
				{"id": 1, "name": "john"},
				{"id": 2, "name": "jane"},
			},
			wantErr: false,
		},
		{
			name: "field does not exist - no transformation",
			processor: &TransformProcessor{
				field:     "missing",
				operation: "uppercase",
			},
			data: []map[string]interface{}{
				{"id": 1, "name": "john"},
			},
			want: []map[string]interface{}{
				{"id": 1, "name": "john"},
			},
			wantErr: false,
		},
		{
			name: "unknown operation - no transformation",
			processor: &TransformProcessor{
				field:     "name",
				operation: "unknown",
			},
			data: []map[string]interface{}{
				{"id": 1, "name": "john"},
			},
			want: []map[string]interface{}{
				{"id": 1, "name": "john"},
			},
			wantErr: false,
		},
		{
			name: "empty data",
			processor: &TransformProcessor{
				field:     "name",
				operation: "uppercase",
			},
			data:    []map[string]interface{}{},
			want:    []map[string]interface{}{},
			wantErr: false,
		},
		{
			name: "preserves other fields",
			processor: &TransformProcessor{
				field:     "name",
				operation: "uppercase",
			},
			data: []map[string]interface{}{
				{"id": 1, "name": "john", "email": "john@example.com", "age": 30},
			},
			want: []map[string]interface{}{
				{"id": 1, "name": "JOHN", "email": "john@example.com", "age": 30},
			},
			wantErr: false,
		},
		{
			name: "numeric value transformation",
			processor: &TransformProcessor{
				field:     "id",
				operation: "uppercase",
			},
			data: []map[string]interface{}{
				{"id": 123, "name": "john"},
			},
			want: []map[string]interface{}{
				{"id": "123", "name": "john"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.processor.Process(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransformProcessor.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("TransformProcessor.Process() returned %d items, want %d", len(got), len(tt.want))
				return
			}
			// Compare results
			for i, item := range got {
				if i >= len(tt.want) {
					t.Errorf("TransformProcessor.Process() returned more items than expected")
					break
				}
				// Compare all fields
				for key, val := range item {
					if wantVal, exists := tt.want[i][key]; exists {
						if val != wantVal {
							t.Errorf("TransformProcessor.Process() item[%d][%s] = %v, want %v", i, key, val, wantVal)
						}
					} else {
						t.Errorf("TransformProcessor.Process() item[%d] has unexpected key: %s", i, key)
					}
				}
				// Check for missing keys
				for key := range tt.want[i] {
					if _, exists := item[key]; !exists {
						t.Errorf("TransformProcessor.Process() item[%d] missing key: %s", i, key)
					}
				}
			}
		})
	}
}

func TestTransformProcessor_transform(t *testing.T) {
	tests := []struct {
		name      string
		processor *TransformProcessor
		value     interface{}
		want      interface{}
	}{
		{
			name: "uppercase - string",
			processor: &TransformProcessor{
				field:     "name",
				operation: "uppercase",
			},
			value: "hello",
			want:  "HELLO",
		},
		{
			name: "uppercase - numeric",
			processor: &TransformProcessor{
				field:     "id",
				operation: "uppercase",
			},
			value: 123,
			want:  "123",
		},
		{
			name: "lowercase - string",
			processor: &TransformProcessor{
				field:     "name",
				operation: "lowercase",
			},
			value: "HELLO",
			want:  "hello",
		},
		{
			name: "trim - with spaces",
			processor: &TransformProcessor{
				field:     "name",
				operation: "trim",
			},
			value: "  hello  ",
			want:  "hello",
		},
		{
			name: "trim - with tabs and newlines",
			processor: &TransformProcessor{
				field:     "name",
				operation: "trim",
			},
			value: "\t\nhello\t\n",
			want:  "hello",
		},
		{
			name: "unknown operation - returns original",
			processor: &TransformProcessor{
				field:     "name",
				operation: "unknown",
			},
			value: "hello",
			want:  "hello",
		},
		{
			name: "trim - already trimmed",
			processor: &TransformProcessor{
				field:     "name",
				operation: "trim",
			},
			value: "hello",
			want:  "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.processor.transform(tt.value)
			if got != tt.want {
				t.Errorf("TransformProcessor.transform() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransformProcessor_Name(t *testing.T) {
	tests := []struct {
		name      string
		processor *TransformProcessor
		want      string
	}{
		{
			name: "uppercase operation",
			processor: &TransformProcessor{
				field:     "name",
				operation: "uppercase",
			},
			want: "Transform[name: uppercase]",
		},
		{
			name: "lowercase operation",
			processor: &TransformProcessor{
				field:     "email",
				operation: "lowercase",
			},
			want: "Transform[email: lowercase]",
		},
		{
			name: "trim operation",
			processor: &TransformProcessor{
				field:     "description",
				operation: "trim",
			},
			want: "Transform[description: trim]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.processor.Name()
			if got != tt.want {
				t.Errorf("TransformProcessor.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}
