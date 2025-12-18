package main

import (
	"sync"
)

// CartService handles shopping cart operations and coupon application
type CartService struct {
	validator  *CouponValidator
	calculator *DiscountCalculator
	carts      map[string]*ShoppingCart
	mu         sync.RWMutex
}

// NewCartService creates a new cart service
func NewCartService(validator *CouponValidator) *CartService {
	return &CartService{
		validator:  validator,
		calculator: NewDiscountCalculator(),
		carts:      make(map[string]*ShoppingCart),
	}
}

// CreateCart creates a new shopping cart for a user
func (cs *CartService) CreateCart(userID string) *ShoppingCart {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cartID := "CART_" + userID
	cart := &ShoppingCart{
		CartID:         cartID,
		UserID:         userID,
		Items:          make([]*CartItem, 0),
		AppliedCoupons: make([]*Coupon, 0),
		Subtotal:       0.0,
		Discount:       0.0,
		Total:          0.0,
	}

	cs.carts[cartID] = cart
	return cart
}

// GetCart retrieves a cart by cart ID
func (cs *CartService) GetCart(cartID string) (*ShoppingCart, error) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	cart, exists := cs.carts[cartID]
	if !exists {
		return nil, ErrCouponNotFound // Reusing error for simplicity
	}

	return cart, nil
}

// AddItem adds an item to the cart
func (cs *CartService) AddItem(cartID string, item *CartItem) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cart, exists := cs.carts[cartID]
	if !exists {
		return ErrCouponNotFound
	}

	// Check if item already exists, update quantity
	for _, existingItem := range cart.Items {
		if existingItem.ProductID == item.ProductID {
			existingItem.Quantity += item.Quantity
			cs.recalculateCart(cart)
			return nil
		}
	}

	cart.Items = append(cart.Items, item)
	cs.recalculateCart(cart)
	return nil
}

// RemoveItem removes an item from the cart
func (cs *CartService) RemoveItem(cartID string, productID string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cart, exists := cs.carts[cartID]
	if !exists {
		return ErrCouponNotFound
	}

	for i, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			cs.recalculateCart(cart)
			return nil
		}
	}

	return nil
}

// UpdateItemQuantity updates the quantity of an item in the cart
func (cs *CartService) UpdateItemQuantity(cartID string, productID string, quantity int) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cart, exists := cs.carts[cartID]
	if !exists {
		return ErrCouponNotFound
	}

	for _, item := range cart.Items {
		if item.ProductID == productID {
			if quantity <= 0 {
				return cs.RemoveItem(cartID, productID)
			}
			item.Quantity = quantity
			cs.recalculateCart(cart)
			return nil
		}
	}

	return nil
}

// ApplyCoupon applies a coupon to the cart
func (cs *CartService) ApplyCoupon(cartID string, couponCode string) (*DiscountResult, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cart, exists := cs.carts[cartID]
	if !exists {
		return nil, ErrCouponNotFound
	}

	// Get coupon by code
	coupon, err := cs.validator.GetCouponByCode(couponCode)
	if err != nil {
		return &DiscountResult{
			CouponID:       "",
			CouponCode:     couponCode,
			DiscountAmount: 0.0,
			Applied:        false,
			Message:        err.Error(),
		}, err
	}

	// Validate coupon
	err = cs.validator.ValidateCoupon(coupon, cart, cart.UserID)
	if err != nil {
		return &DiscountResult{
			CouponID:       coupon.CouponID,
			CouponCode:     couponCode,
			DiscountAmount: 0.0,
			Applied:        false,
			Message:        err.Error(),
		}, err
	}

	// Apply coupon
	cart.AppliedCoupons = append(cart.AppliedCoupons, coupon)
	coupon.UsedCount++

	// Recalculate cart totals
	cs.recalculateCart(cart)

	discountAmount := cs.calculator.CalculateDiscount(coupon, cart)

	return &DiscountResult{
		CouponID:       coupon.CouponID,
		CouponCode:     couponCode,
		DiscountAmount: discountAmount,
		Applied:        true,
		Message:        "Coupon applied successfully",
	}, nil
}

// RemoveCoupon removes a coupon from the cart
func (cs *CartService) RemoveCoupon(cartID string, couponCode string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cart, exists := cs.carts[cartID]
	if !exists {
		return ErrCouponNotFound
	}

	for i, coupon := range cart.AppliedCoupons {
		if coupon.Code == couponCode {
			cart.AppliedCoupons = append(cart.AppliedCoupons[:i], cart.AppliedCoupons[i+1:]...)
			if coupon.UsedCount > 0 {
				coupon.UsedCount--
			}
			cs.recalculateCart(cart)
			return nil
		}
	}

	return nil
}

// GetCartSummary returns a summary of the cart with all discounts applied
func (cs *CartService) GetCartSummary(cartID string) (*CartSummary, error) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	cart, exists := cs.carts[cartID]
	if !exists {
		return nil, ErrCouponNotFound
	}

	subtotal := cart.GetSubtotal()
	discount := cs.calculator.CalculateTotalDiscount(cart.AppliedCoupons, cart)
	shipping := 0.0

	// Calculate shipping (free shipping coupons handled here)
	for _, coupon := range cart.AppliedCoupons {
		if coupon.Type == CouponTypeFreeShipping {
			shipping = 0.0
			break
		}
	}

	// Default shipping cost (can be made configurable)
	if shipping == 0.0 && subtotal < 500.0 {
		shipping = 50.0
	}

	total := subtotal - discount + shipping

	return &CartSummary{
		Subtotal:       subtotal,
		Discount:       discount,
		Shipping:       shipping,
		Total:          total,
		AppliedCoupons: cart.AppliedCoupons,
		Savings:        discount,
	}, nil
}

// GetApplicableCoupons returns all coupons that are applicable to the given cart
func (cs *CartService) GetApplicableCoupons(cartID string) ([]*Coupon, error) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	cart, exists := cs.carts[cartID]
	if !exists {
		return nil, ErrCouponNotFound
	}

	return cs.validator.GetApplicableCoupons(cart, cart.UserID), nil
}

// recalculateCart recalculates cart totals
func (cs *CartService) recalculateCart(cart *ShoppingCart) {
	cart.Subtotal = cart.GetSubtotal()
	cart.Discount = cs.calculator.CalculateTotalDiscount(cart.AppliedCoupons, cart)
	cart.Total = cart.Subtotal - cart.Discount
}
