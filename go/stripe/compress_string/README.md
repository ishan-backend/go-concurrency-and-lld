# Compress String

## Problem Statement

Given a String, split it into major parts separated by special char '/'. For each major part that's split by '/', we can further split it into minor parts separated by '.'.

The compression works as follows:
- For each word in the minor parts, compress it by keeping the first and last character, and replacing the middle characters with their count.
- Example: "stripe" becomes "s4e" (s + 4 middle chars + e)
- Words with 2 or fewer characters remain unchanged.

The compression happens in two phases:
1. **Phase 1**: Compress the specified number of minor parts from the end in each major part
2. **Phase 2**: Compress all minor parts in all major parts

### Example 1
```
str = stripe.com/payments/checkout/customer.john.doe
minor_parts = 2
after Part 1 compression
=>
s4e.c1m/p6s/c6t/c6r.j2n.d1e
after Part 2 compression
=>
s4e.c1m/p6s/c6t/c6r.j5e
```

### Example 2
```
str = www.api.stripe.com/checkout
minor_parts = 3
(after Part 1 compression)
=>
w1w.a1i.s4e.c1m/c6t
(then after Part 2 compression)
=>
w1w.a1i.s7m/c6t
```

Reference: https://leetcode.com/discuss/post/4057853/stripe-phone-interview-by-anonymous_user-v8wr/

## Test Cases

### Test Case 1
- **Input**: `stripe.com/payments/checkout/customer.john.doe`, `minor_parts = 2`
- **Expected Output**: `s4e.c1m/p6s/c6t/c6r.j5e`

### Test Case 2
- **Input**: `www.api.stripe.com/checkout`, `minor_parts = 3`
- **Expected Output**: `w1w.a1i.s7m/c6t`

### Test Case 3
- **Input**: `a.b.c/d.e.f`, `minor_parts = 2`
- **Expected Output**: `a.b.c/d.e.f` (words too short to compress)

## Running the Solutions

### Go

To run the Go solution:

```bash
cd go/stripe/compress_string
go run main.go
```

### C++

To compile and run the C++ solution:

```bash
cd go/stripe/compress_string
g++ -o main main.cpp
./main
```