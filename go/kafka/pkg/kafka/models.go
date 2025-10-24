package kafka

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// Message represents a message in the Kafka system
type Message struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

// NewMessage creates a new message with the given content
func NewMessage(content string) *Message {
	return &Message{
		ID:      uuid.New().String(),
		Content: content,
	}
}

// String returns the string representation of the message
func (m *Message) String() string {
	return m.Content
}

// Topic represents a Kafka topic with multiple partitions
type Topic struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	Partitions    []*Partition `json:"partitions"`
	MaxPartitions int          `json:"max_partitions"`
}

// NewTopic creates a new topic with the specified name and number of partitions
func NewTopic(name string, numPartitions int) *Topic {
	topic := &Topic{
		ID:            uuid.New().String(),
		Name:          name,
		MaxPartitions: numPartitions,
		Partitions:    make([]*Partition, 0, numPartitions),
	}

	// Generate partitions
	for i := 0; i < numPartitions; i++ {
		topic.Partitions = append(topic.Partitions, NewPartition(topic))
	}

	return topic
}

// String returns the string representation of the topic
func (t *Topic) String() string {
	return t.Name
}

// Partition represents a partition within a topic
type Partition struct {
	ID       string     `json:"id"`
	Topic    *Topic     `json:"topic"`
	Messages []*Message `json:"messages"`
	mutex    sync.RWMutex
}

// NewPartition creates a new partition for the given topic
func NewPartition(topic *Topic) *Partition {
	return &Partition{
		ID:       uuid.New().String(),
		Topic:    topic,
		Messages: make([]*Message, 0),
	}
}

// AddMessage adds a message to the partition and returns the offset
func (p *Partition) AddMessage(message *Message) int {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.Messages = append(p.Messages, message)
	return len(p.Messages) - 1
}

// GetMessageAtOffset retrieves a message at the specified offset
func (p *Partition) GetMessageAtOffset(offset int) *Message {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if offset >= len(p.Messages) {
		return nil
	}

	return p.Messages[offset]
}

// Producer represents a message producer
type Producer struct {
	ID string `json:"id"`
}

// NewProducer creates a new producer
func NewProducer() *Producer {
	return &Producer{
		ID: uuid.New().String(),
	}
}

// String returns the string representation of the producer
func (p *Producer) String() string {
	return fmt.Sprintf("Producer-Id = %s", p.ID)
}

// Consumer represents a message consumer
type Consumer struct {
	ID            string         `json:"id"`
	ConsumerGroup *ConsumerGroup `json:"consumer_group"`
}

// NewConsumer creates a new consumer and adds it to the consumer group
func NewConsumer(consumerGroup *ConsumerGroup) *Consumer {
	consumer := &Consumer{
		ID:            uuid.New().String(),
		ConsumerGroup: consumerGroup,
	}

	if consumerGroup != nil {
		consumerGroup.AddConsumer(consumer)
	}

	return consumer
}

// String returns the string representation of the consumer
func (c *Consumer) String() string {
	return c.ID
}

// ConsumerGroup represents a group of consumers
type ConsumerGroup struct {
	GroupID              string                `json:"group_id"`
	Consumers            map[string]*Consumer  `json:"consumers"`
	CurrentOffsets       map[string]int        `json:"current_offsets"`        // partition ID -> offset
	CommittedOffsets     map[string]int        `json:"committed_offsets"`      // partition ID -> committed offset
	ConsumerPartitionMap map[string]*Partition `json:"consumer_partition_map"` // consumer ID -> partition
	mutex                sync.RWMutex
}

var (
	defaultConsumerGroup *ConsumerGroup
	defaultGroupOnce     sync.Once
)

// GetDefaultConsumerGroup returns the default consumer group singleton
func GetDefaultConsumerGroup() *ConsumerGroup {
	defaultGroupOnce.Do(func() {
		defaultConsumerGroup = &ConsumerGroup{
			GroupID:              "default-consumer-group",
			Consumers:            make(map[string]*Consumer),
			CurrentOffsets:       make(map[string]int),
			CommittedOffsets:     make(map[string]int),
			ConsumerPartitionMap: make(map[string]*Partition),
		}
	})
	return defaultConsumerGroup
}

// NewConsumerGroup creates a new consumer group
func NewConsumerGroup() *ConsumerGroup {
	return &ConsumerGroup{
		GroupID:              uuid.New().String(),
		Consumers:            make(map[string]*Consumer),
		CurrentOffsets:       make(map[string]int),
		CommittedOffsets:     make(map[string]int),
		ConsumerPartitionMap: make(map[string]*Partition),
	}
}

// AddConsumer adds a consumer to the group
func (cg *ConsumerGroup) AddConsumer(consumer *Consumer) {
	if consumer == nil {
		return
	}

	cg.mutex.Lock()
	defer cg.mutex.Unlock()

	cg.Consumers[consumer.ID] = consumer
}

// GetConsumerCount returns the number of consumers in the group
func (cg *ConsumerGroup) GetConsumerCount() int {
	cg.mutex.RLock()
	defer cg.mutex.RUnlock()

	return len(cg.Consumers)
}

// IsPartitionAllocated checks if a partition is already allocated to any consumer
func (cg *ConsumerGroup) IsPartitionAllocated(partition *Partition) bool {
	cg.mutex.RLock()
	defer cg.mutex.RUnlock()

	for _, allocatedPartition := range cg.ConsumerPartitionMap {
		if allocatedPartition.ID == partition.ID {
			return true
		}
	}
	return false
}

// AllocatePartition allocates a partition to a consumer
func (cg *ConsumerGroup) AllocatePartition(partition *Partition, consumer *Consumer) *Partition {
	cg.mutex.Lock()
	defer cg.mutex.Unlock()

	cg.CurrentOffsets[partition.ID] = -1
	cg.ConsumerPartitionMap[consumer.ID] = partition
	return partition
}

// GetMessageForConsumer retrieves the next message for a consumer
func (cg *ConsumerGroup) GetMessageForConsumer(consumer *Consumer) *Message {
	cg.mutex.Lock()
	defer cg.mutex.Unlock()

	partition := cg.ConsumerPartitionMap[consumer.ID]
	if partition == nil {
		return nil
	}

	lastOffset := cg.CurrentOffsets[partition.ID]
	nextOffset := lastOffset + 1

	message := partition.GetMessageAtOffset(nextOffset)
	if message == nil {
		return nil
	}

	cg.CurrentOffsets[partition.ID] = nextOffset
	return message
}

// MessageMetadata represents metadata about a published message
type MessageMetadata struct {
	Partition *Partition `json:"partition"`
	Offset    int        `json:"offset"`
}

// NewMessageMetadata creates new message metadata
func NewMessageMetadata(partition *Partition, offset int) *MessageMetadata {
	return &MessageMetadata{
		Partition: partition,
		Offset:    offset,
	}
}

// String returns the string representation of message metadata
func (mm *MessageMetadata) String() string {
	return fmt.Sprintf("Partition = %s offset = %d", mm.Partition.ID, mm.Offset)
}
