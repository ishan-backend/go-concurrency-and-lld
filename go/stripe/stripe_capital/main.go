package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Transaction represents a financial transaction
type Transaction struct {
	Amount      float64
	Date        time.Time
	Description string
}

// StripeCapital manages capital financing for businesses
type StripeCapital struct {
	transactions []Transaction
}

// NewStripeCapital creates a new Stripe Capital instance
func NewStripeCapital() *StripeCapital {
	return &StripeCapital{
		transactions: make([]Transaction, 0),
	}
}

// AddTransaction adds a transaction to the history
func (sc *StripeCapital) AddTransaction(amount float64, date time.Time, description string) {
	sc.transactions = append(sc.transactions, Transaction{
		Amount:      amount,
		Date:        date,
		Description: description,
	})
}

// CalculateTotalCapital calculates total capital from all transactions
func (sc *StripeCapital) CalculateTotalCapital() float64 {
	total := 0.0
	for _, txn := range sc.transactions {
		total += txn.Amount
	}
	return total
}

// CalculateAverageMonthlyRevenue calculates average monthly revenue
func (sc *StripeCapital) CalculateAverageMonthlyRevenue() float64 {
	if len(sc.transactions) == 0 {
		return 0.0
	}
	
	// Group transactions by month
	monthlyRevenue := make(map[string]float64)
	for _, txn := range sc.transactions {
		if txn.Amount > 0 { // Only count positive amounts as revenue
			monthKey := txn.Date.Format("2006-01")
			monthlyRevenue[monthKey] += txn.Amount
		}
	}
	
	if len(monthlyRevenue) == 0 {
		return 0.0
	}
	
	total := 0.0
	for _, revenue := range monthlyRevenue {
		total += revenue
	}
	
	return total / float64(len(monthlyRevenue))
}

// IsEligibleForLoan checks if business is eligible for a loan
// Eligibility criteria: Average monthly revenue >= threshold
func (sc *StripeCapital) IsEligibleForLoan(threshold float64) bool {
	avgRevenue := sc.CalculateAverageMonthlyRevenue()
	return avgRevenue >= threshold
}

// GetLoanAmount calculates potential loan amount based on revenue
// Typically, loan amount is a multiple of average monthly revenue
func (sc *StripeCapital) GetLoanAmount(multiplier float64) float64 {
	avgRevenue := sc.CalculateAverageMonthlyRevenue()
	return avgRevenue * multiplier
}

// GetMonthlyRevenueBreakdown returns revenue breakdown by month
func (sc *StripeCapital) GetMonthlyRevenueBreakdown() map[string]float64 {
	monthlyRevenue := make(map[string]float64)
	for _, txn := range sc.transactions {
		if txn.Amount > 0 {
			monthKey := txn.Date.Format("2006-01")
			monthlyRevenue[monthKey] += txn.Amount
		}
	}
	return monthlyRevenue
}

// GetTopMonths returns top N months by revenue
func (sc *StripeCapital) GetTopMonths(n int) []string {
	breakdown := sc.GetMonthlyRevenueBreakdown()
	
	type monthRevenue struct {
		month   string
		revenue float64
	}
	
	var months []monthRevenue
	for month, revenue := range breakdown {
		months = append(months, monthRevenue{month: month, revenue: revenue})
	}
	
	// Sort by revenue (descending)
	sort.Slice(months, func(i, j int) bool {
		if months[i].revenue != months[j].revenue {
			return months[i].revenue > months[j].revenue
		}
		return months[i].month > months[j].month
	})
	
	result := make([]string, 0)
	for i := 0; i < len(months) && i < n; i++ {
		result = append(result, months[i].month)
	}
	
	return result
}

func main() {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Stripe Capital - Capital Financing System")
	fmt.Println(strings.Repeat("=", 60))
	
	capital := NewStripeCapital()
	
	// Example 1: Basic eligibility check
	fmt.Println("\nExample 1: Basic Eligibility Check")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	baseDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	capital.AddTransaction(12000, baseDate, "Payment 1")
	capital.AddTransaction(15000, baseDate.AddDate(0, 1, 0), "Payment 2")
	capital.AddTransaction(11000, baseDate.AddDate(0, 2, 0), "Payment 3")
	capital.AddTransaction(13000, baseDate.AddDate(0, 3, 0), "Payment 4")
	
	totalCapital := capital.CalculateTotalCapital()
	avgRevenue := capital.CalculateAverageMonthlyRevenue()
	eligible := capital.IsEligibleForLoan(10000)
	
	fmt.Printf("Total Capital: $%.2f\n", totalCapital)
	fmt.Printf("Average Monthly Revenue: $%.2f\n", avgRevenue)
	fmt.Printf("Eligible for Loan (threshold $10,000): %v\n", eligible)
	
	// Example 2: Loan amount calculation
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 2: Loan Amount Calculation")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	loanAmount := capital.GetLoanAmount(3.0) // 3x monthly revenue
	fmt.Printf("Potential Loan Amount (3x monthly revenue): $%.2f\n", loanAmount)
	
	// Example 3: Monthly breakdown
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 3: Monthly Revenue Breakdown")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	breakdown := capital.GetMonthlyRevenueBreakdown()
	fmt.Println("Monthly Revenue Breakdown:")
	for month, revenue := range breakdown {
		fmt.Printf("  %s: $%.2f\n", month, revenue)
	}
	
	// Example 4: Top months
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 4: Top Months by Revenue")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	topMonths := capital.GetTopMonths(3)
	fmt.Println("Top 3 months:")
	for i, month := range topMonths {
		revenue := breakdown[month]
		fmt.Printf("  %d. %s: $%.2f\n", i+1, month, revenue)
	}
	
	// Example 5: Not eligible case
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 5: Low Revenue Business")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	capital2 := NewStripeCapital()
	capital2.AddTransaction(5000, baseDate, "Payment 1")
	capital2.AddTransaction(6000, baseDate.AddDate(0, 1, 0), "Payment 2")
	capital2.AddTransaction(4000, baseDate.AddDate(0, 2, 0), "Payment 3")
	
	avgRevenue2 := capital2.CalculateAverageMonthlyRevenue()
	eligible2 := capital2.IsEligibleForLoan(10000)
	
	fmt.Printf("Average Monthly Revenue: $%.2f\n", avgRevenue2)
	fmt.Printf("Eligible for Loan (threshold $10,000): %v\n", eligible2)
}

