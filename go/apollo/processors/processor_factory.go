package processors

import "fmt"

type ProcessorFactory struct {
	creators map[string]func(config map[string]interface{}) (Processor, error)
}

func NewProcessorFactory() *ProcessorFactory {
	factory := &ProcessorFactory{
		creators: make(map[string]func(config map[string]interface{}) (Processor, error)),
	}

	// Register default processors
	factory.Register("filter", CreateFilterProcessor)
	factory.Register("transform", CreateTransformProcessor)
	factory.Register("aggregate", CreateAggregateProcessor)

	return factory
}

func (f *ProcessorFactory) Register(name string, creator func(config map[string]interface{}) (Processor, error)) {
	f.creators[name] = creator
}

func (f *ProcessorFactory) Create(name string, config map[string]interface{}) (Processor, error) {
	creator, exists := f.creators[name]
	if !exists {
		return nil, fmt.Errorf("unknown processor type: %s", name)
	}
	return creator(config)
}
