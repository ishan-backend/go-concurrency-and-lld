package main

import "time"

// CartItem represents an item in the shopping cart
type CartItem struct {
	ProductID string
	Name      string
	Price     float64
	Quantity  int
	Category  string
}

// GetTotalPrice calculates the total price for this item
func (item *CartItem) GetTotalPrice() float64 {
	return item.Price * float64(item.Quantity)
}

// ShoppingCart represents a shopping cart with items
type ShoppingCart struct {
	CartID         string
	UserID         string
	Items          []*CartItem
	AppliedCoupons []*Coupon
	Subtotal       float64
	Discount       float64
	Total          float64
}

// GetSubtotal calculates the subtotal of all items
func (cart *ShoppingCart) GetSubtotal() float64 {
	subtotal := 0.0
	for _, item := range cart.Items {
		subtotal += item.GetTotalPrice()
	}
	return subtotal
}

// CouponType represents the type of coupon discount
type CouponType string

const (
	CouponTypePercentage   CouponType = "PERCENTAGE"
	CouponTypeFixedAmount  CouponType = "FIXED_AMOUNT"
	CouponTypeBuyXGetY     CouponType = "BUY_X_GET_Y"
	CouponTypeFreeShipping CouponType = "FREE_SHIPPING"
)

// Coupon represents a discount coupon
type Coupon struct {
	CouponID             string
	Code                 string
	Type                 CouponType
	DiscountValue        float64 // Percentage (0-100) or fixed amount
	MinPurchase          float64 // Minimum purchase amount required
	MaxDiscount          float64 // Maximum discount cap (for percentage coupons)
	ValidFrom            time.Time
	ValidUntil           time.Time
	UsageLimit           int             // Maximum number of times coupon can be used
	UsedCount            int             // Current usage count
	ApplicableCategories []string        // Categories this coupon applies to
	ExcludedProducts     []string        // Products excluded from coupon
	UserSpecific         bool            // If true, only specific users can use
	AllowedUsers         []string        // List of user IDs allowed to use
	BuyXGetYConfig       *BuyXGetYConfig // Configuration for BUY_X_GET_Y type
}

// BuyXGetYConfig represents configuration for buy X get Y coupons
type BuyXGetYConfig struct {
	BuyQuantity        int     // Buy X items
	GetQuantity        int     // Get Y items free
	ApplicableCategory string  // Category this applies to
	MinPurchase        float64 // Minimum purchase amount
}

// DiscountResult represents the result of applying a coupon
type DiscountResult struct {
	CouponID       string
	CouponCode     string
	DiscountAmount float64
	Applied        bool
	Message        string
}

// CartSummary represents the final cart summary after applying coupons
type CartSummary struct {
	Subtotal       float64
	Discount       float64
	Shipping       float64
	Total          float64
	AppliedCoupons []*Coupon
	Savings        float64
}
