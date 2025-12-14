# Stripe Capital

## Problem Statement

Implement a capital financing system that helps determine loan eligibility for businesses based on their transaction history. Stripe Capital offers financing to businesses using Stripe, and the system should:

1. Track financial transactions over time
2. Calculate average monthly revenue
3. Determine loan eligibility based on revenue thresholds
4. Calculate potential loan amounts
5. Provide monthly revenue breakdowns
6. Identify top-performing months

### Key Features

- **Transaction Tracking**: Record all financial transactions with amounts and dates
- **Revenue Calculation**: Calculate average monthly revenue from positive transactions
- **Eligibility Check**: Determine if a business qualifies for a loan based on revenue threshold
- **Loan Amount Calculation**: Calculate potential loan amount as a multiple of monthly revenue
- **Monthly Breakdown**: Group transactions by month and calculate monthly revenue
- **Top Months**: Identify the best-performing months by revenue

## Requirements

1. Add transactions with amount, date, and description
2. Calculate total capital from all transactions
3. Calculate average monthly revenue
4. Check eligibility for loan (revenue >= threshold)
5. Calculate potential loan amount (revenue * multiplier)
6. Get monthly revenue breakdown
7. Get top N months by revenue

## Test Cases

### Test Case 1: Basic Eligibility Check
- **Input**:
  - Transactions: $12,000, $15,000, $11,000, $13,000 (one per month)
  - Threshold: $10,000
- **Expected Output**:
  - Average Monthly Revenue: $12,750
  - Eligible: true

### Test Case 2: Loan Amount Calculation
- **Input**: Average monthly revenue = $12,750, multiplier = 3x
- **Expected Output**: Loan Amount = $38,250

### Test Case 3: Monthly Breakdown
- **Input**: Transactions across multiple months
- **Expected Output**: Map of month -> revenue for each month

### Test Case 4: Top Months
- **Input**: Multiple months with varying revenue
- **Expected Output**: List of top N months sorted by revenue (descending)

### Test Case 5: Not Eligible
- **Input**:
  - Transactions: $5,000, $6,000, $4,000 (one per month)
  - Threshold: $10,000
- **Expected Output**:
  - Average Monthly Revenue: $5,000
  - Eligible: false

## Running the Solutions

### Go

To run the Go solution:

```bash
cd go/stripe/stripe_capital
go run main.go
```

### C++

To compile and run the C++ solution:

```bash
cd go/stripe/stripe_capital
g++ -o main main.cpp -std=c++11
./main
```

**Note**: C++ solution requires C++11 or later for chrono library support.

