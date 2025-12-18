# Bloom Filter - Low-Level Design

A thread-safe, space-efficient Bloom Filter implementation in Go.

## Overview

A **Bloom Filter** is a probabilistic data structure used to test whether an element is a member of a set. It is designed to be very fast and extremely space-efficient, especially when working with large volumes of data.

### Key Properties

- **Space Efficient**: Uses significantly less memory than storing actual elements
- **Fast Operations**: Both `Add` and `MightContain` operations are O(1)
- **No False Negatives**: If it says an element is not present, it's guaranteed to be absent
- **Possible False Positives**: May incorrectly say an element is present when it isn't
- **Thread-Safe**: Safe to use in concurrent environments

## Design Components

### 1. BitSet
A thread-safe bit array implementation that efficiently stores bits using `uint64` slices. Each `uint64` holds 64 bits, providing compact storage.

**Key Methods:**
- `Set(index)`: Sets a bit at the given index to 1
- `Get(index)`: Returns true if the bit at the given index is set
- `Size()`: Returns the size of the bit array

### 2. HashStrategy Interface
Defines the contract for hash functions used in the Bloom Filter. Different implementations provide different hashing algorithms.

**Implementations:**
- **FNV1aHashStrategy**: FNV-1a (64-bit) hashing algorithm - fast and simple
- **DJB2HashStrategy**: DJB2 hash function - known for simplicity
- **MurmurHashStrategy**: Simplified MurmurHash-like algorithm - good distribution

### 3. BloomFilter
The main class that exposes the public API and manages internal data structures.

**Key Methods:**
- `Add(element)`: Adds an element to the Bloom Filter
- `MightContain(element)`: Checks if an element might be in the set
- `Size()`: Returns the size of the bit array
- `NumHashFunctions()`: Returns the number of hash functions used

### 4. BloomFilterBuilder
Provides a fluent API for building BloomFilter instances with optimal configuration.

**Builder Methods:**
- `WithSize(size)`: Sets the size of the bit array (m)
- `WithNumHashFunctions(k)`: Sets the number of hash functions (k)
- `WithHashStrategy(strategy)`: Sets the hash strategy to use
- `WithExpectedItems(n, falsePositiveRate)`: Calculates optimal m and k based on expected items and desired false positive rate
- `Build()`: Creates the BloomFilter instance

## Usage Examples

### Basic Usage

```go
// Create a Bloom Filter with manual configuration
bf, err := NewBloomFilterBuilder().
    WithSize(1000).
    WithNumHashFunctions(3).
    WithHashStrategy(NewFNV1aHashStrategy()).
    Build()

// Add elements
bf.Add("apple")
bf.Add("banana")

// Check for presence
mightContain, _ := bf.MightContain("apple")  // true
mightContain, _ := bf.MightContain("grape")  // false (definitely not present)
```

### Optimal Configuration

```go
// Create a Bloom Filter optimized for 1000 items with 1% false positive rate
bf, err := NewBloomFilterBuilder().
    WithExpectedItems(1000, 0.01).
    Build()

// The builder automatically calculates:
// - Optimal bit array size (m)
// - Optimal number of hash functions (k)
```

### Thread-Safe Concurrent Access

```go
// The Bloom Filter is thread-safe and can be used concurrently
var wg sync.WaitGroup

// Concurrent writes
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        bf.Add(fmt.Sprintf("element-%d", id))
    }(i)
}

// Concurrent reads
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        bf.MightContain(fmt.Sprintf("element-%d", id))
    }(i)
}

wg.Wait()
```

## How It Works

1. **Initialization**: A Bloom Filter is initialized with:
   - `m`: Size of the bit array
   - `k`: Number of hash functions

2. **Adding Elements**:
   - When an element is added, `k` hash functions are applied to it
   - Each hash function produces an index in the range `[0, m-1]`
   - The bits at all `k` positions are set to 1

3. **Checking Membership**:
   - When checking if an element might be present:
     - Apply the same `k` hash functions
     - Check if all `k` bits are set to 1
     - If any bit is 0, the element is **definitely not** present
     - If all bits are 1, the element **might be** present (possible false positive)

## Optimal Configuration

The optimal values for `m` (bit array size) and `k` (number of hash functions) can be calculated based on:
- `n`: Expected number of elements
- `p`: Desired false positive rate

**Formulas:**
- `m = -n * ln(p) / (ln(2)^2)`
- `k = (m/n) * ln(2)`

The `WithExpectedItems()` builder method automatically calculates these optimal values.

## Design Patterns Used

1. **Strategy Pattern**: `HashStrategy` interface allows different hashing algorithms to be used interchangeably
2. **Builder Pattern**: `BloomFilterBuilder` provides a fluent API for configuration
3. **Facade Pattern**: `BloomFilter` class hides internal complexity and exposes a simple API

## Thread Safety

The implementation uses `sync.RWMutex` to ensure thread-safety:
- Write operations (`Add`) use write locks
- Read operations (`MightContain`) use read locks
- This allows multiple concurrent reads while ensuring exclusive writes

## Time and Space Complexity

- **Time Complexity**: 
  - `Add`: O(k) where k is the number of hash functions
  - `MightContain`: O(k)
  - Since k is typically small (3-10), operations are effectively O(1)

- **Space Complexity**: 
  - O(m) where m is the size of the bit array
  - Much more space-efficient than storing actual elements

## Trade-offs

- **Space vs Accuracy**: Smaller bit arrays increase false positive rate
- **Hash Functions**: More hash functions reduce false positives but increase computation time
- **No Deletion**: Standard Bloom Filters don't support deletion (would require Counting Bloom Filter)

## Running the Demo

```bash
cd go/bloom-filter
go run .
```

This will run several demos:
1. Basic usage with manual configuration
2. Builder with expected items and false positive rate
3. Demonstration of false positives
4. Thread-safety demonstration

## Future Enhancements

- Counting Bloom Filter (supports deletion)
- Scalable Bloom Filter (grows dynamically)
- Partitioned Bloom Filter (better cache performance)
- Serialization support for persistence

