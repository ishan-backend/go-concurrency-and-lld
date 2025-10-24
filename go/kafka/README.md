# Kafka-like Message Broker in Go

This is a Go implementation of a Kafka-like message broker, converted from the original Java implementation. It demonstrates various design patterns, concurrency concepts, and architectural principles.

## Architecture Overview

### Core Components

#### 1. Models (`pkg/kafka/models.go`)

**Key Structs:**
- `Message`: Represents a message with ID and content
- `Topic`: Contains multiple partitions and metadata
- `Partition`: Stores messages with thread-safe access
- `Producer`: Message producer entity
- `Consumer`: Message consumer entity
- `ConsumerGroup`: Manages multiple consumers and partition allocation
- `MessageMetadata`: Contains partition and offset information

**Concurrency Features:**
- `sync.RWMutex` for thread-safe read/write operations
- `sync.Once` for singleton pattern implementation
- Thread-safe partition message storage

#### 2. Services (`pkg/kafka/services.go`)

**Service Interfaces:**
- `TopicService`: Topic creation and message publishing
- `ConsumerService`: Consumer subscription and message consumption
- `ConsumerGroupService`: Consumer group management
- `RoutingStrategy`: Message routing to partitions

**Design Patterns:**
- **Repository Pattern**: `TopicRepository` for data persistence
- **Strategy Pattern**: `RoutingStrategy` for flexible message routing
- **Service Layer Pattern**: Separation of business logic

#### 3. Workers (`pkg/kafka/workers.go`)

**Worker Types:**
- `ProducerWorker`: Continuously produces messages
- `ConsumerWorker`: Continuously consumes messages
- `WorkerManager`: Manages multiple workers

**Concurrency Features:**
- Goroutines for concurrent execution
- Channels for communication and synchronization
- `time.Ticker` for periodic operations

## Design Patterns Used

### 1. Singleton Pattern
```go
var (
    defaultConsumerGroup *ConsumerGroup
    defaultGroupOnce     sync.Once
)

func GetDefaultConsumerGroup() *ConsumerGroup {
    defaultGroupOnce.Do(func() {
        defaultConsumerGroup = &ConsumerGroup{...}
    })
    return defaultConsumerGroup
}
```

### 2. Strategy Pattern
```go
type RoutingStrategy interface {
    Route(topic *Topic, message *Message) *MessageMetadata
}

type RoundRobinRoutingStrategy struct {
    counter int
    mutex   sync.Mutex
}
```

### 3. Repository Pattern
```go
type TopicRepository struct {
    topics map[string]*Topic
    mutex  sync.RWMutex
}
```

### 4. Service Layer Pattern
- Clear separation between business logic and data access
- Interface-based design for testability and flexibility

## Concurrency Concepts

### 1. Mutex Usage
- **RWMutex**: Used for read-heavy operations (e.g., topic repository)
- **Mutex**: Used for write operations and simple locking

### 2. Goroutines and Channels
```go
func (pw *ProducerWorker) Start() {
    go pw.run()
}

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
```

### 3. Select Statement
- Used for non-blocking operations
- Handles multiple channel operations
- Implements graceful shutdown

### 4. Worker Pattern
- Each producer/consumer runs in its own goroutine
- Centralized worker management
- Clean shutdown mechanism

## Key Differences from Java Implementation

### 1. Concurrency Model
- **Java**: Threads with `ExecutorService`
- **Go**: Goroutines with channels

### 2. Synchronization
- **Java**: `ReentrantLock`, `synchronized`
- **Go**: `sync.Mutex`, `sync.RWMutex`

### 3. Memory Management
- **Java**: Garbage collection with explicit memory management
- **Go**: Automatic memory management with garbage collection

### 4. Error Handling
- **Java**: Exceptions
- **Go**: Error values and explicit error checking

## Component Interactions

### 1. Message Flow
```
Producer -> TopicService -> RoutingStrategy -> Partition -> MessageMetadata
```

### 2. Consumer Flow
```
Consumer -> ConsumerService -> ConsumerGroupService -> ConsumerGroup -> Partition -> Message
```

### 3. Partition Allocation
```
Consumer -> ConsumerGroupService -> ConsumerGroup -> Partition Allocation
```

## Usage Example

```go
func main() {
    // Create services
    topicService := kafka.NewTopicServiceImpl()
    consumerService := kafka.NewConsumerServiceImpl()
    
    // Create topic
    topic := topicService.CreateTopic("my-topic", 3)
    
    // Create producer
    producer := kafka.NewProducer()
    
    // Create consumer group
    consumerGroup := kafka.NewConsumerGroup()
    
    // Create workers
    producerWorker := kafka.NewProducerWorker(topic, producer, topicService)
    consumerWorker := kafka.NewConsumerWorker(consumer, consumerService)
    
    // Start workers
    producerWorker.Start()
    consumerWorker.Start()
    
    // Let it run...
    time.Sleep(30 * time.Second)
    
    // Stop workers
    producerWorker.Stop()
    consumerWorker.Stop()
}
```

## Features Implemented

1. **Topic Management**: Create topics with multiple partitions
2. **Message Publishing**: Publish messages to topics using routing strategies
3. **Consumer Groups**: Multiple consumers in a group with partition allocation
4. **Concurrent Processing**: Producer and consumer workers run concurrently
5. **Thread Safety**: All operations are thread-safe using appropriate synchronization
6. **Graceful Shutdown**: Clean shutdown mechanism for workers

## Benefits of Go Implementation

1. **Lightweight Concurrency**: Goroutines are much lighter than threads
2. **Built-in Synchronization**: Channels provide elegant synchronization primitives
3. **Memory Efficiency**: Better memory usage compared to Java
4. **Simple Deployment**: Single binary deployment
5. **Fast Compilation**: Quick build times
6. **Excellent Concurrency Support**: Native support for concurrent programming

This implementation demonstrates how to build a scalable, concurrent message broker using Go's concurrency primitives and design patterns.

---------------------------------------------

Examining the Java Kafka package to map its structure and interactions.

## Go Kafka Package Structure

```
go/kafka/
├── pkg/kafka/
│   ├── models.go      # Core data structures
│   ├── services.go    # Business logic and interfaces
│   └── workers.go     # Concurrent workers
├── main.go           # Demo application
├── go.mod           # Go module definition
└── README.md        # Comprehensive documentation
```

## Design patterns and concepts

### 1. **Design patterns**
- **Singleton**: Default consumer group with `sync.Once`
- **Strategy**: Routing strategies for message distribution
- **Repository**: Topic persistence
- **Service layer**: Separation of business logic

### 2. **Concurrency**
- **Goroutines**: Replaces Java threads
- **Channels**: Communication and synchronization
- **Mutexes**: `sync.Mutex` and `sync.RWMutex` for thread safety
- **Select**: Non-blocking channel operations
- **Worker pattern**: Producer/consumer workers

### 3. **Component interactions**

**Message flow:**
```
Producer → TopicService → RoutingStrategy → Partition → MessageMetadata
```

**Consumer flow:**
```
Consumer → ConsumerService → ConsumerGroupService → ConsumerGroup → Partition → Message
```

**Partition allocation:**
```
Consumer → ConsumerGroupService → ConsumerGroup → Partition Allocation
```

## Differences from Java

| Aspect | Java | Go |
|--------|------|-----|
| Concurrency | Threads + ExecutorService | Goroutines + Channels |
| Synchronization | ReentrantLock, synchronized | sync.Mutex, sync.RWMutex |
| Memory | GC with explicit management | Automatic GC |
| Error handling | Exceptions | Error values |
| Deployment | JAR files | Single binary |

## Features
1. Thread-safe operations with mutexes
2. Concurrent producer/consumer workers
3. Partition allocation to consumers
4. Round-robin routing
5. Graceful shutdown via channels
6. Interface-based design

## Benefits
- Lightweight concurrency (goroutines)
- Built-in synchronization (channels)
- Efficient memory usage
- Single-binary deployment
- Fast compilation
- Native concurrency support

This Go version preserves the Java design while using Go’s concurrency primitives. The code is structured, documented, and ready to run.
