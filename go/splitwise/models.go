package main

import "time"

// ------------------------- Enums -------------------------

type ExpenseType string

const (
	ExpenseTypeEqual      ExpenseType = "EQUAL"
	ExpenseTypeExact      ExpenseType = "EXACT"
	ExpenseTypePercentage ExpenseType = "PERCENTAGE"
	ExpenseTypeShares     ExpenseType = "SHARES"
)

type SplitType string

const (
	SplitTypeEqual      SplitType = "EQUAL"
	SplitTypeUnequal    SplitType = "UNEQUAL"
	SplitTypePercentage SplitType = "PERCENTAGE"
	SplitTypeShares     SplitType = "SHARES"
)

type ActivityType string

const (
	ActivityTypeExpenseAdded  ActivityType = "EXPENSE_ADDED"
	ActivityTypeSettlement    ActivityType = "SETTLEMENT"
	ActivityTypeMemberAdded   ActivityType = "MEMBER_ADDED"
	ActivityTypeMemberRemoved ActivityType = "MEMBER_REMOVED"
	ActivityTypeGroupCreated  ActivityType = "GROUP_CREATED"
)

// ------------------------- Models -------------------------

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type Group struct {
	ID                string             `json:"id"`
	Name              string             `json:"name"`
	Members           []string           `json:"members"`
	Admin             string             `json:"admin"`
	DefaultSplitType  SplitType          `json:"default_split_type"`
	DefaultSplitRatio map[string]float64 `json:"default_split_ratio"` // e.g., {"user1": 0.6, "user2": 0.4}
	CreatedAt         time.Time          `json:"created_at"`
}

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

type Balance struct {
	UserID  string  `json:"user_id"`
	OwedTo  string  `json:"owed_to"`
	Amount  float64 `json:"amount"` // positive = owes, negative = is owed
	GroupID *string `json:"group_id,omitempty"`
}

type Activity struct {
	ID          string       `json:"id"`
	Type        ActivityType `json:"type"`
	UserID      string       `json:"user_id"`
	GroupID     *string      `json:"group_id,omitempty"`
	ExpenseID   *string      `json:"expense_id,omitempty"`
	Description string       `json:"description"`
	CreatedAt   time.Time    `json:"created_at"`
}

type Transaction struct {
	ID        string    `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Amount    float64   `json:"amount"`
	GroupID   *string   `json:"group_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ------------------------- Request/Response Models -------------------------

type AddExpenseRequest struct {
	Title        string             `json:"title"`
	Description  string             `json:"description,omitempty"`
	Amount       float64            `json:"amount"`
	Payer        string             `json:"payer"`
	CreatedBy    string             `json:"created_by"`
	GroupID      *string            `json:"group_id,omitempty"`
	ExpenseType  ExpenseType        `json:"expense_type"`
	SplitType    SplitType          `json:"split_type"`
	Participants []string           `json:"participants"`
	Shares       map[string]float64 `json:"shares,omitempty"`        // For SHARES split type
	Percentages  map[string]float64 `json:"percentages,omitempty"`   // For PERCENTAGE split type
	ExactAmounts map[string]float64 `json:"exact_amounts,omitempty"` // For EXACT split type
}

type BalanceResponse struct {
	UserID         string    `json:"user_id"`
	GroupID        *string   `json:"group_id,omitempty"`
	GlobalBalances []Balance `json:"global_balances,omitempty"`
	GroupBalances  []Balance `json:"group_balances,omitempty"`
	TotalOwed      float64   `json:"total_owed"`
	TotalOwedToMe  float64   `json:"total_owed_to_me"`
	NetBalance     float64   `json:"net_balance"`
}

type SimplifyDebtsResponse struct {
	SimplifiedTransactions     []Transaction `json:"simplified_transactions"`
	OriginalTransactionCount   int           `json:"original_transaction_count"`
	SimplifiedTransactionCount int           `json:"simplified_transaction_count"`
}

type ActivityResponse struct {
	Activities []Activity `json:"activities"`
	Total      int        `json:"total"`
	Limit      int        `json:"limit"`
	Offset     int        `json:"offset"`
}

type SetDefaultSplitRequest struct {
	SplitType  SplitType          `json:"split_type"`
	SplitRatio map[string]float64 `json:"split_ratio"`
}
