package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// Initialize repository and service
	repository := NewInMemoryCouponRepository()
	service := NewCouponRecommendationService(repository)
	
	// Create sample rules
	userTypeRule := &UserTypeRule{RequiredUserType: UserTypePremium}
	productTypeRule := &ProductTypeRule{RequiredProductType: ProductTypeGroceries}
	minAmountRule := &MinimumAmountRule{MinimumAmount: 300.0}
	paymentMethodRule := &PaymentMethodRule{RequiredPaymentMethod: PaymentMethodUPI}
	
	// Create coupons with different rule combinations
	coupon1 := &Coupon{
		CouponID:       "C1",
		DiscountAmount: 100.0,
		Rules:          []Rule{userTypeRule, productTypeRule},
		RulesFormula:   "AND(0, 1)",
	}
	
	coupon2 := &Coupon{
		CouponID:       "C2",
		DiscountAmount: 50.0,
		Rules:          []Rule{productTypeRule},
		RulesFormula:   "0",
	}
	
	coupon3 := &Coupon{
		CouponID:       "C3",
		DiscountAmount: 150.0,
		Rules:          []Rule{minAmountRule, paymentMethodRule},
		RulesFormula:   "AND(0, 1)",
	}
	
	coupon4 := &Coupon{
		CouponID:       "C4",
		DiscountAmount: 75.0,
		Rules:          []Rule{productTypeRule, minAmountRule},
		RulesFormula:   "OR(0, 1)",
	}
	
	// Add coupons to repository
	service.AddCoupon(coupon1)
	service.AddCoupon(coupon2)
	service.AddCoupon(coupon3)
	service.AddCoupon(coupon4)
	
	// Create sample order
	user := &User{
		UserID:        "U1",
		UserType:      UserTypePremium,
		PastPurchases: []string{"P1", "P2"},
	}
	
	cartItems := []CartItem{
		{ProductID: "P1", ProductType: ProductTypeGroceries, Quantity: 2, Price: 200.0},
		{ProductID: "P2", ProductType: ProductTypeBeverages, Quantity: 1, Price: 100.0},
	}
	
	cart := &Cart{
		CartID: "CART1",
		Items:  cartItems,
	}
	
	order := &Order{
		OrderID:      "O1",
		User:         user,
		Cart:         cart,
		TotalAmount:  cart.GetTotalAmount(),
		PaymentMethod: PaymentMethodUPI,
	}
	
	// Test single order recommendation
	fmt.Println("=== Single Order Recommendation ===")
	applicableCoupons, err := service.GetApplicableCoupons(order)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("Order Total: $%.2f\n", order.TotalAmount)
	fmt.Printf("Applicable Coupons: %d\n", len(applicableCoupons))
	for _, coupon := range applicableCoupons {
		fmt.Printf("  - Coupon %s: $%.2f discount\n", coupon.CouponID, coupon.DiscountAmount)
	}
	
	bestCoupon, err := service.GetBestCoupon(order)
	if err == nil {
		fmt.Printf("Best Coupon: %s with $%.2f discount\n", bestCoupon.CouponID, bestCoupon.DiscountAmount)
	}
	
	// Test concurrent requests
	fmt.Println("\n=== Concurrent Request Test ===")
	testConcurrentRequests(service)
	
	// Test performance with multiple orders
	fmt.Println("\n=== Performance Test ===")
	testPerformance(service)
}

// testConcurrentRequests tests the system with concurrent requests
func testConcurrentRequests(service *CouponRecommendationService) {
	orders := createSampleOrders(10)
	
	var wg sync.WaitGroup
	results := make(chan string, len(orders))
	
	startTime := time.Now()
	
	for i, order := range orders {
		wg.Add(1)
		go func(idx int, o *Order) {
			defer wg.Done()
			coupons, err := service.GetApplicableCoupons(o)
			if err != nil {
				results <- fmt.Sprintf("Order %d: Error - %v", idx, err)
				return
			}
			results <- fmt.Sprintf("Order %d: Found %d applicable coupons", idx, len(coupons))
		}(i, order)
	}
	
	wg.Wait()
	close(results)
	
	elapsed := time.Since(startTime)
	
	fmt.Printf("Processed %d concurrent orders in %v\n", len(orders), elapsed)
	for result := range results {
		fmt.Printf("  %s\n", result)
	}
}

// testPerformance tests the performance of the recommendation system
func testPerformance(service *CouponRecommendationService) {
	orders := createSampleOrders(100)
	
	startTime := time.Now()
	for _, order := range orders {
		_, _ = service.GetApplicableCoupons(order)
	}
	elapsed := time.Since(startTime)
	
	fmt.Printf("Processed %d orders sequentially in %v\n", len(orders), elapsed)
	fmt.Printf("Average time per order: %v\n", elapsed/time.Duration(len(orders)))
}

// createSampleOrders creates sample orders for testing
func createSampleOrders(count int) []*Order {
	orders := make([]*Order, count)
	
	userTypes := []UserType{UserTypeRegular, UserTypePremium}
	productTypes := []ProductType{ProductTypeGroceries, ProductTypeBeverages, ProductTypeElectronics}
	paymentMethods := []PaymentMethod{PaymentMethodUPI, PaymentMethodCreditCard, PaymentMethodDebitCard}
	
	for i := 0; i < count; i++ {
		user := &User{
			UserID:        fmt.Sprintf("U%d", i+1),
			UserType:      userTypes[i%len(userTypes)],
			PastPurchases: []string{fmt.Sprintf("P%d", i+1)},
		}
		
		cartItems := []CartItem{
			{
				ProductID:   fmt.Sprintf("P%d", i+1),
				ProductType: productTypes[i%len(productTypes)],
				Quantity:    (i % 5) + 1,
				Price:       float64((i%10)+1) * 50.0,
			},
		}
		
		cart := &Cart{
			CartID: fmt.Sprintf("CART%d", i+1),
			Items:  cartItems,
		}
		
		orders[i] = &Order{
			OrderID:      fmt.Sprintf("O%d", i+1),
			User:         user,
			Cart:         cart,
			TotalAmount:  cart.GetTotalAmount(),
			PaymentMethod: paymentMethods[i%len(paymentMethods)],
		}
	}
	
	return orders
}

