package kafka

import (
	"fmt"
	"math/rand"
	"time"
)

// ProducerWorker represents a worker that produces messages
type ProducerWorker struct {
	topic       *Topic
	producer    *Producer
	topicService TopicService
	stopChan    chan bool
}

// NewProducerWorker creates a new producer worker
func NewProducerWorker(topic *Topic, producer *Producer, topicService TopicService) *ProducerWorker {
	return &ProducerWorker{
		topic:       topic,
		producer:    producer,
		topicService: topicService,
		stopChan:    make(chan bool),
	}
}

// Start starts the producer worker in a goroutine
func (pw *ProducerWorker) Start() {
	go pw.run()
}

// Stop stops the producer worker
func (pw *ProducerWorker) Stop() {
	close(pw.stopChan)
}

// run is the main loop for the producer worker
func (pw *ProducerWorker) run() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-pw.stopChan:
			return
		case <-ticker.C:
			pw.produceMessage()
		}
	}
}

// produceMessage produces a random message
func (pw *ProducerWorker) produceMessage() {
	message := NewMessage(fmt.Sprintf("random message %d", rand.Intn(100000)))
	fmt.Printf("[Producer] %s producing message = %s to topic %s\n", 
		pw.producer.String(), message.String(), pw.topic.Name)
	
	messageMetadata := pw.topicService.PublishMessage(pw.topic, message)
	if messageMetadata != nil {
		fmt.Printf("Message published to: %s\n", messageMetadata.String())
	}
}

// ConsumerWorker represents a worker that consumes messages
type ConsumerWorker struct {
	consumer        *Consumer
	consumerService ConsumerService
	stopChan        chan bool
}

// NewConsumerWorker creates a new consumer worker
func NewConsumerWorker(consumer *Consumer, consumerService ConsumerService) *ConsumerWorker {
	return &ConsumerWorker{
		consumer:        consumer,
		consumerService: consumerService,
		stopChan:        make(chan bool),
	}
}

// Start starts the consumer worker in a goroutine
func (cw *ConsumerWorker) Start() {
	go cw.run()
}

// Stop stops the consumer worker
func (cw *ConsumerWorker) Stop() {
	close(cw.stopChan)
}

// run is the main loop for the consumer worker
func (cw *ConsumerWorker) run() {
	ticker := time.NewTicker(1200 * time.Millisecond)
	defer ticker.Stop()
	
	for {
		select {
		case <-cw.stopChan:
			return
		case <-ticker.C:
			cw.consumeMessage()
		}
	}
}

// consumeMessage consumes a message
func (cw *ConsumerWorker) consumeMessage() {
	message := cw.consumerService.Consume(cw.consumer)
	if message != nil {
		fmt.Printf("[Consumer] %s consumed message = %s\n", 
			cw.consumer.String(), message.String())
	}
}

// WorkerManager manages multiple workers
type WorkerManager struct {
	producerWorkers []*ProducerWorker
	consumerWorkers []*ConsumerWorker
}

// NewWorkerManager creates a new worker manager
func NewWorkerManager() *WorkerManager {
	return &WorkerManager{
		producerWorkers: make([]*ProducerWorker, 0),
		consumerWorkers: make([]*ConsumerWorker, 0),
	}
}

// AddProducerWorker adds a producer worker to the manager
func (wm *WorkerManager) AddProducerWorker(worker *ProducerWorker) {
	wm.producerWorkers = append(wm.producerWorkers, worker)
}

// AddConsumerWorker adds a consumer worker to the manager
func (wm *WorkerManager) AddConsumerWorker(worker *ConsumerWorker) {
	wm.consumerWorkers = append(wm.consumerWorkers, worker)
}

// StartAll starts all workers
func (wm *WorkerManager) StartAll() {
	for _, worker := range wm.producerWorkers {
		worker.Start()
	}
	
	for _, worker := range wm.consumerWorkers {
		worker.Start()
	}
}

// StopAll stops all workers
func (wm *WorkerManager) StopAll() {
	for _, worker := range wm.producerWorkers {
		worker.Stop()
	}
	
	for _, worker := range wm.consumerWorkers {
		worker.Stop()
	}
}
