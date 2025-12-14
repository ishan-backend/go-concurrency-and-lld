package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypeEarn   TransactionType = "earn"
	TransactionTypeSpend  TransactionType = "spend"
	TransactionTypeExpire TransactionType = "expire"
)

// Transaction represents a points transaction
type Transaction struct {
	UserID      string
	Type        TransactionType
	Amount      int
	Timestamp   time.Time
	Description string
}

// UserPointsManager manages user points with expiration support
type UserPointsManager struct {
	transactions []Transaction
	// userID -> list of point balances with expiration times
	// Each entry: {amount, expirationTime}
	userBalances map[string][]struct {
		amount        int
		expirationTime time.Time
	}
}

// NewUserPointsManager creates a new points manager
func NewUserPointsManager() *UserPointsManager {
	return &UserPointsManager{
		transactions:  make([]Transaction, 0),
		userBalances: make(map[string][]struct {
			amount        int
			expirationTime time.Time
		}),
	}
}

// AddTransaction adds a transaction and updates balances
func (u *UserPointsManager) AddTransaction(userID string, txnType TransactionType, amount int, description string, expirationTime *time.Time) {
	txn := Transaction{
		UserID:      userID,
		Type:        txnType,
		Amount:      amount,
		Timestamp:   time.Now(),
		Description: description,
	}
	u.transactions = append(u.transactions, txn)
	
	if u.userBalances[userID] == nil {
		u.userBalances[userID] = make([]struct {
			amount        int
			expirationTime time.Time
		}, 0)
	}
	
	switch txnType {
	case TransactionTypeEarn:
		expTime := time.Now().Add(365 * 24 * time.Hour) // Default 1 year
		if expirationTime != nil {
			expTime = *expirationTime
		}
		u.userBalances[userID] = append(u.userBalances[userID], struct {
			amount        int
			expirationTime time.Time
		}{amount: amount, expirationTime: expTime})
		
	case TransactionTypeSpend:
		u.spendPoints(userID, amount)
		
	case TransactionTypeExpire:
		u.expirePoints(userID, amount)
	}
}

// spendPoints spends points using FIFO (oldest points first)
func (u *UserPointsManager) spendPoints(userID string, amount int) {
	balance := u.userBalances[userID]
	if balance == nil {
		return
	}
	
	remaining := amount
	for i := 0; i < len(balance) && remaining > 0; {
		if balance[i].amount <= remaining {
			remaining -= balance[i].amount
			balance = append(balance[:i], balance[i+1:]...)
		} else {
			balance[i].amount -= remaining
			remaining = 0
		}
	}
	
	u.userBalances[userID] = balance
}

// expirePoints removes expired points
func (u *UserPointsManager) expirePoints(userID string, amount int) {
	balance := u.userBalances[userID]
	if balance == nil {
		return
	}
	
	now := time.Now()
	remaining := amount
	for i := 0; i < len(balance) && remaining > 0; {
		if balance[i].expirationTime.Before(now) || balance[i].expirationTime.Equal(now) {
			if balance[i].amount <= remaining {
				remaining -= balance[i].amount
				balance = append(balance[:i], balance[i+1:]...)
			} else {
				balance[i].amount -= remaining
				remaining = 0
			}
		} else {
			i++
		}
	}
	
	u.userBalances[userID] = balance
}

// GetBalance gets the current balance for a user (excluding expired points)
func (u *UserPointsManager) GetBalance(userID string) int {
	u.cleanExpiredPoints(userID)
	
	balance := u.userBalances[userID]
	if balance == nil {
		return 0
	}
	
	total := 0
	for _, entry := range balance {
		total += entry.amount
	}
	return total
}

// cleanExpiredPoints removes expired points for a user
func (u *UserPointsManager) cleanExpiredPoints(userID string) {
	balance := u.userBalances[userID]
	if balance == nil {
		return
	}
	
	now := time.Now()
	filtered := make([]struct {
		amount        int
		expirationTime time.Time
	}, 0)
	
	for _, entry := range balance {
		if entry.expirationTime.After(now) {
			filtered = append(filtered, entry)
		}
	}
	
	u.userBalances[userID] = filtered
}

// GetTransactions returns all transactions for a user
func (u *UserPointsManager) GetTransactions(userID string) []Transaction {
	var result []Transaction
	for _, txn := range u.transactions {
		if txn.UserID == userID {
			result = append(result, txn)
		}
	}
	return result
}

// GetAllBalances returns balances for all users
func (u *UserPointsManager) GetAllBalances() map[string]int {
	result := make(map[string]int)
	for userID := range u.userBalances {
		result[userID] = u.GetBalance(userID)
	}
	return result
}

// GetTopUsers returns users sorted by balance (highest first)
func (u *UserPointsManager) GetTopUsers(limit int) []string {
	balances := u.GetAllBalances()
	
	type userBalance struct {
		userID  string
		balance int
	}
	
	var userBalances []userBalance
	for userID, balance := range balances {
		userBalances = append(userBalances, userBalance{
			userID:  userID,
			balance: balance,
		})
	}
	
	// Sort by balance (highest first), then by userID
	sort.Slice(userBalances, func(i, j int) bool {
		if userBalances[i].balance != userBalances[j].balance {
			return userBalances[i].balance > userBalances[j].balance
		}
		return userBalances[i].userID < userBalances[j].userID
	})
	
	result := make([]string, 0)
	for i, ub := range userBalances {
		if i >= limit {
			break
		}
		result = append(result, ub.userID)
	}
	
	return result
}

func main() {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("User Points Manager")
	fmt.Println(strings.Repeat("=", 60))
	
	manager := NewUserPointsManager()
	
	// Example 1: Basic earn and spend
	fmt.Println("\nExample 1: Basic Earn and Spend")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	manager.AddTransaction("user1", TransactionTypeEarn, 100, "Sign up bonus", nil)
	manager.AddTransaction("user1", TransactionTypeEarn, 50, "Referral bonus", nil)
	manager.AddTransaction("user1", TransactionTypeSpend, 75, "Purchase item", nil)
	
	fmt.Printf("User1 balance: %d points\n", manager.GetBalance("user1"))
	// Expected: 100 + 50 - 75 = 75
	
	// Example 2: Points expiration
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 2: Points Expiration")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	manager2 := NewUserPointsManager()
	now := time.Now()
	
	// Add points with different expiration times
	exp1 := now.Add(30 * 24 * time.Hour) // 30 days
	exp2 := now.Add(60 * 24 * time.Hour) // 60 days
	
	manager2.AddTransaction("user2", TransactionTypeEarn, 100, "Bonus 1", &exp1)
	manager2.AddTransaction("user2", TransactionTypeEarn, 200, "Bonus 2", &exp2)
	
	fmt.Printf("User2 balance: %d points\n", manager2.GetBalance("user2"))
	// Expected: 100 + 200 = 300
	
	// Simulate expiration (30 days later)
	expiredTime := now.Add(31 * 24 * time.Hour)
	// In real scenario, we'd check expiration when getting balance
	// For demo, we manually expire
	manager2.AddTransaction("user2", TransactionTypeExpire, 100, "Expired points", &expiredTime)
	
	fmt.Printf("User2 balance after expiration: %d points\n", manager2.GetBalance("user2"))
	// Expected: 200 (100 expired)
	
	// Example 3: FIFO spending
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 3: FIFO Spending (Oldest Points First)")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	manager3 := NewUserPointsManager()
	
	manager3.AddTransaction("user3", TransactionTypeEarn, 100, "Earn 1", nil)
	manager3.AddTransaction("user3", TransactionTypeEarn, 50, "Earn 2", nil)
	manager3.AddTransaction("user3", TransactionTypeEarn, 25, "Earn 3", nil)
	
	fmt.Printf("User3 balance before spend: %d points\n", manager3.GetBalance("user3"))
	
	manager3.AddTransaction("user3", TransactionTypeSpend, 75, "Spend 75", nil)
	
	fmt.Printf("User3 balance after spending 75: %d points\n", manager3.GetBalance("user3"))
	// Expected: 100 (first 100 spent, then 50 partially spent, leaving 25+25=50)
	// Actually: 100 + 50 + 25 - 75 = 100
	
	// Example 4: Top users
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 4: Top Users by Balance")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	manager4 := NewUserPointsManager()
	manager4.AddTransaction("alice", TransactionTypeEarn, 500, "Bonus", nil)
	manager4.AddTransaction("bob", TransactionTypeEarn, 300, "Bonus", nil)
	manager4.AddTransaction("bob", TransactionTypeEarn, 200, "Referral", nil)
	manager4.AddTransaction("charlie", TransactionTypeEarn, 100, "Bonus", nil)
	manager4.AddTransaction("charlie", TransactionTypeSpend, 50, "Purchase", nil)
	
	topUsers := manager4.GetTopUsers(3)
	fmt.Println("Top 3 users by balance:")
	for i, userID := range topUsers {
		balance := manager4.GetBalance(userID)
		fmt.Printf("  %d. %s: %d points\n", i+1, userID, balance)
	}
	// Expected:
	// 1. bob: 500 (300 + 200)
	// 2. alice: 500
	// 3. charlie: 50 (100 - 50)
	
	// Example 5: Transaction history
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 5: Transaction History")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	transactions := manager4.GetTransactions("bob")
	fmt.Println("Transactions for bob:")
	for i, txn := range transactions {
		fmt.Printf("  %d. %s: %d points - %s\n", i+1, txn.Type, txn.Amount, txn.Description)
	}
}

