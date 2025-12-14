# Server Penalty

## Problem Statement

Implement a system to calculate penalties for servers based on their performance metrics. Servers can have various events such as:

1. **Response Time Events**: When a server's response time exceeds a threshold
2. **Error Events**: When a server encounters errors

### Penalty Calculation Rules

- **Response Time Penalty**: 
  - If response time > threshold: `penalty = (response_time - threshold) * multiplier`
  - Example: If threshold is 100ms, multiplier is $0.50/ms, and response time is 150ms:
    - Penalty = (150 - 100) * 0.50 = $25

- **Error Penalty**: 
  - Each error incurs a fixed penalty: `penalty = error_count * error_penalty`
  - Example: If error penalty is $10 and there are 3 errors:
    - Penalty = 3 * $10 = $30

- **Total Penalty**: Sum of all penalties for a server

## Requirements

1. Add events for servers (response time or errors)
2. Calculate penalty for a specific server
3. Get penalties for all servers
4. Get servers sorted by penalty (highest first)
5. Optionally filter events by time window (e.g., last 5 minutes)

## Test Cases

### Test Case 1: Response Time Penalties
- **Input**:
  - Server1: response_time = 150ms, 120ms (threshold = 100ms, multiplier = $0.50/ms)
  - Server2: response_time = 200ms, 110ms
  - Server3: response_time = 90ms
- **Expected Output**:
  - Server1: $35.00 (50ms * $0.50 + 20ms * $0.50)
  - Server2: $55.00 (100ms * $0.50 + 10ms * $0.50)
  - Server3: $0.00 (all under threshold)

### Test Case 2: Combined Response Time and Errors
- **Input**:
  - ServerA: response_time = 150ms, 2 errors (error_penalty = $10)
  - ServerB: response_time = 120ms, 130ms
  - ServerC: 3 errors
- **Expected Output**:
  - ServerA: $45.00 (50ms * $0.50 + 2 * $10)
  - ServerB: $25.00 (20ms * $0.50 + 30ms * $0.50)
  - ServerC: $30.00 (3 * $10)

### Test Case 3: Servers Sorted by Penalty
- **Input**: Same as Test Case 2
- **Expected Output**: 
  - ServerA: $45.00 (highest)
  - ServerC: $30.00
  - ServerB: $25.00 (lowest)

### Test Case 4: Time-Windowed Penalties
- **Input**:
  - ServerX: response_time = 150ms (2 min ago), 200ms (10 min ago)
  - ServerY: 2 errors (1 min ago, 6 min ago)
  - Time window: last 5 minutes
- **Expected Output**:
  - ServerX: $25.00 (only 150ms event counted)
  - ServerY: $10.00 (only 1 error counted)

## Running the Solutions

### Go

To run the Go solution:

```bash
cd go/stripe/server_penalty
go run main.go
```

### C++

To compile and run the C++ solution:

```bash
cd go/stripe/server_penalty
g++ -o main main.cpp -std=c++11
./main
```

**Note**: C++ solution requires C++11 or later for chrono library support.

