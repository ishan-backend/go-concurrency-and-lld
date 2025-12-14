# User Points

## Problem Statement

Implement a user points management system that tracks points earned, spent, and expired for users. The system should support:

1. **Earning Points**: Users can earn points through various activities
2. **Spending Points**: Users can spend points (FIFO - oldest points first)
3. **Point Expiration**: Points can expire after a certain time period
4. **Balance Tracking**: Track current balance for each user
5. **Leaderboard**: Get top users by balance

### Key Features

- **FIFO Spending**: When spending points, use the oldest points first
- **Expiration Handling**: Points expire based on their expiration time
- **Transaction History**: Track all transactions for audit purposes
- **Balance Calculation**: Calculate current balance excluding expired points

## Requirements

1. Add transactions (earn, spend, expire)
2. Get balance for a specific user
3. Get balances for all users
4. Get top N users by balance
5. Get transaction history for a user
6. Handle point expiration automatically

## Test Cases

### Test Case 1: Basic Earn and Spend
- **Input**:
  - User1 earns 100 points
  - User1 earns 50 points
  - User1 spends 75 points
- **Expected Output**: Balance = 75 points (100 + 50 - 75)

### Test Case 2: Points Expiration
- **Input**:
  - User2 earns 100 points (expires in 30 days)
  - User2 earns 200 points (expires in 60 days)
  - After 31 days, check balance
- **Expected Output**: Balance = 200 points (100 expired)

### Test Case 3: FIFO Spending
- **Input**:
  - User3 earns 100 points
  - User3 earns 50 points
  - User3 earns 25 points
  - User3 spends 75 points
- **Expected Output**: Balance = 100 points
  - First 100 points are spent (75 from first batch, leaving 25)
  - Remaining: 25 + 50 + 25 = 100

### Test Case 4: Top Users
- **Input**:
  - Alice earns 500 points
  - Bob earns 300 points, then 200 points
  - Charlie earns 100 points, then spends 50 points
- **Expected Output**:
  - Top 3 users:
    1. Bob: 500 points
    2. Alice: 500 points
    3. Charlie: 50 points

### Test Case 5: Transaction History
- **Input**: Multiple transactions for a user
- **Expected Output**: List of all transactions in chronological order

## Running the Solutions

### Go

To run the Go solution:

```bash
cd go/stripe/user_points
go run main.go
```

### C++

To compile and run the C++ solution:

```bash
cd go/stripe/user_points
g++ -o main main.cpp -std=c++11
./main
```

**Note**: C++ solution requires C++11 or later for chrono library support.

