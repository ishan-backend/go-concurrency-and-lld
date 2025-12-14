Designing flexible data processing solutions using AI tooling while interpreting ambiguous requirements.

# Design Patterns in the Data Processing Framework

## Overview
This framework uses 6 key design patterns to achieve flexibility and extensibility.


https://www.concurrency.rocks/tutorials/pipeline
---

## 1. Strategy Pattern
**Purpose**: Define a family of algorithms (processors) that are interchangeable

**Implementation**:
```go
type Processor interface {
    Process(data []map[string]interface{}) ([]map[string]interface{}, error)
    Name() string
}
```

**Why it's flexible**:
- Easy to add new processors without modifying existing code
- Processors can be swapped at runtime
- Each processor has a single responsibility

**How to extend**:
```go
// Add a new processor
type DeduplicateProcessor struct {
    field string
}

func (dp *DeduplicateProcessor) Process(data []map[string]interface{}) ([]map[string]interface{}, error) {
    seen := make(map[string]bool)
    result := make([]map[string]interface{}, 0)
    
    for _, row := range data {
        key := fmt.Sprintf("%v", row[dp.field])
        if !seen[key] {
            seen[key] = true
            result = append(result, row)
        }
    }
    return result, nil
}

func (dp *DeduplicateProcessor) Name() string {
    return fmt.Sprintf("Deduplicate[%s]", dp.field)
}
```

---

## 2. Chain of Responsibility Pattern
**Purpose**: Pass data through a chain of processors

**Implementation**:
```go
type Pipeline struct {
    processors []Processor
}

func (p *Pipeline) Execute(data []map[string]interface{}) ([]map[string]interface{}, error) {
    result := data
    for _, processor := range p.processors {
        result, err = processor.Process(result)
        if err != nil {
            return nil, err
        }
    }
    return result, nil
}
```

**Why it's flexible**:
- Dynamic processor chains
- Easy to add/remove/reorder steps
- Each processor can modify or pass data along

---

## 3. Factory Pattern
**Purpose**: Create processors dynamically based on configuration

**Implementation**:
```go
type ProcessorFactory struct {
    creators map[string]func(config map[string]interface{}) (Processor, error)
}

func (f *ProcessorFactory) Create(name string, config map[string]interface{}) (Processor, error) {
    creator, exists := f.creators[name]
    if !exists {
        return nil, fmt.Errorf("unknown processor type: %s", name)
    }
    return creator(config)
}
```

**Why it's flexible**:
- Create processors from configuration (JSON, YAML, etc.)
- Register new processor types at runtime
- Centralized processor creation logic

**How to use with config files**:
```go
// Load from JSON config
type Config struct {
    Processors []struct {
        Type   string                 `json:"type"`
        Config map[string]interface{} `json:"config"`
    } `json:"processors"`
}

func BuildPipelineFromConfig(configPath string) (*Pipeline, error) {
    // Parse config, create processors using factory
    factory := NewProcessorFactory()
    pipeline := NewPipeline()
    
    for _, proc := range config.Processors {
        processor, err := factory.Create(proc.Type, proc.Config)
        if err != nil {
            return nil, err
        }
        pipeline.AddProcessor(processor)
    }
    return pipeline, nil
}
```

---

## 4. Adapter Pattern
**Purpose**: Support multiple data sources and sinks with a common interface

**Implementation**:
```go
type DataSource interface {
    Read() ([]map[string]interface{}, error)
    Close() error
}

type DataSink interface {
    Write(data []map[string]interface{}) error
    Close() error
}
```

**Why it's flexible**:
- Add new sources (database, API, Kafka, etc.) without changing pipeline code
- Swap input/output formats easily
- Consistent interface for all data I/O

**How to extend with new sources**:
```go
// Database source
type PostgresDataSource struct {
    db    *sql.DB
    query string
}

func (ds *PostgresDataSource) Read() ([]map[string]interface{}, error) {
    rows, err := ds.db.Query(ds.query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    // Convert rows to []map[string]interface{}
    // ... implementation
}

// API source
type APIDataSource struct {
    url    string
    client *http.Client
}

func (ds *APIDataSource) Read() ([]map[string]interface{}, error) {
    resp, err := ds.client.Get(ds.url)
    // ... parse response
}
```

---

## 5. Builder Pattern
**Purpose**: Fluent API for constructing pipelines

**Implementation**:
```go
type PipelineBuilder struct {
    pipeline *Pipeline
    factory  *ProcessorFactory
}

func (pb *PipelineBuilder) Filter(field, operation string, value interface{}) *PipelineBuilder {
    // Create and add processor
    return pb // Return self for chaining
}
```

**Why it's flexible**:
- Readable, chainable API
- Hide complexity of processor creation
- Easy to construct complex pipelines

**Usage**:
```go
pipeline := NewPipelineBuilder().
    Filter("age", "greater_than", 18).
    Transform("email", "lowercase").
    Deduplicate("user_id").
    Aggregate("country").
    Build()
```

---

## 6. Template Method Pattern (Implicit)
**Purpose**: Define the skeleton of the processing algorithm

**Implementation**:
```go
func (dp *DataProcessor) Execute() error {
    // 1. Extract
    data, err := dp.source.Read()
    
    // 2. Transform
    result, err := dp.pipeline.Execute(data)
    
    // 3. Load
    err = dp.sink.Write(result)
    
    return nil
}
```

**Why it's flexible**:
- Fixed workflow, flexible implementation
- Classic ETL pattern
- Each step can be customized independently

---

## Common Extension Scenarios

### Adding a New Processor Type

1. **Implement the Processor interface**
2. **Create a factory function**
3. **Register with the factory**
4. **Add builder method (optional)**

```go
// 1. Implement interface
type SortProcessor struct {
    field string
    order string
}

func (sp *SortProcessor) Process(data []map[string]interface{}) ([]map[string]interface{}, error) {
    // Sort implementation
    return sortedData, nil
}

func (sp *SortProcessor) Name() string {
    return fmt.Sprintf("Sort[%s %s]", sp.field, sp.order)
}

// 2. Create factory function
func CreateSortProcessor(config map[string]interface{}) (Processor, error) {
    field := config["field"].(string)
    order := config["order"].(string)
    return &SortProcessor{field: field, order: order}, nil
}

// 3. Register in factory
factory.Register("sort", CreateSortProcessor)

// 4. Add builder method
func (pb *PipelineBuilder) Sort(field, order string) *PipelineBuilder {
    config := map[string]interface{}{
        "field": field,
        "order": order,
    }
    processor, _ := pb.factory.Create("sort", config)
    pb.pipeline.AddProcessor(processor)
    return pb
}
```

### Adding a New Data Source

```go
type KafkaDataSource struct {
    topic    string
    consumer *kafka.Consumer
}

func NewKafkaDataSource(topic string) *KafkaDataSource {
    // Initialize Kafka consumer
    return &KafkaDataSource{topic: topic}
}

func (ds *KafkaDataSource) Read() ([]map[string]interface{}, error) {
    // Read from Kafka, convert to standard format
}

func (ds *KafkaDataSource) Close() error {
    return ds.consumer.Close()
}
```

### Configuration-Driven Pipeline

```json
{
  "source": {
    "type": "csv",
    "path": "input.csv"
  },
  "processors": [
    {
      "type": "filter",
      "config": {
        "field": "status",
        "operation": "equals",
        "value": "active"
      }
    },
    {
      "type": "transform",
      "config": {
        "field": "name",
        "operation": "uppercase"
      }
    }
  ],
  "sink": {
    "type": "json",
    "path": "output.json"
  }
}
```

---

## Key Principles for Flexibility

1. **Program to interfaces, not implementations**
2. **Open/Closed Principle**: Open for extension, closed for modification
3. **Single Responsibility**: Each processor does one thing
4. **Dependency Inversion**: Depend on abstractions (interfaces)
5. **Composition over inheritance**: Build complex behavior by combining simple pieces

---

## Testing Strategy

```go
func TestFilterProcessor(t *testing.T) {
    processor := &FilterProcessor{
        field:     "status",
        operation: "equals",
        value:     "active",
    }
    
    input := []map[string]interface{}{
        {"id": 1, "status": "active"},
        {"id": 2, "status": "inactive"},
    }
    
    result, err := processor.Process(input)
    
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    if len(result) != 1 {
        t.Errorf("expected 1 result, got %d", len(result))
    }
}
```

---

## Interview Tips

1. **Start with interfaces**: Define contracts first
2. **Think about extension points**: Where will requirements change?
3. **Use composition**: Build complex from simple
4. **Configuration over code**: Make behavior configurable
5. **Handle errors gracefully**: Every step can fail
6. **Make it testable**: Inject dependencies, use interfaces

## Common AI Prompts for Extensions

- "Add a processor that validates data against a JSON schema"
- "Implement a rate limiter for API data sources"
- "Add support for streaming large files in chunks"
- "Create a processor that enriches data from an external API"
- "Add retry logic with exponential backoff for failed operations"

-------------------------------------------------------------------------------------------

# SOLID Principles in Go

A comprehensive guide to applying SOLID principles in Go, with practical examples and idiomatic patterns.

## Table of Contents
1. Single Responsibility Principle (SRP)
2. Open/Closed Principle (OCP)
3. Liskov Substitution Principle (LSP)
4. Interface Segregation Principle (ISP)
5. Dependency Inversion Principle (DIP)

---

## 1. Single Responsibility Principle (SRP)

**"A class should have one, and only one, reason to change."**

In Go, this applies to packages, types, and functions. Each should have a single, well-defined purpose.

### ❌ Bad Example

```go
type User struct {
    Name  string
    Email string
}

// This type has too many responsibilities
func (u *User) Save() error {
    // Database logic
    db := connectDB()
    return db.Save(u)
}

func (u *User) SendEmail(message string) error {
    // Email sending logic
    smtp := connectSMTP()
    return smtp.Send(u.Email, message)
}

func (u *User) ValidateEmail() bool {
    // Validation logic
    return strings.Contains(u.Email, "@")
}
```

### ✅ Good Example

```go
// User struct focuses only on data
type User struct {
    Name  string
    Email string
}

// Separate repository for persistence
type UserRepository struct {
    db *sql.DB
}

func (r *UserRepository) Save(user *User) error {
    // Database logic only
    return r.db.Exec("INSERT INTO users...")
}

// Separate service for email
type EmailService struct {
    smtp *smtp.Client
}

func (s *EmailService) Send(to, message string) error {
    // Email logic only
    return s.smtp.SendMail(to, message)
}

// Separate validator
type UserValidator struct{}

func (v *UserValidator) ValidateEmail(email string) bool {
    // Validation logic only
    return strings.Contains(email, "@")
}
```

---

## 2. Open/Closed Principle (OCP)

**"Software entities should be open for extension but closed for modification."**

Use interfaces and composition to extend behavior without modifying existing code.

### ❌ Bad Example

```go
type PaymentProcessor struct{}

func (p *PaymentProcessor) ProcessPayment(method string, amount float64) error {
    if method == "credit_card" {
        // Process credit card
        return processCreditCard(amount)
    } else if method == "paypal" {
        // Process PayPal
        return processPayPal(amount)
    } else if method == "crypto" {
        // Process crypto - requires modifying this function
        return processCrypto(amount)
    }
    return errors.New("unknown payment method")
}
```

### ✅ Good Example

```go
// Define an interface for payment methods
type PaymentMethod interface {
    Process(amount float64) error
}

// Concrete implementations
type CreditCardPayment struct {
    CardNumber string
}

func (c *CreditCardPayment) Process(amount float64) error {
    // Credit card processing logic
    return nil
}

type PayPalPayment struct {
    Email string
}

func (p *PayPalPayment) Process(amount float64) error {
    // PayPal processing logic
    return nil
}

// New payment method - no modification needed!
type CryptoPayment struct {
    WalletAddress string
}

func (c *CryptoPayment) Process(amount float64) error {
    // Crypto processing logic
    return nil
}

// Processor works with any PaymentMethod
type PaymentProcessor struct{}

func (p *PaymentProcessor) ProcessPayment(method PaymentMethod, amount float64) error {
    return method.Process(amount)
}
```

---

## 3. Liskov Substitution Principle (LSP)

**"Objects of a superclass should be replaceable with objects of a subclass without breaking the application."**

In Go, this means any implementation of an interface should be interchangeable.

### ❌ Bad Example

```go
type Bird interface {
    Fly() error
}

type Sparrow struct{}

func (s *Sparrow) Fly() error {
    fmt.Println("Sparrow flying")
    return nil
}

type Penguin struct{}

func (p *Penguin) Fly() error {
    // Violates LSP - penguins can't fly!
    return errors.New("penguins cannot fly")
}

// This function expects all birds to fly
func MakeBirdFly(b Bird) error {
    return b.Fly() // Breaks for Penguin!
}
```

### ✅ Good Example

```go
// Separate capabilities into different interfaces
type Swimmer interface {
    Swim() error
}

type Flyer interface {
    Fly() error
}

type Sparrow struct{}

func (s *Sparrow) Fly() error {
    fmt.Println("Sparrow flying")
    return nil
}

type Penguin struct{}

func (p *Penguin) Swim() error {
    fmt.Println("Penguin swimming")
    return nil
}

// Functions work with appropriate interfaces
func MakeFly(f Flyer) error {
    return f.Fly()
}

func MakeSwim(s Swimmer) error {
    return s.Swim()
}
```

---

## 4. Interface Segregation Principle (ISP)

**"No client should be forced to depend on methods it does not use."**

Create small, focused interfaces rather than large, monolithic ones.

### ❌ Bad Example

```go
// Fat interface forces implementations to provide all methods
type Worker interface {
    Work() error
    Eat() error
    Sleep() error
    TakePTO() error
}

// Robot has to implement methods it doesn't need
type Robot struct{}

func (r *Robot) Work() error {
    return nil
}

func (r *Robot) Eat() error {
    // Robots don't eat!
    return nil
}

func (r *Robot) Sleep() error {
    // Robots don't sleep!
    return nil
}

func (r *Robot) TakePTO() error {
    // Robots don't take time off!
    return nil
}
```

### ✅ Good Example

```go
// Small, focused interfaces
type Worker interface {
    Work() error
}

type Eater interface {
    Eat() error
}

type Sleeper interface {
    Sleep() error
}

type PTOTaker interface {
    TakePTO() error
}

// Robot only implements what it needs
type Robot struct{}

func (r *Robot) Work() error {
    return nil
}

// Human implements multiple interfaces
type Human struct{}

func (h *Human) Work() error {
    return nil
}

func (h *Human) Eat() error {
    return nil
}

func (h *Human) Sleep() error {
    return nil
}

func (h *Human) TakePTO() error {
    return nil
}

// Functions depend on specific interfaces
func DoWork(w Worker) error {
    return w.Work()
}

func FeedCreature(e Eater) error {
    return e.Eat()
}
```

---

## 5. Dependency Inversion Principle (DIP)

**"High-level modules should not depend on low-level modules. Both should depend on abstractions."**

Depend on interfaces, not concrete implementations. Use dependency injection.

### ❌ Bad Example

```go
// Tight coupling to concrete MySQL implementation
type UserService struct {
    db *MySQLDatabase // Direct dependency on concrete type
}

func NewUserService() *UserService {
    return &UserService{
        db: NewMySQLDatabase(), // Hard-coded dependency
    }
}

func (s *UserService) GetUser(id int) (*User, error) {
    return s.db.Query(id) // Tied to MySQL
}

type MySQLDatabase struct{}

func NewMySQLDatabase() *MySQLDatabase {
    return &MySQLDatabase{}
}

func (db *MySQLDatabase) Query(id int) (*User, error) {
    // MySQL-specific logic
    return nil, nil
}
```

### ✅ Good Example

```go
// Define abstraction
type UserRepository interface {
    GetByID(id int) (*User, error)
    Save(user *User) error
}

// High-level service depends on abstraction
type UserService struct {
    repo UserRepository // Depends on interface
}

// Dependency injection via constructor
func NewUserService(repo UserRepository) *UserService {
    return &UserService{
        repo: repo,
    }
}

func (s *UserService) GetUser(id int) (*User, error) {
    return s.repo.GetByID(id)
}

// Concrete implementations
type MySQLRepository struct {
    db *sql.DB
}

func (r *MySQLRepository) GetByID(id int) (*User, error) {
    // MySQL implementation
    return nil, nil
}

func (r *MySQLRepository) Save(user *User) error {
    return nil
}

type PostgresRepository struct {
    db *sql.DB
}

func (r *PostgresRepository) GetByID(id int) (*User, error) {
    // Postgres implementation
    return nil, nil
}

func (r *PostgresRepository) Save(user *User) error {
    return nil
}

// Easy to swap implementations
func main() {
    mysqlRepo := &MySQLRepository{}
    service := NewUserService(mysqlRepo)
    
    // Or use Postgres
    postgresRepo := &PostgresRepository{}
    service = NewUserService(postgresRepo)
}
```

---

## Complete Example: E-Commerce Order System

Here's how all SOLID principles work together:

```go
package main

import (
    "errors"
    "fmt"
)

// Models
type Order struct {
    ID     string
    Total  float64
    Status string
}

// ISP: Small, focused interfaces
type OrderRepository interface {
    Save(order *Order) error
    FindByID(id string) (*Order, error)
}

type PaymentProcessor interface {
    Process(amount float64) error
}

type NotificationSender interface {
    Send(message string) error
}

// SRP: Each service has a single responsibility
type OrderService struct {
    repo         OrderRepository      // DIP: Depend on abstraction
    payment      PaymentProcessor     // DIP: Depend on abstraction
    notification NotificationSender   // DIP: Depend on abstraction
}

// DIP: Constructor injection
func NewOrderService(
    repo OrderRepository,
    payment PaymentProcessor,
    notification NotificationSender,
) *OrderService {
    return &OrderService{
        repo:         repo,
        payment:      payment,
        notification: notification,
    }
}

func (s *OrderService) PlaceOrder(order *Order) error {
    // Process payment
    if err := s.payment.Process(order.Total); err != nil {
        return fmt.Errorf("payment failed: %w", err)
    }
    
    // Save order
    order.Status = "confirmed"
    if err := s.repo.Save(order); err != nil {
        return fmt.Errorf("save failed: %w", err)
    }
    
    // Send notification
    msg := fmt.Sprintf("Order %s confirmed", order.ID)
    if err := s.notification.Send(msg); err != nil {
        // Log but don't fail
        fmt.Printf("notification failed: %v\n", err)
    }
    
    return nil
}

// OCP & LSP: Implementations can be extended without modification
type StripePayment struct {
    apiKey string
}

func (s *StripePayment) Process(amount float64) error {
    fmt.Printf("Processing $%.2f via Stripe\n", amount)
    return nil
}

type PayPalPayment struct {
    email string
}

func (p *PayPalPayment) Process(amount float64) error {
    fmt.Printf("Processing $%.2f via PayPal\n", amount)
    return nil
}

type InMemoryOrderRepository struct {
    orders map[string]*Order
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
    return &InMemoryOrderRepository{
        orders: make(map[string]*Order),
    }
}

func (r *InMemoryOrderRepository) Save(order *Order) error {
    r.orders[order.ID] = order
    return nil
}

func (r *InMemoryOrderRepository) FindByID(id string) (*Order, error) {
    if order, ok := r.orders[id]; ok {
        return order, nil
    }
    return nil, errors.New("order not found")
}

type EmailNotification struct {
    smtpServer string
}

func (e *EmailNotification) Send(message string) error {
    fmt.Printf("Sending email: %s\n", message)
    return nil
}

type SMSNotification struct {
    phoneNumber string
}

func (s *SMSNotification) Send(message string) error {
    fmt.Printf("Sending SMS: %s\n", message)
    return nil
}

func main() {
    // Wire up dependencies
    repo := NewInMemoryOrderRepository()
    payment := &StripePayment{apiKey: "sk_test_123"}
    notification := &EmailNotification{smtpServer: "smtp.example.com"}
    
    service := NewOrderService(repo, payment, notification)
    
    // Place an order
    order := &Order{
        ID:    "ORD-001",
        Total: 99.99,
    }
    
    if err := service.PlaceOrder(order); err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Println("Order placed successfully!")
    
    // Easy to swap implementations
    service2 := NewOrderService(
        repo,
        &PayPalPayment{email: "user@example.com"},
        &SMSNotification{phoneNumber: "+1234567890"},
    )
    
    order2 := &Order{
        ID:    "ORD-002",
        Total: 149.99,
    }
    
    service2.PlaceOrder(order2)
}
```

---

## Key Takeaways for Go

1. **Use interfaces liberally** - Go's implicit interface satisfaction makes SOLID principles natural
2. **Keep interfaces small** - Go encourages many small interfaces over large ones
3. **Favor composition over inheritance** - Go doesn't have inheritance; use embedding and composition
4. **Dependency injection** - Pass dependencies as parameters, typically through constructors
5. **Package organization** - Apply SRP at the package level; each package should have a clear purpose

## Common Go Patterns Aligned with SOLID

- **Repository Pattern** - DIP, SRP
- **Strategy Pattern** - OCP, ISP
- **Dependency Injection** - DIP
- **Interface Composition** - ISP
- **Table-Driven Tests** - Makes testing SOLID code easier

These principles help create maintainable, testable, and flexible Go applications!