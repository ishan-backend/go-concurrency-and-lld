# Find Pairs

## Problem Statement

Given an array of integers and a target sum, find all pairs of indices (i, j) where i < j and the values at those indices sum to the target.

This is a variation of the classic "Two Sum" problem where we need to find all pairs, not just one.

### Key Features

- **Find All Pairs**: Return all pairs of indices whose values sum to target
- **Avoid Duplicates**: Each pair (i, j) is only included once where i < j
- **Handle Duplicates**: Correctly handle arrays with duplicate values
- **Multiple Approaches**: Implement both hash map and two-pointer approaches

## Requirements

1. Find all pairs of indices (i, j) where numbers[i] + numbers[j] == target and i < j
2. Return pairs sorted by first index, then second index
3. Count the number of pairs
4. Handle duplicate values correctly
5. Support two-pointer approach for sorted arrays

## Test Cases

### Test Case 1: Basic Two Sum
- **Input**: `numbers = [2, 7, 11, 15]`, `target = 9`
- **Expected Output**: `[(0, 1)]` (indices 0 and 1: 2 + 7 = 9)

### Test Case 2: Multiple Pairs
- **Input**: `numbers = [1, 2, 3, 4, 5, 6]`, `target = 7`
- **Expected Output**: 
  - `[(0, 5)]` -> 1 + 6 = 7
  - `[(1, 4)]` -> 2 + 5 = 7
  - `[(2, 3)]` -> 3 + 4 = 7
- **Count**: 3 pairs

### Test Case 3: Duplicate Values
- **Input**: `numbers = [3, 3, 4, 5, 6, 3]`, `target = 6`
- **Expected Output**: 
  - `[(0, 1)]` -> 3 + 3 = 6
  - `[(0, 5)]` -> 3 + 3 = 6
  - `[(1, 5)]` -> 3 + 3 = 6
- **Note**: All pairs with duplicate values should be included

### Test Case 4: Two-Pointer Approach
- **Input**: `numbers = [1, 3, 4, 5, 6, 8, 9]`, `target = 10`
- **Expected Output**:
  - `[(1, 4)]` -> 3 + 7 = 10 (if 7 exists)
  - Or other valid pairs summing to 10

### Test Case 5: No Pairs Found
- **Input**: `numbers = [1, 2, 3, 4]`, `target = 10`
- **Expected Output**: `[]` (empty list, no pairs found)

## Algorithm Approaches

### Approach 1: Hash Map (O(n) time, O(n) space)
1. Create a map: value -> list of indices
2. For each number, calculate complement = target - number
3. Look up complement in map and create pairs
4. Ensure i < j to avoid duplicates

### Approach 2: Two Pointers (O(n log n) time, O(1) space)
1. Sort array while preserving original indices
2. Use two pointers (left, right)
3. If sum < target, move left pointer right
4. If sum > target, move right pointer left
5. If sum == target, add pair and move both pointers

## Running the Solutions

### Go

To run the Go solution:

```bash
cd go/stripe/find_pairs
go run main.go
```

### C++

To compile and run the C++ solution:

```bash
cd go/stripe/find_pairs
g++ -o main main.cpp
./main
```

