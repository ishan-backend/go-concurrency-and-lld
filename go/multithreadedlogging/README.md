# Multithreaded Logging System - Go Implementation

This is a Go implementation of a multithreaded logging system, equivalent to the Java version in `src/com/lld/multithreadedlogging/`.

## Design Requirements

1. The system supports multiple log levels (error, debug, warning, info)
2. Multiple log mediums are supported (file-based logs, network logs, DB logs)
3. Multiple producers can produce to the queue simultaneously
4. Irrespective of the rate of production, the system should consume messages as fast as possible from the queue and flush them to the log medium
5. A single application message can get published to multiple mediums

## Go-Specific Implementation Details

### ArrayBlockingQueue Equivalent: Go Channels

In Java, the system uses `ArrayBlockingQueue<WritableMessage>` for thread-safe message queuing. In Go, this is implemented using **buffered channels**:

```go
// Java: ArrayBlockingQueue<WritableMessage> blockingQueue
// Go: chan *model.WritableMessage
queue: make(chan *model.WritableMessage, capacity)
```

#### Key Differences:

1. **Blocking Operations**:
   - Java: `blockingQueue.offer(message, timeout, timeUnit)` and `blockingQueue.poll(timeout, timeUnit)`
   - Go: `select` statement with `case q.queue <- message:` and `case message := <-q.queue:`

2. **Timeout Handling**:
   - Java: Built-in timeout support in `offer()` and `poll()`
   - Go: Implemented using `time.After()` in `select` statements

3. **Thread Safety**:
   - Java: `ArrayBlockingQueue` is inherently thread-safe
   - Go: Channels are inherently thread-safe and goroutine-safe

### Concurrency Model

#### Java vs Go Approach:

**Java**:
```java
ExecutorService executorService = Executors.newFixedThreadPool(3);
executorService.execute(new Producer(queueService));
executorService.execute(new Producer(queueService));
executorService.execute(new Consumer(queueService));
```

**Go**:
```go
var wg sync.WaitGroup
wg.Add(2)
go func() {
    defer wg.Done()
    producer := worker.NewProducer(queueService)
    producer.Run(ctx)
}()
```

#### Key Go Concurrency Features Used:

1. **Goroutines**: Lightweight threads for producers and consumers
2. **Channels**: For thread-safe communication between goroutines
3. **Context**: For graceful shutdown and cancellation
4. **WaitGroup**: For synchronizing goroutine completion
5. **Mutex/Sync.Once**: For singleton pattern implementation

### Singleton Pattern Implementation

**Java**:
```java
private static DBInternalLogger instance;
public static DBInternalLogger getInstance() {
    if (instance == null) {
        synchronized (DBInternalLogger.class) {
            if (instance == null) {
                instance = new DBInternalLogger();
            }
        }
    }
    return instance;
}
```

**Go**:
```go
var (
    dbLoggerInstance *dbInternalLogger
    dbLoggerOnce     sync.Once
)

func GetDBInternalLogger() logger.InternalLogger {
    dbLoggerOnce.Do(func() {
        dbLoggerInstance = &dbInternalLogger{}
    })
    return dbLoggerInstance
}
```

### Error Handling

**Java**: Uses `InterruptedException` for thread interruption
**Go**: Uses `context.Context` for cancellation and `error` interface for error handling

### Message Processing Flow

1. **Producer**: Creates `ApplicationAwareMessage` → Converts to `WritableMessage` for each medium → Publishes to queue
2. **Queue Service**: Buffered channel acts as message buffer
3. **Consumer**: Consumes messages → Gets appropriate logger via factory → Logs message

## Running the Application

```bash
cd go/multithreadedlogging
go mod tidy
go run main.go
```

## Architecture Components

### Core Components

1. **Models**:
   - `ApplicationAwareMessage`: Application-level message
   - `WritableMessage`: Internal message with metadata

2. **Services**:
   - `QueueService`: Interface for message queuing
   - `queueServiceImpl`: Channel-based implementation

3. **Loggers**:
   - `ApplicationAwareLogger`: Application-facing interface
   - `InternalLogger`: Internal logging interface
   - Implementations: `DBInternalLogger`, `FileInternalLogger`, `NetworkInternalLogger`

4. **Workers**:
   - `Producer`: Message producer goroutine
   - `Consumer`: Message consumer goroutine

5. **Factory**:
   - `InternalLoggerFactory`: Creates logger instances based on medium

### Configuration

- `Config`: Holds timeout and capacity settings
- Default capacity: 100 messages
- Default timeouts: 100ms for publish/consume operations

## Go Advantages Over Java Implementation

1. **Lower Memory Overhead**: Goroutines are much lighter than Java threads
2. **Better Concurrency**: Go's channel-based communication is more efficient than Java's blocking queues
3. **Simpler Error Handling**: Context-based cancellation is more elegant than exception handling
4. **Compile-time Safety**: Go's type system catches many issues at compile time
5. **Built-in Concurrency Primitives**: Channels, goroutines, and sync primitives are first-class citizens

## Performance Considerations

1. **Channel Buffer Size**: Configurable based on expected message volume
2. **Goroutine Pool**: Could be implemented for high-throughput scenarios
3. **Batch Processing**: Could be added for better throughput
4. **Connection Pooling**: For database and network loggers in production

## Production Enhancements

1. **Structured Logging**: Use structured logging libraries like `logrus` or `zap`
2. **Metrics**: Add Prometheus metrics for monitoring
3. **Configuration**: Use `viper` for configuration management
4. **Graceful Shutdown**: Implement proper cleanup on shutdown
5. **Circuit Breakers**: For external service logging
6. **Retry Logic**: For failed log operations

----------------------------------

Reviewing the Java multithreaded logging implementation to identify core components and patterns.
## Summary

Implemented a Go version of the multithreaded logging system from `src/com/lld/multithreadedlogging/`.

### Go implementation structure

```
go/multithreadedlogging/
├── main.go                           # Main application entry point
├── go.mod                           # Go module definition
├── README.md                        # Comprehensive documentation
├── constant/
│   ├── log_level.go                 # Log level enum
│   └── supported_log_medium.go      # Log medium enum
├── configuration/
│   └── config.go                    # Configuration constants
├── model/
│   ├── application_aware_message.go # Application message model
│   └── writable_message.go         # Internal message model
├── service/
│   └── queue_service.go             # Queue service implementation
├── logger/
│   ├── application_aware_logger.go  # Application logger interface
│   ├── internal_logger.go           # Internal logger interface
│   ├── factory/
│   │   └── internal_logger_factory.go # Logger factory
│   └── impl/
│       ├── application_aware_logger_impl.go # Application logger implementation
│       ├── db_internal_logger.go    # Database logger
│       ├── file_internal_logger.go  # File logger
│       └── network_internal_logger.go # Network logger
└── worker/
    ├── producer.go                   # Message producer worker
    └── consumer.go                   # Message consumer worker
```

### ArrayBlockingQueue equivalent: Go channels

Implemented with buffered channels:

```go
// Java: ArrayBlockingQueue<WritableMessage> blockingQueue
// Go: chan *model.WritableMessage
queue: make(chan *model.WritableMessage, capacity)
```

#### Differences

1. Blocking operations:
   - Java: `blockingQueue.offer(message, timeout, timeUnit)` and `blockingQueue.poll(timeout, timeUnit)`
   - Go: `select` with `case q.queue <- message:` and `case message := <-q.queue:`

2. Timeout handling:
   - Java: Built-in timeout support
   - Go: Implemented using `time.After()` in `select` statements

3. Thread safety:
   - Java: `ArrayBlockingQueue` is thread-safe
   - Go: Channels are goroutine-safe

### Concurrency model

#### Java vs Go

**Java**:
```java
ExecutorService executorService = Executors.newFixedThreadPool(3);
executorService.execute(new Producer(queueService));
executorService.execute(new Producer(queueService));
executorService.execute(new Consumer(queueService));
```

**Go**:
```go
var wg sync.WaitGroup
wg.Add(2)
go func() {
    defer wg.Done()
    producer := worker.NewProducer(queueService)
    producer.Run(ctx)
}()
```

### Go concurrency features used

1. Goroutines
2. Channels
3. Context for cancellation
4. WaitGroup for synchronization
5. sync.Once for singletons

### Singleton pattern

**Java**:
```java
private static DBInternalLogger instance;
public static DBInternalLogger getInstance() {
    if (instance == null) {
        synchronized (DBInternalLogger.class) {
            if (instance == null) {
                instance = new DBInternalLogger();
            }
        }
    }
    return instance;
}
```

**Go**:
```go
var (
    dbLoggerInstance *dbInternalLogger
    dbLoggerOnce     sync.Once
)

func GetDBInternalLogger() logger.InternalLogger {
    dbLoggerOnce.Do(func() {
        dbLoggerInstance = &dbInternalLogger{}
    })
    return dbLoggerInstance
}
```

### How to run

```bash
cd go/multithreadedlogging
go mod tidy
go run main.go
```

### Advantages

1. Lower memory overhead
2. More efficient concurrency with channels
3. Context-based cancellation
4. Compile-time safety
5. Built-in concurrency primitives

The implementation preserves the original design while using Go idioms. Channels replace `ArrayBlockingQueue` for thread-safe queuing, and goroutines replace Java threads.