package main

import (
	"fmt"
	"time"
)

func main() {
	// Initialize services
	validator := NewCouponValidator()
	cartService := NewCartService(validator)

	// Create sample coupons
	setupCoupons(validator)

	// Create a shopping cart
	cart := cartService.CreateCart("USER123")
	fmt.Println("=== Shopping Cart Coupon Application System ===")
	fmt.Println()

	// Add items to cart
	fmt.Println("1. Adding items to cart...")
	cartService.AddItem(cart.CartID, &CartItem{
		ProductID: "P1",
		Name:      "Laptop",
		Price:     50000.0,
		Quantity:  1,
		Category:  "Electronics",
	})

	cartService.AddItem(cart.CartID, &CartItem{
		ProductID: "P2",
		Name:      "Mouse",
		Price:     500.0,
		Quantity:  2,
		Category:  "Electronics",
	})

	cartService.AddItem(cart.CartID, &CartItem{
		ProductID: "P3",
		Name:      "Keyboard",
		Price:     2000.0,
		Quantity:  1,
		Category:  "Electronics",
	})

	displayCart(cartService, cart.CartID)

	// Test 1: Apply percentage coupon
	fmt.Println("\n2. Applying PERCENTAGE coupon (SAVE10)...")
	result, err := cartService.ApplyCoupon(cart.CartID, "SAVE10")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Result: %s - Discount: $%.2f\n", result.Message, result.DiscountAmount)
	}
	displayCart(cartService, cart.CartID)

	// Test 2: Try to apply expired coupon
	fmt.Println("\n3. Trying to apply EXPIRED coupon...")
	result, err = cartService.ApplyCoupon(cart.CartID, "EXPIRED")
	if err != nil {
		fmt.Printf("Expected Error: %v\n", err)
	}

	// Test 3: Apply fixed amount coupon
	fmt.Println("\n4. Applying FIXED_AMOUNT coupon (FLAT500)...")
	result, err = cartService.ApplyCoupon(cart.CartID, "FLAT500")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Result: %s - Discount: $%.2f\n", result.Message, result.DiscountAmount)
	}
	displayCart(cartService, cart.CartID)

	// Test 4: Remove a coupon
	fmt.Println("\n5. Removing SAVE10 coupon...")
	err = cartService.RemoveCoupon(cart.CartID, "SAVE10")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Coupon removed successfully")
	}
	displayCart(cartService, cart.CartID)

	// Test 5: Buy X Get Y coupon
	fmt.Println("\n6. Creating new cart for Buy X Get Y test...")
	cart2 := cartService.CreateCart("USER456")
	cartService.AddItem(cart2.CartID, &CartItem{
		ProductID: "P4",
		Name:      "T-Shirt",
		Price:     500.0,
		Quantity:  5,
		Category:  "Clothing",
	})

	fmt.Println("Applying BUY_X_GET_Y coupon (BUY3GET1)...")
	result, err = cartService.ApplyCoupon(cart2.CartID, "BUY3GET1")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Result: %s - Discount: $%.2f\n", result.Message, result.DiscountAmount)
	}
	displayCart(cartService, cart2.CartID)

	// Test 6: Free shipping coupon
	fmt.Println("\n7. Creating cart for free shipping test...")
	cart3 := cartService.CreateCart("USER789")
	cartService.AddItem(cart3.CartID, &CartItem{
		ProductID: "P5",
		Name:      "Book",
		Price:     300.0,
		Quantity:  1,
		Category:  "Books",
	})

	fmt.Println("Applying FREE_SHIPPING coupon...")
	result, err = cartService.ApplyCoupon(cart3.CartID, "FREESHIP")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Result: %s\n", result.Message)
	}
	displayCart(cartService, cart3.CartID)

	// Test 7: User-specific coupon
	fmt.Println("\n8. Testing user-specific coupon...")
	cart4 := cartService.CreateCart("VIP_USER")
	cartService.AddItem(cart4.CartID, &CartItem{
		ProductID: "P6",
		Name:      "Premium Product",
		Price:     10000.0,
		Quantity:  1,
		Category:  "Premium",
	})

	fmt.Println("Applying VIP coupon (VIP20)...")
	result, err = cartService.ApplyCoupon(cart4.CartID, "VIP20")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Result: %s - Discount: $%.2f\n", result.Message, result.DiscountAmount)
	}
	displayCart(cartService, cart4.CartID)

	// Test 8: Minimum purchase validation
	fmt.Println("\n9. Testing minimum purchase validation...")
	cart5 := cartService.CreateCart("USER999")
	cartService.AddItem(cart5.CartID, &CartItem{
		ProductID: "P7",
		Name:      "Cheap Item",
		Price:     100.0,
		Quantity:  1,
		Category:  "Misc",
	})

	fmt.Println("Trying to apply MIN500 coupon (requires $500 minimum)...")
	result, err = cartService.ApplyCoupon(cart5.CartID, "MIN500")
	if err != nil {
		fmt.Printf("Expected Error: %v\n", err)
	}

	// Test 9: Get all applicable coupons
	fmt.Println("\n10. Getting all applicable coupons for cart...")
	applicableCoupons, err := cartService.GetApplicableCoupons(cart.CartID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Found %d applicable coupon(s):\n", len(applicableCoupons))
		currentCart, _ := cartService.GetCart(cart.CartID)
		for _, coupon := range applicableCoupons {
			// Check if already applied
			alreadyApplied := false
			if currentCart != nil {
				for _, appliedCoupon := range currentCart.AppliedCoupons {
					if appliedCoupon.CouponID == coupon.CouponID {
						alreadyApplied = true
						break
					}
				}
			}
			status := "Available"
			if alreadyApplied {
				status = "Already Applied"
			}
			fmt.Printf("  - %s (%s) - %s\n", coupon.Code, coupon.Type, status)
		}
	}
}

func setupCoupons(validator *CouponValidator) {
	now := time.Now()

	// Percentage coupon
	validator.RegisterCoupon(&Coupon{
		CouponID:      "C1",
		Code:          "SAVE10",
		Type:          CouponTypePercentage,
		DiscountValue: 10.0,
		MinPurchase:   1000.0,
		MaxDiscount:   2000.0,
		ValidFrom:     now.AddDate(0, -1, 0),
		ValidUntil:    now.AddDate(0, 1, 0),
		UsageLimit:    100,
		UsedCount:     0,
	})

	// Fixed amount coupon
	validator.RegisterCoupon(&Coupon{
		CouponID:      "C2",
		Code:          "FLAT500",
		Type:          CouponTypeFixedAmount,
		DiscountValue: 500.0,
		MinPurchase:   2000.0,
		ValidFrom:     now.AddDate(0, -1, 0),
		ValidUntil:    now.AddDate(0, 1, 0),
		UsageLimit:    50,
		UsedCount:     0,
	})

	// Expired coupon
	validator.RegisterCoupon(&Coupon{
		CouponID:      "C3",
		Code:          "EXPIRED",
		Type:          CouponTypePercentage,
		DiscountValue: 20.0,
		MinPurchase:   0.0,
		ValidFrom:     now.AddDate(0, -2, 0),
		ValidUntil:    now.AddDate(0, -1, 0), // Expired
		UsageLimit:    100,
		UsedCount:     0,
	})

	// Buy X Get Y coupon
	validator.RegisterCoupon(&Coupon{
		CouponID:    "C4",
		Code:        "BUY3GET1",
		Type:        CouponTypeBuyXGetY,
		MinPurchase: 0.0,
		ValidFrom:   now.AddDate(0, -1, 0),
		ValidUntil:  now.AddDate(0, 1, 0),
		UsageLimit:  100,
		UsedCount:   0,
		BuyXGetYConfig: &BuyXGetYConfig{
			BuyQuantity:        3,
			GetQuantity:        1,
			ApplicableCategory: "Clothing",
			MinPurchase:        0.0,
		},
	})

	// Free shipping coupon
	validator.RegisterCoupon(&Coupon{
		CouponID:    "C5",
		Code:        "FREESHIP",
		Type:        CouponTypeFreeShipping,
		MinPurchase: 0.0,
		ValidFrom:   now.AddDate(0, -1, 0),
		ValidUntil:  now.AddDate(0, 1, 0),
		UsageLimit:  1000,
		UsedCount:   0,
	})

	// User-specific coupon
	validator.RegisterCoupon(&Coupon{
		CouponID:      "C6",
		Code:          "VIP20",
		Type:          CouponTypePercentage,
		DiscountValue: 20.0,
		MinPurchase:   5000.0,
		ValidFrom:     now.AddDate(0, -1, 0),
		ValidUntil:    now.AddDate(0, 1, 0),
		UsageLimit:    10,
		UsedCount:     0,
		UserSpecific:  true,
		AllowedUsers:  []string{"VIP_USER"},
	})

	// Minimum purchase coupon
	validator.RegisterCoupon(&Coupon{
		CouponID:      "C7",
		Code:          "MIN500",
		Type:          CouponTypePercentage,
		DiscountValue: 15.0,
		MinPurchase:   500.0,
		ValidFrom:     now.AddDate(0, -1, 0),
		ValidUntil:    now.AddDate(0, 1, 0),
		UsageLimit:    100,
		UsedCount:     0,
	})
}

func displayCart(cartService *CartService, cartID string) {
	summary, err := cartService.GetCartSummary(cartID)
	if err != nil {
		fmt.Printf("Error getting cart summary: %v\n", err)
		return
	}

	cart, _ := cartService.GetCart(cartID)

	fmt.Println("\n--- Cart Details ---")
	fmt.Printf("Cart ID: %s\n", cart.CartID)
	fmt.Printf("User ID: %s\n", cart.UserID)
	fmt.Println("\nItems:")
	for _, item := range cart.Items {
		fmt.Printf("  - %s (x%d) @ $%.2f = $%.2f\n",
			item.Name, item.Quantity, item.Price, item.GetTotalPrice())
	}

	fmt.Println("\n--- Summary ---")
	fmt.Printf("Subtotal:     $%.2f\n", summary.Subtotal)
	fmt.Printf("Discount:     $%.2f\n", summary.Discount)
	fmt.Printf("Shipping:     $%.2f\n", summary.Shipping)
	fmt.Printf("Total:        $%.2f\n", summary.Total)
	fmt.Printf("Savings:      $%.2f\n", summary.Savings)

	if len(summary.AppliedCoupons) > 0 {
		fmt.Println("\nApplied Coupons:")
		for _, coupon := range summary.AppliedCoupons {
			fmt.Printf("  - %s (%s)\n", coupon.Code, coupon.Type)
		}
	}
	fmt.Println()
}
