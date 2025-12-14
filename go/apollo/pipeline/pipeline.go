package pipeline

import (
	"apollo_test/processors"
	"fmt"
)

// ============================================================================
// CHAIN OF RESPONSIBILITY PATTERN - Pipeline processing
// ============================================================================

// Pipeline chains multiple processors together
type Pipeline struct {
	processors []processors.Processor
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		processors: make([]processors.Processor, 0),
	}
}

func (p *Pipeline) AddProcessor(processor processors.Processor) *Pipeline {
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
