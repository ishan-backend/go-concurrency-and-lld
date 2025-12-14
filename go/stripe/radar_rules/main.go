package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Charge represents a payment charge
type Charge struct {
	CardCountry string
	Currency    string
	Amount      int
	IPCountry   string
}

// ParseCharge parses a charge string in format: "CHARGE: card_country=US&currency=USD&amount=150&ip_country=CA"
func ParseCharge(chargeStr string) *Charge {
	charge := &Charge{}
	
	// Remove "CHARGE: " prefix
	chargeStr = strings.TrimPrefix(chargeStr, "CHARGE: ")
	
	// Split by &
	parts := strings.Split(chargeStr, "&")
	for _, part := range parts {
		kv := strings.Split(part, "=")
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			
			switch key {
			case "card_country":
				charge.CardCountry = value
			case "currency":
				charge.Currency = value
			case "amount":
				if amount, err := strconv.Atoi(value); err == nil {
					charge.Amount = amount
				}
			case "ip_country":
				charge.IPCountry = value
			}
		}
	}
	
	return charge
}

// Rule represents a rule (ALLOW or BLOCK)
type Rule struct {
	Type      string // "ALLOW" or "BLOCK"
	Condition string // e.g., "amount<100", "card_country != ip_country AND amount > 100"
}

// ParseRule parses a rule string
func ParseRule(ruleStr string) *Rule {
	rule := &Rule{}
	
	if strings.HasPrefix(ruleStr, "ALLOW:") {
		rule.Type = "ALLOW"
		rule.Condition = strings.TrimPrefix(ruleStr, "ALLOW:")
	} else if strings.HasPrefix(ruleStr, "BLOCK:") {
		rule.Type = "BLOCK"
		rule.Condition = strings.TrimPrefix(ruleStr, "BLOCK:")
	}
	
	return rule
}

// EvaluateCondition evaluates a condition against a charge
// Returns true if condition is satisfied
func EvaluateCondition(condition string, charge *Charge) bool {
	condition = strings.TrimSpace(condition)
	
	// Handle AND conditions
	if strings.Contains(condition, " AND ") {
		parts := strings.Split(condition, " AND ")
		for _, part := range parts {
			if !EvaluateCondition(strings.TrimSpace(part), charge) {
				return false
			}
		}
		return true
	}
	
	// Handle OR conditions
	if strings.Contains(condition, " OR ") {
		parts := strings.Split(condition, " OR ")
		for _, part := range parts {
			if EvaluateCondition(strings.TrimSpace(part), charge) {
				return true
			}
		}
		return false
	}
	
	// Handle comparison operators
	operators := []string{"!=", "==", ">=", "<=", ">", "<"}
	for _, op := range operators {
		if strings.Contains(condition, op) {
			return evaluateComparison(condition, op, charge)
		}
	}
	
	return false
}

// evaluateComparison evaluates a comparison condition
func evaluateComparison(condition, operator string, charge *Charge) bool {
	parts := strings.Split(condition, operator)
	if len(parts) != 2 {
		return false
	}
	
	left := strings.TrimSpace(parts[0])
	right := strings.TrimSpace(parts[1])
	
	// Handle amount comparisons
	if left == "amount" {
		rightVal, err := strconv.Atoi(right)
		if err != nil {
			return false
		}
		return compareInt(charge.Amount, operator, rightVal)
	}
	
	// Handle country comparisons
	if left == "card_country" && right == "ip_country" {
		return compareString(charge.CardCountry, operator, charge.IPCountry)
	}
	
	// Handle direct value comparisons
	if strings.HasPrefix(left, "card_country=") || strings.HasPrefix(left, "ip_country=") {
		// Extract value from left side
		leftParts := strings.Split(left, "=")
		if len(leftParts) == 2 {
			field := leftParts[0]
			value := leftParts[1]
			
			var chargeValue string
			if field == "card_country" {
				chargeValue = charge.CardCountry
			} else if field == "ip_country" {
				chargeValue = charge.IPCountry
			}
			
			return compareString(chargeValue, operator, value)
		}
	}
	
	return false
}

// compareInt compares two integers with the given operator
func compareInt(a int, operator string, b int) bool {
	switch operator {
	case ">":
		return a > b
	case "<":
		return a < b
	case ">=":
		return a >= b
	case "<=":
		return a <= b
	case "==":
		return a == b
	case "!=":
		return a != b
	}
	return false
}

// compareString compares two strings with the given operator
func compareString(a, operator, b string) bool {
	switch operator {
	case "==":
		return a == b
	case "!=":
		return a != b
	}
	return false
}

// ProcessRadarRules processes charge and rules, returns 1 if allowed, 0 if blocked
func ProcessRadarRules(input string) int {
	parts := strings.Split(input, ",")
	if len(parts) == 0 {
		return 1 // Default allow
	}
	
	// Parse charge (first part)
	charge := ParseCharge(parts[0])
	
	// Process rules in order
	for i := 1; i < len(parts); i++ {
		ruleStr := strings.TrimSpace(parts[i])
		if ruleStr == "" {
			continue
		}
		
		rule := ParseRule(ruleStr)
		if rule == nil {
			continue
		}
		
		conditionMet := EvaluateCondition(rule.Condition, charge)
		
		if rule.Type == "ALLOW" {
			if conditionMet {
				return 1 // Allowed
			}
		} else if rule.Type == "BLOCK" {
			if conditionMet {
				return 0 // Blocked
			}
		}
	}
	
	// Default: if no rule matched, allow
	return 1
}

func main() {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Stripe Radar Rules Engine")
	fmt.Println(strings.Repeat("=", 60))
	
	// Example 1: Basic block
	fmt.Println("\nExample 1: Basic Block")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	input1 := "CHARGE: card_country=US&currency=USD&amount=150&ip_country=CA,BLOCK:amount > 100"
	result1 := ProcessRadarRules(input1)
	fmt.Printf("Input: %s\n", input1)
	fmt.Printf("Result: %d (0 = blocked, 1 = allowed)\n", result1)
	fmt.Printf("Expected: 0 (blocked because amount 150 > 100)\n")
	
	// Example 2: Compound block with ALLOW and BLOCK
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 2: Compound Block")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	input2 := "CHARGE: card_country=US&currency=USD&amount=150&ip_country=CA,ALLOW:amount<100,BLOCK:card_country != ip_country AND amount > 100"
	result2 := ProcessRadarRules(input2)
	fmt.Printf("Input: %s\n", input2)
	fmt.Printf("Result: %d (0 = blocked, 1 = allowed)\n", result2)
	fmt.Printf("Expected: 0 (blocked: card_country(US) != ip_country(CA) AND amount(150) > 100)\n")
	
	// Example 3: ALLOW rule matches
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 3: ALLOW Rule Matches")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	input3 := "CHARGE: card_country=US&currency=USD&amount=50&ip_country=CA,ALLOW:amount<100,BLOCK:card_country != ip_country AND amount > 100"
	result3 := ProcessRadarRules(input3)
	fmt.Printf("Input: %s\n", input3)
	fmt.Printf("Result: %d (0 = blocked, 1 = allowed)\n", result3)
	fmt.Printf("Expected: 1 (allowed because amount 50 < 100)\n")
	
	// Example 4: BLOCK rule with country mismatch only
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 4: Country Mismatch Only")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	input4 := "CHARGE: card_country=US&currency=USD&amount=50&ip_country=CA,BLOCK:card_country != ip_country"
	result4 := ProcessRadarRules(input4)
	fmt.Printf("Input: %s\n", input4)
	fmt.Printf("Result: %d (0 = blocked, 1 = allowed)\n", result4)
	fmt.Printf("Expected: 0 (blocked: US != CA)\n")
	
	// Example 5: No rules match
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 5: No Rules Match")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	input5 := "CHARGE: card_country=US&currency=USD&amount=50&ip_country=US,ALLOW:amount>100,BLOCK:card_country != ip_country"
	result5 := ProcessRadarRules(input5)
	fmt.Printf("Input: %s\n", input5)
	fmt.Printf("Result: %d (0 = blocked, 1 = allowed)\n", result5)
	fmt.Printf("Expected: 1 (default allow, no rules matched)\n")
	
	// Example 6: Multiple BLOCK rules
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 6: Multiple BLOCK Rules")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	input6 := "CHARGE: card_country=US&currency=USD&amount=200&ip_country=US,BLOCK:amount>150,BLOCK:amount<50"
	result6 := ProcessRadarRules(input6)
	fmt.Printf("Input: %s\n", input6)
	fmt.Printf("Result: %d (0 = blocked, 1 = allowed)\n", result6)
	fmt.Printf("Expected: 0 (blocked by first rule: amount 200 > 150)\n")
	
	// Example 7: Complex AND condition
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 7: Complex AND Condition")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	input7 := "CHARGE: card_country=US&currency=USD&amount=120&ip_country=CA,BLOCK:card_country != ip_country AND amount > 100"
	result7 := ProcessRadarRules(input7)
	fmt.Printf("Input: %s\n", input7)
	fmt.Printf("Result: %d (0 = blocked, 1 = allowed)\n", result7)
	fmt.Printf("Expected: 0 (blocked: US != CA AND 120 > 100)\n")
}

