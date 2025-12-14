package pipeline

import (
	"apollo_test/processors"
	"fmt"
)

type PipelineBuilder struct {
	pipeline *Pipeline
	factory  *processors.ProcessorFactory
}

func NewPipelineBuilder() *PipelineBuilder {
	return &PipelineBuilder{
		pipeline: NewPipeline(),
		factory:  processors.NewProcessorFactory(),
	}
}

func (pb *PipelineBuilder) Filter(field, operation string, value interface{}) *PipelineBuilder {
	config := map[string]interface{}{
		"field":     field,
		"operation": operation,
		"value":     value,
	}
	processor, err := pb.factory.Create("filter", config)
	if err != nil {
		fmt.Printf("Error creating filter: %v\n", err)
		return pb
	}
	pb.pipeline.AddProcessor(processor)
	return pb
}

func (pb *PipelineBuilder) Transform(field, operation string) *PipelineBuilder {
	config := map[string]interface{}{
		"field":     field,
		"operation": operation,
	}
	processor, err := pb.factory.Create("transform", config)
	if err != nil {
		fmt.Printf("Error creating transform: %v\n", err)
		return pb
	}
	pb.pipeline.AddProcessor(processor)
	return pb
}

func (pb *PipelineBuilder) Aggregate(groupBy string) *PipelineBuilder {
	config := map[string]interface{}{
		"group_by": groupBy,
	}
	processor, err := pb.factory.Create("aggregate", config)
	if err != nil {
		fmt.Printf("Error creating aggregate: %v\n", err)
		return pb
	}
	pb.pipeline.AddProcessor(processor)
	return pb
}

func (pb *PipelineBuilder) Build() *Pipeline {
	return pb.pipeline
}
