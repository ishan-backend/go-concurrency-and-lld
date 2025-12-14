# Cost of Shipping

## Problem Statement

Given an input string containing shipping routes in the format: `"origin,destination,carrier,cost:origin,destination,carrier,cost:..."`, implement a system to:

1. **Part 1**: Find the direct cost to ship an item from origin X to destination Y
2. **Part 2**: Find the cost if you need to use intermediate steps (e.g., US -> CA -> UK)

### Input Format
The input string contains multiple routes separated by `:`. Each route has the format:
`origin,destination,carrier,cost`

Example: `"US,UK,UPS,5:US,CA,FedEx,3:CA,UK,DHL,7"`

This means:
- US -> UK via UPS costs $5
- US -> CA via FedEx costs $3
- CA -> UK via DHL costs $7

Reference: https://leetcode.com/discuss/post/5883672/stripe-phone-screen-by-anonymous_user-0kk5/

## Test Cases

### Test Case 1: Direct Routes
- **Input**: `"US,UK,UPS,5:US,CA,FedEx,3:CA,UK,DHL,7"`
- **Query**: US -> UK
- **Expected Output**: Cost = $5 (Carrier: UPS)

### Test Case 2: No Direct Route
- **Input**: `"US,UK,UPS,5:US,CA,FedEx,3:CA,UK,DHL,7"`
- **Query**: UK -> US
- **Expected Output**: No direct route found

### Test Case 3: Intermediate Route
- **Input**: `"US,UK,UPS,5:US,CA,FedEx,3:CA,UK,DHL,7"`
- **Query**: US -> UK (with intermediate steps)
- **Expected Output**: 
  - Direct: US -> UK = $5
  - Via intermediate: US -> CA ($3) -> UK ($7) = Total $10

## Running the Solutions

### Go

To run the Go solution:

```bash
cd go/stripe/cost_of_shipping
go run main.go
```

### C++

To compile and run the C++ solution:

```bash
cd go/stripe/cost_of_shipping
g++ -o main main.cpp
./main
```
