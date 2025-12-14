# Coupon Recommendation System - Low Level Design in Go

## Problem Statement

Design a Coupon Recommendation System for Zepto that fulfills the following requirements:

### Key Requirements

1. **Coupon Recommendation**: Recommend applicable coupons for a given order
2. **Dynamic Rule Evaluation**: Evaluate rules dynamically based on the user's order
3. **Flexible Conditions**: Support conditions like product type, user type, and payment method
4. **Concurrency**: Ensure thread safety for rule evaluation when multiple users access the system
5. **Optimization**: Optimize the rule evaluation process to minimize redundant checks

---

## Architecture Overview

This implementation follows a clean architecture with separation of concerns, using several design patterns to ensure maintainability, scalability, and thread safety.

### Design Patterns Used

1. **Strategy Pattern**: For different rule types (UserTypeRule, ProductTypeRule, etc.)
2. **Repository Pattern**: For data access abstraction
3. **Service Layer Pattern**: For business logic encapsulation
4. **Interface Segregation**: Clear interfaces for extensibility

---

## Core Components

### 1. Models (`models.go`)

The domain models represent the core entities:

- **User**: Represents a user with their type and purchase history
- **CartItem**: Individual items in the shopping cart
- **Cart**: Collection of cart items with total calculation
- **Order**: Complete order with user, cart, and payment information
- **Coupon**: Discount coupon with rules and formula

```go
type User struct {
    UserID        string
    UserType      UserType
    PastPurchases []string
}

type Coupon struct {
    CouponID       string
    DiscountAmount float64
    Rules          []Rule
    RulesFormula   string // e.g., "AND(0, 1)"
}
```

### 2. Rules System (`rules.go`)

The rules system uses the **Strategy Pattern** to allow different types of eligibility rules:

#### Rule Interface

```go
type Rule interface {
    Check(order *Order) bool
    GetRuleType() string
}
```

#### Implemented Rule Types

1. **UserTypeRule**: Checks if user type matches (e.g., Premium users only)
2. **ProductTypeRule**: Checks if cart contains specific product types
3. **MinimumAmountRule**: Validates minimum order amount
4. **PaymentMethodRule**: Validates payment method
5. **PastPurchaseRule**: Checks user's purchase history

#### Concurrent Rule Evaluation

The `RuleEvaluator` class handles concurrent rule evaluation:

```go
type RuleEvaluator struct {
    mu sync.RWMutex
}

func (re *RuleEvaluator) EvaluateRulesConcurrently(order *Order, rules []Rule) []bool {
    results := make([]bool, len(rules))
    var wg sync.WaitGroup
    
    for i, rule := range rules {
        wg.Add(1)
        go func(index int, r Rule) {
            defer wg.Done()
            results[index] = r.Check(order)
        }(i, rule)
    }
    
    wg.Wait()
    return results
}
```

**Key Features**:

- Uses `sync.RWMutex` for thread-safe read operations
- Evaluates multiple rules concurrently using goroutines
- Uses `sync.WaitGroup` to synchronize goroutine completion
- Returns results in the same order as input rules

### 3. Repository Pattern (`repository.go`)

The repository pattern abstracts data access and provides thread-safe operations:

```go
type CouponRepository interface {
    AddCoupon(coupon *Coupon) error
    GetCoupon(couponID string) (*Coupon, error)
    GetAllCoupons() []*Coupon
    RemoveCoupon(couponID string) error
}
```

**Thread Safety**:

- Uses `sync.RWMutex` for concurrent read/write operations
- Read operations use `RLock()` for better performance (multiple concurrent reads)
- Write operations use `Lock()` for exclusive access

### 4. Service Layer (`service.go`)

The service layer orchestrates business logic and handles concurrent coupon evaluation:

```go
type CouponRecommendationService struct {
    repository    CouponRepository
    ruleEvaluator *RuleEvaluator
    mu            sync.RWMutex
}
```

#### Concurrent Coupon Evaluation

The service evaluates multiple coupons concurrently for a single order:

```go
func (s *CouponRecommendationService) GetApplicableCoupons(order *Order) ([]*Coupon, error) {
    allCoupons := s.repository.GetAllCoupons()
    
    resultChan := make(chan couponResult, len(allCoupons))
    var wg sync.WaitGroup
    
    // Evaluate each coupon concurrently
    for _, coupon := range allCoupons {
        wg.Add(1)
        go func(c *Coupon) {
            defer wg.Done()
            applicable := s.isCouponApplicable(order, c)
            resultChan <- couponResult{coupon: c, applicable: applicable}
        }(coupon)
    }
    
    // Collect results
    go func() {
        wg.Wait()
        close(resultChan)
    }()
    
    // Filter applicable coupons
    var applicableCoupons []*Coupon
    for result := range resultChan {
        if result.applicable {
            applicableCoupons = append(applicableCoupons, result.coupon)
        }
    }
    
    return applicableCoupons, nil
}
```

**Performance Benefits**:

- Parallel evaluation of all coupons
- Non-blocking channel operations
- Efficient result collection

---

## Concurrency Handling

### Thread Safety Mechanisms

1. **Read-Write Mutex (`sync.RWMutex`)**:
   - Allows multiple concurrent readers
   - Exclusive access for writers
   - Used in Repository and Service layers

2. **Wait Groups (`sync.WaitGroup`)**:
   - Synchronizes goroutine completion
   - Ensures all rule evaluations complete before proceeding

3. **Channels**:
   - Safe communication between goroutines
   - Buffered channels for non-blocking operations

### Concurrent Request Handling

The system handles multiple concurrent requests safely:

```go
func testConcurrentRequests(service *CouponRecommendationService) {
    orders := createSampleOrders(10)
    var wg sync.WaitGroup
    
    for i, order := range orders {
        wg.Add(1)
        go func(idx int, o *Order) {
            defer wg.Done()
            coupons, _ := service.GetApplicableCoupons(o)
            // Process results
        }(i, order)
    }
    
    wg.Wait()
}
```

---

## Rule Formula System

The system supports logical formulas for rule combinations:

- **AND(0, 1)**: Both rule 0 and rule 1 must pass
- **OR(0, 1)**: Either rule 0 or rule 1 must pass
- **0**: Single rule (rule at index 0)

Example:

```go
coupon := &Coupon{
    CouponID:       "C1",
    DiscountAmount: 100.0,
    Rules:          []Rule{userTypeRule, productTypeRule},
    RulesFormula:   "AND(0, 1)", // Both rules must pass
}
```

---

## Usage Example

```go
// Initialize
repository := NewInMemoryCouponRepository()
service := NewCouponRecommendationService(repository)

        // Create rules
userTypeRule := &UserTypeRule{RequiredUserType: UserTypePremium}
productTypeRule := &ProductTypeRule{RequiredProductType: ProductTypeGroceries}

// Create coupon
coupon := &Coupon{
    CouponID:       "C1",
    DiscountAmount: 100.0,
    Rules:          []Rule{userTypeRule, productTypeRule},
    RulesFormula:   "AND(0, 1)",
}

// Add coupon
service.AddCoupon(coupon)

// Get applicable coupons
applicableCoupons, err := service.GetApplicableCoupons(order)

// Get best coupon (highest discount)
bestCoupon, err := service.GetBestCoupon(order)
```

---

## Performance Optimizations

1. **Concurrent Rule Evaluation**: Rules are evaluated in parallel using goroutines
2. **Concurrent Coupon Evaluation**: Multiple coupons are checked simultaneously
3. **Read-Write Locks**: Multiple readers can access data concurrently
4. **Buffered Channels**: Non-blocking communication between goroutines
5. **Efficient Data Structures**: Maps for O(1) lookups where applicable

---

## Testing

The implementation includes:

1. **Single Order Test**: Basic functionality verification
2. **Concurrent Request Test**: Tests thread safety with multiple simultaneous requests
3. **Performance Test**: Measures throughput and average response time

Run the tests:

```bash
go run main.go
```

---

## Key Takeaways

1. **Design Patterns**: Strategy, Repository, and Service patterns ensure clean architecture
2. **Thread Safety**: Proper use of mutexes, wait groups, and channels
3. **Concurrency**: Parallel evaluation of rules and coupons for better performance
4. **Extensibility**: Easy to add new rule types by implementing the `Rule` interface
5. **Maintainability**: Clear separation of concerns and well-defined interfaces

---

## Future Enhancements

1. **Caching**: Cache rule evaluation results for frequently accessed orders
2. **Rule Engine**: Implement a more sophisticated formula parser (e.g., using ANTLR)
3. **Database Integration**: Replace in-memory repository with database-backed storage
4. **Metrics & Monitoring**: Add observability for production use
5. **Rate Limiting**: Implement rate limiting for API endpoints
6. **Coupon Expiry**: Add expiration date handling
7. **Usage Limits**: Track and enforce coupon usage limits per user

---

## Conclusion

This implementation provides a robust, thread-safe, and performant coupon recommendation system. The use of design patterns ensures maintainability, while concurrent evaluation optimizes performance. The system is ready for production use with proper error handling and extensibility built-in.
