# Splitwise-like Expense Sharing System - Low Level Design

## Table of Contents
1. [Requirements & Use Cases](#requirements--use-cases)
2. [Actors](#actors)
3. [Data Model](#data-model)
4. [APIs](#apis)
5. [Follow-up Features](#follow-up-features)
6. [Design Decisions](#design-decisions)

List requirements/use-cases
	•	List all the actors
	•	Define the data model (classes like User, Expense, Group, etc)
	•	Write an API for Add expense to a Group (Need to handle the expense type, split type, users involved etc)
	•	Write an API to view balances for a given user (globally and within a group)
	•	Follow up (to understand if their LLD is extensible or not)
	◦	Introduce simplify debts feature
	◦	Custom Split default strategy (example: 60%:40% split)
	◦	Tracking activity

---

## Requirements & Use Cases

### Core Requirements
1. **User Management**
   - Users can register and create accounts
   - Users can join groups
   - Users can be part of multiple groups

2. **Group Management**
   - Users can create groups
   - Users can add/remove members from groups
   - Groups can have custom default split strategies

3. **Expense Management**
   - Users can add expenses to groups
   - Support multiple expense types (EQUAL, EXACT, PERCENTAGE, SHARES)
   - Support multiple split types (EQUAL, UNEQUAL, PERCENTAGE, SHARES)
   - Track who paid and who owes what

4. **Balance Tracking**
   - View balances globally (across all groups)
   - View balances within a specific group
   - Track who owes whom

5. **Debt Simplification**
   - Automatically simplify debts to minimize transactions
   - Reduce the number of payments needed to settle all debts

6. **Activity Tracking**
   - Track all expense additions
   - Track all settlements
   - Track group membership changes

7. **Custom Split Strategies**
   - Groups can have default split strategies (e.g., 60:40 split)
   - Apply default strategy when adding expenses

---

## Actors

1. **User**
   - Can create groups
   - Can add expenses
   - Can view balances
   - Can settle debts
   - Can view activity

2. **Group Admin** (subset of User)
   - All User capabilities
   - Can add/remove members
   - Can set default split strategy

3. **System**
   - Manages balances
   - Simplifies debts
   - Tracks activities

---

## Data Model

### User
```go
type User struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}
```

### Group
```go
type Group struct {
    ID                string                 `json:"id"`
    Name              string                 `json:"name"`
    Members           []string               `json:"members"`
    Admin             string                 `json:"admin"`
    DefaultSplitType  SplitType              `json:"default_split_type"`
    DefaultSplitRatio map[string]float64     `json:"default_split_ratio"` // e.g., {"user1": 0.6, "user2": 0.4}
    CreatedAt         time.Time              `json:"created_at"`
}
```

### Expense
```go
type Expense struct {
    ID          string             `json:"id"`
    GroupID     *string            `json:"group_id,omitempty"`
    Title       string             `json:"title"`
    Description string             `json:"description,omitempty"`
    Amount      float64            `json:"amount"`
    Payer       string             `json:"payer"`
    CreatedBy   string             `json:"created_by"`
    ExpenseType ExpenseType        `json:"expense_type"`
    SplitType   SplitType          `json:"split_type"`
    Shares      map[string]float64 `json:"shares"` // user_id -> amount/percentage/share
    CreatedAt   time.Time          `json:"created_at"`
}
```

### ExpenseType
- `EQUAL`: Split equally among all participants
- `EXACT`: Exact amounts specified for each user
- `PERCENTAGE`: Split by percentage
- `SHARES`: Split by shares (e.g., 2:1:1)

### SplitType
- `EQUAL`: Equal split
- `UNEQUAL`: Unequal split (exact amounts)
- `PERCENTAGE`: Percentage-based split
- `SHARES`: Share-based split

### Balance
```go
type Balance struct {
    UserID    string  `json:"user_id"`
    OwedTo    string  `json:"owed_to"`
    Amount    float64 `json:"amount"` // positive = owes, negative = is owed
    GroupID   *string `json:"group_id,omitempty"`
}
```

### Activity
```go
type Activity struct {
    ID          string      `json:"id"`
    Type        ActivityType `json:"type"`
    UserID      string      `json:"user_id"`
    GroupID     *string     `json:"group_id,omitempty"`
    ExpenseID   *string     `json:"expense_id,omitempty"`
    Description string      `json:"description"`
    CreatedAt   time.Time   `json:"created_at"`
}
```

### ActivityType
- `EXPENSE_ADDED`: Expense was added
- `SETTLEMENT`: Debt was settled
- `MEMBER_ADDED`: Member added to group
- `MEMBER_REMOVED`: Member removed from group
- `GROUP_CREATED`: Group was created

### Transaction (Settlement)
```go
type Transaction struct {
    ID        string    `json:"id"`
    From      string    `json:"from"`
    To        string    `json:"to"`
    Amount    float64   `json:"amount"`
    GroupID   *string   `json:"group_id,omitempty"`
    CreatedAt time.Time `json:"created_at"`
}
```

---

## APIs

### 1. Add Expense to Group

**Endpoint:** `POST /api/expenses`

**Request Body:**
```json
{
    "title": "Dinner at Restaurant",
    "description": "Group dinner",
    "amount": 300.00,
    "payer": "user_123",
    "created_by": "user_123",
    "group_id": "group_456",
    "expense_type": "EQUAL",
    "split_type": "EQUAL",
    "participants": ["user_123", "user_789", "user_101"],
    "shares": null,
    "percentages": null,
    "exact_amounts": null
}
```

**Response:**
```json
{
    "id": "exp_789",
    "group_id": "group_456",
    "title": "Dinner at Restaurant",
    "amount": 300.00,
    "payer": "user_123",
    "created_by": "user_123",
    "expense_type": "EQUAL",
    "split_type": "EQUAL",
    "shares": {
        "user_123": 100.00,
        "user_789": 100.00,
        "user_101": 100.00
    },
    "created_at": "2024-01-15T10:30:00Z"
}
```

**Features:**
- Supports EQUAL, EXACT, PERCENTAGE, SHARES expense types
- Validates that shares/percentages/exact amounts sum to total
- Updates balances automatically
- Records activity

### 2. View Balances

**Endpoint:** `GET /api/users/{user_id}/balances`

**Query Parameters:**
- `group_id` (optional): Filter balances for a specific group
- `global` (optional, default: true): Include global balances

**Response (Global):**
```json
{
    "user_id": "user_123",
    "global_balances": [
        {
            "owed_to": "user_789",
            "amount": 150.00,
            "group_id": null
        },
        {
            "owed_to": "user_101",
            "amount": 50.00,
            "group_id": "group_456"
        }
    ],
    "total_owed": 200.00,
    "total_owed_to_me": 75.00,
    "net_balance": -125.00
}
```

**Response (Group-specific):**
```json
{
    "user_id": "user_123",
    "group_id": "group_456",
    "group_balances": [
        {
            "owed_to": "user_101",
            "amount": 50.00
        }
    ],
    "total_owed": 50.00,
    "total_owed_to_me": 0.00,
    "net_balance": -50.00
}
```

### 3. Simplify Debts

**Endpoint:** `POST /api/groups/{group_id}/simplify-debts`

**Response:**
```json
{
    "simplified_transactions": [
        {
            "from": "user_123",
            "to": "user_789",
            "amount": 150.00
        },
        {
            "from": "user_101",
            "to": "user_789",
            "amount": 50.00
        }
    ],
    "original_transaction_count": 5,
    "simplified_transaction_count": 2
}
```

### 4. Get Activity Feed

**Endpoint:** `GET /api/groups/{group_id}/activity`

**Query Parameters:**
- `limit` (optional, default: 50): Number of activities to return
- `offset` (optional, default: 0): Pagination offset

**Response:**
```json
{
    "activities": [
        {
            "id": "act_001",
            "type": "EXPENSE_ADDED",
            "user_id": "user_123",
            "group_id": "group_456",
            "expense_id": "exp_789",
            "description": "user_123 added expense 'Dinner at Restaurant' of $300.00",
            "created_at": "2024-01-15T10:30:00Z"
        }
    ],
    "total": 25,
    "limit": 50,
    "offset": 0
}
```

### 5. Set Default Split Strategy

**Endpoint:** `PUT /api/groups/{group_id}/default-split`

**Request Body:**
```json
{
    "split_type": "PERCENTAGE",
    "split_ratio": {
        "user_123": 0.6,
        "user_789": 0.4
    }
}
```

---

## Follow-up Features

### 1. Simplify Debts Feature
- **Algorithm**: Uses greedy approach to minimize transactions
- **Time Complexity**: O(N log N) where N is number of users
- **Space Complexity**: O(N)
- **Implementation**: 
  - Calculate net balances for each user
  - Sort creditors and debtors
  - Match largest creditor with largest debtor
  - Continue until all debts are settled

### 2. Custom Split Default Strategy
- Groups can have default split strategies
- When adding expense, if split details not provided, use group default
- Supports:
  - Equal split (default)
  - Percentage-based (e.g., 60:40)
  - Share-based (e.g., 2:1:1)
  - Exact amounts

### 3. Activity Tracking
- Tracks all major events:
  - Expense additions
  - Settlements
  - Group membership changes
  - Group creation
- Provides activity feed per group
- Supports pagination

---

## Design Decisions

### 1. Balance Storage
- **Choice**: Store balances as `map[user_id]map[owed_to_user_id]amount`
- **Rationale**: 
  - Fast lookups for "who owes whom"
  - Easy to update when expenses are added
  - Supports both global and group-level balances

### 2. Expense Types vs Split Types
- **ExpenseType**: How the expense is categorized (EQUAL, EXACT, etc.)
- **SplitType**: How the amount is split among participants
- **Rationale**: Separation of concerns allows flexibility

### 3. Activity Tracking
- **Choice**: Store activities in-memory with timestamps
- **Rationale**: 
  - Simple implementation for LLD
  - In production, would use database with indexing
  - Supports pagination and filtering

### 4. Debt Simplification
- **Choice**: Greedy algorithm over optimal (DFS-based)
- **Rationale**:
  - Greedy: O(N log N) - scales better
  - Optimal: O(N!) - only feasible for small groups
  - Greedy produces near-optimal results in practice

### 5. Default Split Strategy
- **Choice**: Store at group level
- **Rationale**:
  - Groups often have consistent split patterns
  - Reduces user input when adding expenses
  - Can be overridden per expense

---

## Extensibility

The design is extensible to support:

1. **Recurring Expenses**: Add `recurrence` field to Expense
2. **Expense Categories**: Add `category` field to Expense
3. **Currency Support**: Add `currency` field to Expense and Balance
4. **Notifications**: Hook into Activity system
5. **Multi-currency**: Extend Balance model
6. **Expense Approval**: Add `status` field to Expense
7. **Partial Settlements**: Extend Transaction model
8. **Expense Splitting Rules**: Add rules engine for complex splits
9. **Analytics**: Query Activity and Expense data
10. **Export/Import**: Serialize/deserialize data models

---

## Running the Application

```bash
cd go/splitwise
go run main.go
```

The server will start on `http://localhost:8080`

### Example API Calls

```bash
# Create user
curl -X GET "http://localhost:8080/api/users/create?name=Alice"

# Create group
curl -X POST "http://localhost:8080/api/groups/create" \
  -H "Content-Type: application/json" \
  -d '{"name": "Trip", "members": ["user_1", "user_2"]}'

# Add expense
curl -X POST "http://localhost:8080/api/expenses" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Dinner",
    "amount": 300.00,
    "payer": "user_1",
    "created_by": "user_1",
    "group_id": "group_1",
    "expense_type": "EQUAL",
    "split_type": "EQUAL",
    "participants": ["user_1", "user_2"]
  }'

# View balances
curl -X GET "http://localhost:8080/api/users/user_1/balances?group_id=group_1"

# Simplify debts
curl -X POST "http://localhost:8080/api/groups/group_1/simplify-debts"

# Get activity
curl -X GET "http://localhost:8080/api/groups/group_1/activity"
```

