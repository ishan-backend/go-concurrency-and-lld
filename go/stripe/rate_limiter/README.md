# Rate Limiter

## Problem Statement

Implement a rate limiter that restricts the number of requests a user can make within a given time window. The system should support two algorithms:

1. **Sliding Window**: Tracks requests in a fixed time window
2. **Token Bucket**: Uses tokens that refill over time

### Sliding Window Algorithm
- Maintains a list of request timestamps for each user
- Removes requests outside the time window
- Allows requests if count < maxRequests
- Example: 5 requests per 10 seconds

### Token Bucket Algorithm
- Each user has a bucket with tokens
- Each request consumes 1 token
- Tokens refill at a constant rate over time
- Allows requests if tokens >= 1
- Example: Max 5 tokens, refill 1 token per second

## Requirements

1. Check if a request is allowed for a user
2. Get remaining requests/tokens for a user
3. Support multiple users independently
4. Thread-safe implementation
5. Handle time-based expiration/refill

## Test Cases

### Test Case 1: Sliding Window - Basic
- **Input**: Limit = 5 requests per 10 seconds
- **Actions**: 
  - Make 7 requests quickly
- **Expected Output**:
  - First 5 requests: Allowed = true
  - Requests 6-7: Allowed = false
  - Remaining decreases from 4 to 0

### Test Case 2: Sliding Window - Rate Limit Exceeded
- **Input**: Limit = 3 requests per 5 seconds
- **Actions**: Make 5 requests immediately
- **Expected Output**:
  - First 3 requests: Allowed = true
  - Requests 4-5: Allowed = false

### Test Case 3: Multiple Users
- **Input**: Limit = 2 requests per 3 seconds per user
- **Actions**: 
  - Alice makes 3 requests
  - Bob makes 2 requests
- **Expected Output**:
  - Alice: First 2 allowed, 3rd blocked
  - Bob: Both requests allowed (independent limits)

### Test Case 4: Token Bucket - Basic
- **Input**: Max tokens = 5, Refill rate = 1 token/second
- **Actions**: Make 7 requests
- **Expected Output**:
  - First 5 requests: Allowed = true
  - Requests 6-7: Allowed = false initially
  - After waiting, tokens refill and requests become allowed

### Test Case 5: Token Bucket - Slow Refill
- **Input**: Max tokens = 3, Refill rate = 0.5 tokens/second
- **Actions**: 
  - Make 3 requests quickly (all allowed)
  - Wait 3 seconds
  - Make 1 more request
- **Expected Output**:
  - First 3 requests: Allowed = true
  - After 3 seconds: ~1.5 tokens refilled
  - Next request: Allowed = true (1 token available)

## Running the Solutions

### Go

To run the Go solution:

```bash
cd go/stripe/rate_limiter
go run main.go
```

### C++

To compile and run the C++ solution:

```bash
cd go/stripe/rate_limiter
g++ -o main main.cpp -std=c++11 -pthread
./main
```

**Note**: 
- C++ solution requires C++11 or later
- Requires `-pthread` flag for mutex support
- Uses `std::this_thread::sleep_for` for timing

## Algorithm Comparison

| Feature | Sliding Window | Token Bucket |
|---------|---------------|--------------|
| Burst Handling | Limited by window | Allows bursts up to max tokens |
| Smooth Rate | Yes | Yes (with refill) |
| Memory Usage | O(requests in window) | O(1) per user |
| Complexity | O(n) cleanup | O(1) per request |

