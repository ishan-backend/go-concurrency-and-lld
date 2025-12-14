package processors

import (
	"testing"
)

func TestNewProcessorFactory(t *testing.T) {
	tests := []struct {
		name             string
		wantDefaultProcs []string
		checkFactoryNil  bool
		checkCreatorsNil bool
	}{
		{
			name:             "factory initialization",
			wantDefaultProcs: []string{"filter", "transform", "aggregate"},
			checkFactoryNil:  false,
			checkCreatorsNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewProcessorFactory()

			if factory == nil {
				t.Fatal("NewProcessorFactory() returned nil")
			}

			if factory.creators == nil {
				t.Fatal("NewProcessorFactory() creators map is nil")
			}

			// Verify default processors are registered
			for _, procType := range tt.wantDefaultProcs {
				if _, exists := factory.creators[procType]; !exists {
					t.Errorf("NewProcessorFactory() missing default processor: %s", procType)
				}
			}
		})
	}
}

func TestProcessorFactory_Register(t *testing.T) {
	tests := []struct {
		name          string
		processorName string
		creator       func(config map[string]interface{}) (Processor, error)
		wantErr       bool
	}{
		{
			name:          "register custom processor",
			processorName: "custom",
			creator: func(config map[string]interface{}) (Processor, error) {
				return &FilterProcessor{
					field:     "test",
					operation: "equals",
					value:     "value",
				}, nil
			},
			wantErr: false,
		},
		{
			name:          "register another custom processor",
			processorName: "custom2",
			creator: func(config map[string]interface{}) (Processor, error) {
				return &TransformProcessor{
					field:     "name",
					operation: "uppercase",
				}, nil
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewProcessorFactory()

			factory.Register(tt.processorName, tt.creator)

			if _, exists := factory.creators[tt.processorName]; !exists {
				t.Errorf("ProcessorFactory.Register() did not register processor: %s", tt.processorName)
			}

			// Verify it can be created
			processor, err := factory.Create(tt.processorName, map[string]interface{}{})
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessorFactory.Register() registered processor failed to create: %v", err)
			}
			if !tt.wantErr && processor == nil {
				t.Error("ProcessorFactory.Register() registered processor returned nil")
			}
		})
	}
}

func TestProcessorFactory_Create(t *testing.T) {
	factory := NewProcessorFactory()

	tests := []struct {
		name     string
		procType string
		config   map[string]interface{}
		wantErr  bool
	}{
		{
			name:     "create filter processor",
			procType: "filter",
			config: map[string]interface{}{
				"field":     "status",
				"operation": "equals",
				"value":     "active",
			},
			wantErr: false,
		},
		{
			name:     "create transform processor",
			procType: "transform",
			config: map[string]interface{}{
				"field":     "name",
				"operation": "uppercase",
			},
			wantErr: false,
		},
		{
			name:     "create aggregate processor",
			procType: "aggregate",
			config: map[string]interface{}{
				"group_by": "status",
			},
			wantErr: false,
		},
		{
			name:     "unknown processor type",
			procType: "unknown",
			config:   map[string]interface{}{},
			wantErr:  true,
		},
		{
			name:     "filter processor with invalid config",
			procType: "filter",
			config:   map[string]interface{}{
				// missing required fields
			},
			wantErr: true,
		},
		{
			name:     "transform processor with invalid config",
			procType: "transform",
			config:   map[string]interface{}{
				// missing required fields
			},
			wantErr: true,
		},
		{
			name:     "aggregate processor with invalid config",
			procType: "aggregate",
			config:   map[string]interface{}{
				// missing required fields
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor, err := factory.Create(tt.procType, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessorFactory.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && processor == nil {
				t.Error("ProcessorFactory.Create() returned nil processor without error")
			}
			if !tt.wantErr && processor != nil {
				// Verify processor has Name method
				name := processor.Name()
				if name == "" {
					t.Error("ProcessorFactory.Create() returned processor with empty name")
				}
			}
		})
	}
}

// func TestProcessorFactory_TestSomething(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		config map[string]interface{}
// 		wantErr bool
// 	}{
// 		{

// 		}
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 		})
// 	}
// }

func TestProcessorFactory_Create_Integration(t *testing.T) {
	tests := []struct {
		name       string
		procType   string
		config     map[string]interface{}
		testData   []map[string]interface{}
		wantLen    int
		wantField  string
		wantValue  interface{}
		checkField string
		wantErr    bool
	}{
		{
			name:     "filter processor integration",
			procType: "filter",
			config: map[string]interface{}{
				"field":     "status",
				"operation": "equals",
				"value":     "active",
			},
			testData: []map[string]interface{}{
				{"id": 1, "status": "active"},
				{"id": 2, "status": "inactive"},
				{"id": 3, "status": "active"},
			},
			wantLen:    2,
			wantField:  "",
			wantValue:  nil,
			checkField: "",
			wantErr:    false,
		},
		{
			name:     "transform processor integration",
			procType: "transform",
			config: map[string]interface{}{
				"field":     "name",
				"operation": "uppercase",
			},
			testData: []map[string]interface{}{
				{"id": 1, "name": "john"},
			},
			wantLen:    1,
			wantField:  "name",
			wantValue:  "JOHN",
			checkField: "name",
			wantErr:    false,
		},
		{
			name:     "aggregate processor integration",
			procType: "aggregate",
			config: map[string]interface{}{
				"group_by": "status",
			},
			testData: []map[string]interface{}{
				{"id": 1, "status": "active"},
				{"id": 2, "status": "inactive"},
				{"id": 3, "status": "active"},
			},
			wantLen:    2,
			wantField:  "",
			wantValue:  nil,
			checkField: "",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewProcessorFactory()

			processor, err := factory.Create(tt.procType, tt.config)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ProcessorFactory.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			result, err := processor.Process(tt.testData)
			if err != nil {
				t.Fatalf("Processor.Process() error = %v", err)
			}

			if len(result) != tt.wantLen {
				t.Errorf("Processor.Process() returned %d items, want %d", len(result), tt.wantLen)
			}

			if tt.checkField != "" && len(result) > 0 {
				if result[0][tt.checkField] != tt.wantValue {
					t.Errorf("Processor.Process() result[0][%s] = %v, want %v", tt.checkField, result[0][tt.checkField], tt.wantValue)
				}
			}
		})
	}
}

func TestProcessorFactory_Register_Overwrite(t *testing.T) {
	tests := []struct {
		name          string
		processorName string
		creator       func(config map[string]interface{}) (Processor, error)
		wantField     string
		wantOperation string
		wantValue     interface{}
		wantErr       bool
	}{
		{
			name:          "overwrite filter processor",
			processorName: "filter",
			creator: func(config map[string]interface{}) (Processor, error) {
				return &FilterProcessor{
					field:     "custom",
					operation: "equals",
					value:     "custom",
				}, nil
			},
			wantField:     "custom",
			wantOperation: "equals",
			wantValue:     "custom",
			wantErr:       false,
		},
		{
			name:          "overwrite transform processor",
			processorName: "transform",
			creator: func(config map[string]interface{}) (Processor, error) {
				return &TransformProcessor{
					field:     "custom_field",
					operation: "lowercase",
				}, nil
			},
			wantField:     "custom_field",
			wantOperation: "lowercase",
			wantValue:     nil,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewProcessorFactory()

			factory.Register(tt.processorName, tt.creator)

			// Verify it uses the new creator
			processor, err := factory.Create(tt.processorName, map[string]interface{}{})
			if (err != nil) != tt.wantErr {
				t.Fatalf("ProcessorFactory.Create() after overwrite error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			// Check filter processor
			if filterProc, ok := processor.(*FilterProcessor); ok {
				if filterProc.field != tt.wantField {
					t.Errorf("ProcessorFactory.Register() overwrite failed, field = %v, want %v", filterProc.field, tt.wantField)
				}
				if filterProc.operation != tt.wantOperation {
					t.Errorf("ProcessorFactory.Register() overwrite failed, operation = %v, want %v", filterProc.operation, tt.wantOperation)
				}
			}

			// Check transform processor
			if transformProc, ok := processor.(*TransformProcessor); ok {
				if transformProc.field != tt.wantField {
					t.Errorf("ProcessorFactory.Register() overwrite failed, field = %v, want %v", transformProc.field, tt.wantField)
				}
				if transformProc.operation != tt.wantOperation {
					t.Errorf("ProcessorFactory.Register() overwrite failed, operation = %v, want %v", transformProc.operation, tt.wantOperation)
				}
			}
		})
	}
}
