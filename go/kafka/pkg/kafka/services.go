package kafka

import (
	"math/rand"
	"sync"
)

// Repository pattern for managing topics
type TopicRepository struct {
	topics map[string]*Topic
	mutex  sync.RWMutex
}

// NewTopicRepository creates a new topic repository
func NewTopicRepository() *TopicRepository {
	return &TopicRepository{
		topics: make(map[string]*Topic),
	}
}

// AddTopic adds a topic to the repository
func (tr *TopicRepository) AddTopic(topic *Topic) *Topic {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	if _, exists := tr.topics[topic.Name]; exists {
		return nil // Topic already exists
	}

	tr.topics[topic.Name] = topic
	return topic
}

// GetTopic retrieves a topic by name
func (tr *TopicRepository) GetTopic(name string) *Topic {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()

	return tr.topics[name]
}

// TopicService interface defines operations for topic management
type TopicService interface {
	CreateTopic(name string, numPartitions int) *Topic
	PublishMessage(topic *Topic, message *Message) *MessageMetadata
}

// TopicServiceImpl implements the TopicService interface
type TopicServiceImpl struct {
	topicRepository *TopicRepository
	routingStrategy RoutingStrategy
}

// NewTopicServiceImpl creates a new topic service implementation
func NewTopicServiceImpl() *TopicServiceImpl {
	return &TopicServiceImpl{
		topicRepository: NewTopicRepository(),
		routingStrategy: NewRoundRobinRoutingStrategy(),
	}
}

// CreateTopic creates a new topic with the specified number of partitions
func (ts *TopicServiceImpl) CreateTopic(name string, numPartitions int) *Topic {
	topic := NewTopic(name, numPartitions)
	return ts.topicRepository.AddTopic(topic)
}

// PublishMessage publishes a message to a topic using the routing strategy
func (ts *TopicServiceImpl) PublishMessage(topic *Topic, message *Message) *MessageMetadata {
	if topic == nil || message == nil {
		return nil
	}

	return ts.routingStrategy.Route(topic, message)
}

// RoutingStrategy interface defines how messages are routed to partitions
type RoutingStrategy interface {
	Route(topic *Topic, message *Message) *MessageMetadata
}

// RoundRobinRoutingStrategy implements round-robin routing
type RoundRobinRoutingStrategy struct {
	counter int
	mutex   sync.Mutex
}

// NewRoundRobinRoutingStrategy creates a new round-robin routing strategy
func NewRoundRobinRoutingStrategy() *RoundRobinRoutingStrategy {
	return &RoundRobinRoutingStrategy{}
}

// Route routes a message to a partition using round-robin strategy
func (rr *RoundRobinRoutingStrategy) Route(topic *Topic, message *Message) *MessageMetadata {
	rr.mutex.Lock()
	defer rr.mutex.Unlock()

	if len(topic.Partitions) == 0 {
		return nil
	}

	// Use random selection for simplicity (as in original Java code)
	partitionIndex := rand.Intn(len(topic.Partitions))
	partition := topic.Partitions[partitionIndex]

	offset := partition.AddMessage(message)
	return NewMessageMetadata(partition, offset)
}

// ConsumerService interface defines operations for consumer management
type ConsumerService interface {
	Subscribe(topic *Topic, consumer *Consumer) *Partition
	Consume(consumer *Consumer) *Message
}

// ConsumerServiceImpl implements the ConsumerService interface
type ConsumerServiceImpl struct {
	consumerGroupService ConsumerGroupService
}

// NewConsumerServiceImpl creates a new consumer service implementation
func NewConsumerServiceImpl() *ConsumerServiceImpl {
	return &ConsumerServiceImpl{
		consumerGroupService: NewConsumerGroupServiceImpl(),
	}
}

// Subscribe subscribes a consumer to a topic and returns the allocated partition
func (cs *ConsumerServiceImpl) Subscribe(topic *Topic, consumer *Consumer) *Partition {
	if topic == nil || consumer == nil {
		return nil
	}

	return cs.consumerGroupService.RegisterConsumer(topic, consumer)
}

// Consume consumes a message for the given consumer
func (cs *ConsumerServiceImpl) Consume(consumer *Consumer) *Message {
	return cs.consumerGroupService.GetAMessage(consumer)
}

// ConsumerGroupService interface defines operations for consumer group management
type ConsumerGroupService interface {
	RegisterConsumer(topic *Topic, consumer *Consumer) *Partition
	GetAMessage(consumer *Consumer) *Message
}

// ConsumerGroupServiceImpl implements the ConsumerGroupService interface
type ConsumerGroupServiceImpl struct{}

// NewConsumerGroupServiceImpl creates a new consumer group service implementation
func NewConsumerGroupServiceImpl() *ConsumerGroupServiceImpl {
	return &ConsumerGroupServiceImpl{}
}

// RegisterConsumer registers a consumer with a topic and allocates a partition
func (cgs *ConsumerGroupServiceImpl) RegisterConsumer(topic *Topic, consumer *Consumer) *Partition {
	var allocatedPartition *Partition
	consumerGroup := cgs.fetchConsumerGroup(consumer)

	if consumerGroup.GetConsumerCount() > topic.MaxPartitions {
		// No need to allocate partition to the new consumer as there are already consumers consuming them.
		// This consumer will remain idle
		return allocatedPartition
	}

	for _, topicPartition := range topic.Partitions {
		if !consumerGroup.IsPartitionAllocated(topicPartition) {
			allocatedPartition = consumerGroup.AllocatePartition(topicPartition, consumer)
			break
		}
	}

	return allocatedPartition
}

// fetchConsumerGroup gets the consumer group for a consumer
func (cgs *ConsumerGroupServiceImpl) fetchConsumerGroup(consumer *Consumer) *ConsumerGroup {
	consumerGroup := consumer.ConsumerGroup
	if consumerGroup == nil {
		consumerGroup = GetDefaultConsumerGroup()
	}
	return consumerGroup
}

// GetAMessage retrieves a message for the given consumer
func (cgs *ConsumerGroupServiceImpl) GetAMessage(consumer *Consumer) *Message {
	consumerGroup := cgs.fetchConsumerGroup(consumer)
	return consumerGroup.GetMessageForConsumer(consumer)
}
