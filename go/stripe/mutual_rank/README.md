# Mutual Rank

## Problem Statement

Given a set of user rankings where each user ranks other users in order of preference, calculate the mutual rank between pairs of users.

The mutual rank between two users A and B is defined as:
- **Mutual Rank = Rank(B in A's list) + Rank(A in B's list)**

Where rank is 1-indexed (first preference = rank 1, second = rank 2, etc.).

A lower mutual rank indicates a better mutual match (both users rank each other highly).

### Example
If Alice ranks: `[Bob, Charlie, Dave]`
And Bob ranks: `[Alice, Charlie, Dave]`

Then:
- Rank of Bob in Alice's list = 1
- Rank of Alice in Bob's list = 1\
 
- Mutual rank = 1 + 1 = 2

## Requirements

1. Add rankings for users (comma-separated list of user IDs in preference order)
2. Get the rank of a user in another user's ranking
3. Calculate mutual rank between two users
4. Find best matches for a user (sorted by mutual rank, lowest first)

## Test Cases

### Test Case 1: Symmetric Rankings
- **Input**:
  - Alice: `bob, charlie, dave`
  - Bob: `alice, charlie, dave`
  - Charlie: `alice, bob, dave`
  - Dave: `alice, bob, charlie`
- **Query**: Mutual rank between Alice and Bob
- **Expected Output**: 
  - Rank(Alice -> Bob) = 1
  - Rank(Bob -> Alice) = 1
  - Mutual rank = 2

### Test Case 2: Asymmetric Rankings
- **Input**:
  - User1: `user2, user3, user4`
  - User2: `user4, user1, user3`
  - User3: `user1, user2, user4`
  - User4: `user2, user1, user3`
- **Query**: Mutual rank between User1 and User2
- **Expected Output**:
  - Rank(User1 -> User2) = 1
  - Rank(User2 -> User1) = 2
  - Mutual rank = 3

### Test Case 3: One-Way Ranking
- **Input**:
  - A: `B, C`
  - B: `C, D` (B doesn't rank A)
- **Query**: Mutual rank between A and B
- **Expected Output**: No mutual ranking (returns false/error)

### Test Case 4: Best Matches
- **Input**: Same as Test Case 2
- **Query**: Find best matches for User1
- **Expected Output**: 
  - User3 (mutual rank = 3: User1 ranks User3 as 2, User3 ranks User1 as 1)
  - User2 (mutual rank = 3: User1 ranks User2 as 1, User2 ranks User1 as 2)
  - User4 (mutual rank = 4: User1 ranks User4 as 3, User4 ranks User1 as 1)

## Running the Solutions

### Go

To run the Go solution:

```bash
cd go/stripe/mutual_rank
go run main.go
```

### C++

To compile and run the C++ solution:

```bash
cd go/stripe/mutual_rank
g++ -o main main.cpp
./main
```

