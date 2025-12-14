# Analytical Database

## Problem Statement

Build an analytical database system that works with records represented as maps with string keys and integer values. The database supports functions to find records based on key values and sorting directions.

This is a multi-part problem that builds incrementally:

### Part 1: min_by_key

Implement `min_by_key` function that scans an array of records and returns the record with the minimum value for a specified key.

**Rules:**
- Records that do not contain the specified key are considered to have value 0 for that key
- Keys may map to negative values
- If several records share the same minimum value, return any of them
- Handle empty arrays idiomatically

### Part 2: first_by_key

Implement `first_by_key` function that takes:
- A string key
- A string sort direction ("asc" or "desc")
- An array of records

**Rules:**
- If direction is "asc", return the minimum record
- If direction is "desc", return the maximum record
- Records without the key are treated as having value 0
- Re-implement `min_by_key` using `first_by_key`

### Part 3: RecordComparator

Implement a `RecordComparator` class that:
- Constructor accepts: key (string) and direction (string)
- `compare` method takes two records and returns:
  - `-1` if first record comes before second record
  - `0` if records are equal
  - `1` if first record comes after second record
- Use the comparator in `first_by_key` implementation

Reference: [Stripe Technical Screen](https://github.com/freecodemessiah/stripe-interview/blob/master/Stripe%20Technical%20Screen.ipynb)

## Test Cases

### Part 1: min_by_key

#### Test Case 1: Basic minimum
- **Input**: `min_by_key("a", [{"a": 1, "b": 2}, {"a": 2}])`
- **Expected Output**: `{"a": 1, "b": 2}`

#### Test Case 2: Different order
- **Input**: `min_by_key("a", [{"a": 2}, {"a": 1, "b": 2}])`
- **Expected Output**: `{"a": 1, "b": 2}`

#### Test Case 3: Missing key (treated as 0)
- **Input**: `min_by_key("b", [{"a": 1, "b": 2}, {"a": 2}])`
- **Expected Output**: `{"a": 2}` (missing "b" key = 0, which is minimum)

#### Test Case 4: Empty record
- **Input**: `min_by_key("a", [{}])`
- **Expected Output**: `{}`

#### Test Case 5: Negative values
- **Input**: `min_by_key("b", [{"a": -1}, {"b": -1}])`
- **Expected Output**: `{"b": -1}`

### Part 2: first_by_key

#### Test Case 6: Ascending
- **Input**: `first_by_key("a", "asc", [{"a": 1}])`
- **Expected Output**: `{"a": 1}`

#### Test Case 7: Missing key in ascending
- **Input**: `first_by_key("a", "asc", [{"b": 1}, {"b": -2}, {"a": 10}])`
- **Expected Output**: `{"b": 1}` or `{"b": -2}` (both have a=0, which is minimum)

#### Test Case 8: Descending
- **Input**: `first_by_key("a", "desc", [{"b": 1}, {"b": -2}, {"a": 10}])`
- **Expected Output**: `{"a": 10}`

#### Test Case 9: Different key
- **Input**: `first_by_key("b", "asc", [{"b": 1}, {"b": -2}, {"a": 10}])`
- **Expected Output**: `{"b": -2}`

- **Input**: `first_by_key("b", "desc", [{"b": 1}, {"b": -2}, {"a": 10}])`
- **Expected Output**: `{"b": 1}`

#### Test Case 10: Empty records
- **Input**: `first_by_key("a", "desc", [{}, {"a": 10, "b": -10}, {}, {"a": 3, "c": 3}])`
- **Expected Output**: `{"a": 10, "b": -10}`

### Part 3: RecordComparator

#### Test Case 11: Comparator ascending
- **Input**: 
  - `cmp = Comparator("a", "asc")`
  - `cmp.compare({"a": 1}, {"a": 2})`
- **Expected Output**: `-1`

- **Input**: `cmp.compare({"a": 2}, {"a": 1})`
- **Expected Output**: `1`

- **Input**: `cmp.compare({"a": 1}, {"a": 1})`
- **Expected Output**: `0`

#### Test Case 12: Comparator descending
- **Input**: 
  - `cmp = Comparator("a", "desc")`
  - `cmp.compare({"a": 1}, {"a": 2})`
- **Expected Output**: `1` (reversed from ascending)

## Running the Solutions

### Go

To run the Go solution:

```bash
cd go/stripe/analytical_database
go run main.go
```

### C++

To compile and run the C++ solution:

```bash
cd go/stripe/analytical_database
g++ -o main main.cpp
./main
```

