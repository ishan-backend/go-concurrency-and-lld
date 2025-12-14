# Radar Rules Engine

## Problem Statement

Implement a rule engine for Stripe Radar that processes payment charges and applies ALLOW/BLOCK rules based on conditions. The system should:

1. Parse charge information from a formatted string
2. Parse and evaluate rules (ALLOW/BLOCK) with conditions
3. Return 1 if the charge is allowed, 0 if blocked
4. Process rules in order (first matching rule wins)

### Input Format

The input is a comma-separated string containing:
1. **CHARGE**: Charge information in format `CHARGE: card_country=X&currency=Y&amount=Z&ip_country=W`
2. **Rules**: One or more rules in format `ALLOW:condition` or `BLOCK:condition`

### Rule Conditions

Conditions can include:
- **Amount comparisons**: `amount > 100`, `amount < 50`, `amount == 100`, `amount != 200`
- **Country comparisons**: `card_country != ip_country`, `card_country == ip_country`
- **Compound conditions**: `card_country != ip_country AND amount > 100`

### Processing Logic

1. Rules are processed in order
2. If an ALLOW rule's condition is met, return 1 (allowed)
3. If a BLOCK rule's condition is met, return 0 (blocked)
4. If no rules match, default to 1 (allowed)

Reference: [RadarRules.java](https://github.com/kylelong/stripe-interview/blob/master/RadarRules.java)

## Test Cases

### Test Case 1: Basic Block
- **Input**: `"CHARGE: card_country=US&currency=USD&amount=150&ip_country=CA,BLOCK:amount > 100"`
- **Expected Output**: `0` (blocked because amount 150 > 100)

### Test Case 2: Compound Block
- **Input**: `"CHARGE: card_country=US&currency=USD&amount=150&ip_country=CA,ALLOW:amount<100,BLOCK:card_country != ip_country AND amount > 100"`
- **Expected Output**: `0` (blocked: card_country(US) != ip_country(CA) AND amount(150) > 100)
- **Note**: The ALLOW rule doesn't match (150 is not < 100), so the BLOCK rule is evaluated

### Test Case 3: ALLOW Rule Matches
- **Input**: `"CHARGE: card_country=US&currency=USD&amount=50&ip_country=CA,ALLOW:amount<100,BLOCK:card_country != ip_country AND amount > 100"`
- **Expected Output**: `1` (allowed because amount 50 < 100, ALLOW rule matches first)

### Test Case 4: Country Mismatch Only
- **Input**: `"CHARGE: card_country=US&currency=USD&amount=50&ip_country=CA,BLOCK:card_country != ip_country"`
- **Expected Output**: `0` (blocked: US != CA)

### Test Case 5: No Rules Match
- **Input**: `"CHARGE: card_country=US&currency=USD&amount=50&ip_country=US,ALLOW:amount>100,BLOCK:card_country != ip_country"`
- **Expected Output**: `1` (default allow, no rules matched)

### Test Case 6: Multiple BLOCK Rules
- **Input**: `"CHARGE: card_country=US&currency=USD&amount=200&ip_country=US,BLOCK:amount>150,BLOCK:amount<50"`
- **Expected Output**: `0` (blocked by first rule: amount 200 > 150)

### Test Case 7: Complex AND Condition
- **Input**: `"CHARGE: card_country=US&currency=USD&amount=120&ip_country=CA,BLOCK:card_country != ip_country AND amount > 100"`
- **Expected Output**: `0` (blocked: US != CA AND 120 > 100)

## Supported Operators

- **Comparison**: `>`, `<`, `>=`, `<=`, `==`, `!=`
- **Logical**: `AND` (both conditions must be true)

## Running the Solutions

### Go

To run the Go solution:

```bash
cd go/stripe/radar_rules
go run main.go
```

### C++

To compile and run the C++ solution:

```bash
cd go/stripe/radar_rules
g++ -o main main.cpp
./main
```

