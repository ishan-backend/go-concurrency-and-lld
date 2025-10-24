package main

import (
	"fmt"
	"time"

	"github.com/ishanpandey/LowLevelDesign/go/kafka/pkg/kafka"
)

func testKafka() {
	// Create topic service
	topicService := kafka.NewTopicServiceImpl()
	topic := topicService.CreateTopic("my-topic", 3)

	if topic == nil {
		fmt.Println("Failed to create topic")
		return
	}

	fmt.Printf("Created topic: %s with %d partitions\n", topic.Name, topic.MaxPartitions)

	// Create producer
	producer := kafka.NewProducer()

	// Create consumer group
	consumerGroup := kafka.NewConsumerGroup()
	consumerService := kafka.NewConsumerServiceImpl()

	// Create worker manager
	workerManager := kafka.NewWorkerManager()

	// Create producer worker
	producerWorker := kafka.NewProducerWorker(topic, producer, topicService)
	workerManager.AddProducerWorker(producerWorker)

	// Create consumers and their workers
	consumer1 := kafka.NewConsumer(consumerGroup)
	consumerService.Subscribe(topic, consumer1)
	consumerWorker1 := kafka.NewConsumerWorker(consumer1, consumerService)
	workerManager.AddConsumerWorker(consumerWorker1)

	consumer2 := kafka.NewConsumer(consumerGroup)
	consumerService.Subscribe(topic, consumer2)
	consumerWorker2 := kafka.NewConsumerWorker(consumer2, consumerService)
	workerManager.AddConsumerWorker(consumerWorker2)

	consumer3 := kafka.NewConsumer(consumerGroup)
	consumerService.Subscribe(topic, consumer3)
	consumerWorker3 := kafka.NewConsumerWorker(consumer3, consumerService)
	workerManager.AddConsumerWorker(consumerWorker3)

	// Consumer 4 should remain idle as we have only 3 partitions in the topic
	consumer4 := kafka.NewConsumer(consumerGroup)
	partition4 := consumerService.Subscribe(topic, consumer4)
	if partition4 == nil {
		fmt.Printf("Consumer %s remains idle (no available partitions)\n", consumer4.String())
	}
	consumerWorker4 := kafka.NewConsumerWorker(consumer4, consumerService)
	workerManager.AddConsumerWorker(consumerWorker4)

	// Start all workers
	workerManager.StartAll()

	// Let the system run for some time
	fmt.Println("Kafka system running... Press Ctrl+C to stop")
	time.Sleep(30 * time.Second)

	// Stop all workers
	workerManager.StopAll()
	fmt.Println("Kafka system stopped")
}

func main() {
	fmt.Println("Starting Kafka-like Message Broker Demo")
	testKafka()
}
