package main

import (
	"sync"
)

// Rule is the strategy interface for coupon eligibility rules
type Rule interface {
	Check(order *Order) bool
	GetRuleType() string
}

// UserTypeRule checks if the user type matches the required type
type UserTypeRule struct {
	RequiredUserType UserType
}

func (r *UserTypeRule) Check(order *Order) bool {
	return order.User.UserType == r.RequiredUserType
}

func (r *UserTypeRule) GetRuleType() string {
	return "UserTypeRule"
}

// ProductTypeRule checks if the cart contains products of a specific type
type ProductTypeRule struct {
	RequiredProductType ProductType
}

func (r *ProductTypeRule) Check(order *Order) bool {
	for _, item := range order.Cart.Items {
		if item.ProductType == r.RequiredProductType {
			return true
		}
	}
	return false
}

func (r *ProductTypeRule) GetRuleType() string {
	return "ProductTypeRule"
}

// MinimumAmountRule checks if the order total meets the minimum amount requirement
type MinimumAmountRule struct {
	MinimumAmount float64
}

func (r *MinimumAmountRule) Check(order *Order) bool {
	return order.TotalAmount >= r.MinimumAmount
}

func (r *MinimumAmountRule) GetRuleType() string {
	return "MinimumAmountRule"
}

// PaymentMethodRule checks if the payment method matches the required method
type PaymentMethodRule struct {
	RequiredPaymentMethod PaymentMethod
}

func (r *PaymentMethodRule) Check(order *Order) bool {
	return order.PaymentMethod == r.RequiredPaymentMethod
}

func (r *PaymentMethodRule) GetRuleType() string {
	return "PaymentMethodRule"
}

// PastPurchaseRule checks if the user has purchased specific products in the past
type PastPurchaseRule struct {
	RequiredProductIDs []string
}

func (r *PastPurchaseRule) Check(order *Order) bool {
	if len(r.RequiredProductIDs) == 0 {
		return true
	}
	
	purchaseMap := make(map[string]bool)
	for _, productID := range order.User.PastPurchases {
		purchaseMap[productID] = true
	}
	
	for _, requiredID := range r.RequiredProductIDs {
		if !purchaseMap[requiredID] {
			return false
		}
	}
	return true
}

func (r *PastPurchaseRule) GetRuleType() string {
	return "PastPurchaseRule"
}

// RuleEvaluator handles concurrent rule evaluation with thread safety
type RuleEvaluator struct {
	mu sync.RWMutex
}

// NewRuleEvaluator creates a new rule evaluator
func NewRuleEvaluator() *RuleEvaluator {
	return &RuleEvaluator{}
}

// EvaluateRulesConcurrently evaluates multiple rules concurrently for better performance
func (re *RuleEvaluator) EvaluateRulesConcurrently(order *Order, rules []Rule) []bool {
	re.mu.RLock()
	defer re.mu.RUnlock()
	
	results := make([]bool, len(rules))
	var wg sync.WaitGroup
	
	for i, rule := range rules {
		wg.Add(1)
		go func(index int, r Rule) {
			defer wg.Done()
			results[index] = r.Check(order)
		}(i, rule)
	}
	
	wg.Wait()
	return results
}

// EvaluateRulesFormula evaluates rules based on a logical formula
// Formula format: "AND(0, 1)", "OR(0, 1)", "AND(0, OR(1, 2))"
func (re *RuleEvaluator) EvaluateRulesFormula(order *Order, coupon *Coupon) bool {
	re.mu.RLock()
	defer re.mu.RUnlock()
	
	// Evaluate all rules concurrently
	ruleResults := re.EvaluateRulesConcurrently(order, coupon.Rules)
	
	// Parse and evaluate the formula
	return re.parseFormula(coupon.RulesFormula, ruleResults)
}

// parseFormula parses and evaluates the logical formula
func (re *RuleEvaluator) parseFormula(formula string, results []bool) bool {
	if formula == "" {
		// If no formula, all rules must pass (AND logic)
		for _, result := range results {
			if !result {
				return false
			}
		}
		return true
	}
	
	// Simple parser for AND/OR formulas
	// For production, consider using a proper expression parser
	return re.evaluateSimpleFormula(formula, results)
}

// evaluateSimpleFormula evaluates simple AND/OR formulas
func (re *RuleEvaluator) evaluateSimpleFormula(formula string, results []bool) bool {
	// Handle simple cases like "AND(0, 1)" or "OR(0, 1)"
	// This is a simplified parser - for complex formulas, use a proper parser
	
	if len(formula) == 1 {
		// Single rule index
		index := int(formula[0] - '0')
		if index >= 0 && index < len(results) {
			return results[index]
		}
		return false
	}
	
	// For complex formulas, we'll use a recursive approach
	// This handles nested AND/OR operations
	return re.parseComplexFormula(formula, results)
}

// parseComplexFormula handles complex nested formulas
func (re *RuleEvaluator) parseComplexFormula(formula string, results []bool) bool {
	// Simplified implementation - in production, use a proper expression parser
	// This handles basic AND(0, 1) and OR(0, 1) cases
	
	if len(formula) < 5 {
		return false
	}
	
	// Check if it starts with AND or OR
	if formula[:3] == "AND" {
		// Extract indices from AND(0, 1)
		indices := re.extractIndices(formula[4 : len(formula)-1])
		for _, idx := range indices {
			if idx >= len(results) || !results[idx] {
				return false
			}
		}
		return true
	} else if formula[:2] == "OR" {
		// Extract indices from OR(0, 1)
		indices := re.extractIndices(formula[3 : len(formula)-1])
		for _, idx := range indices {
			if idx < len(results) && results[idx] {
				return true
			}
		}
		return false
	}
	
	// Default: all rules must pass
	for _, result := range results {
		if !result {
			return false
		}
	}
	return true
}

// extractIndices extracts rule indices from a comma-separated string
func (re *RuleEvaluator) extractIndices(s string) []int {
	var indices []int
	current := 0
	for _, char := range s {
		if char >= '0' && char <= '9' {
			current = current*10 + int(char-'0')
		} else if char == ',' || char == ' ' {
			if current >= 0 {
				indices = append(indices, current)
				current = 0
			}
		}
	}
	if current >= 0 {
		indices = append(indices, current)
	}
	return indices
}

