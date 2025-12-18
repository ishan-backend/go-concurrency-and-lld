package main

import (
	"math"
	"sort"
)

// DiscountCalculator calculates discounts for different coupon types
type DiscountCalculator struct{}

// NewDiscountCalculator creates a new discount calculator
func NewDiscountCalculator() *DiscountCalculator {
	return &DiscountCalculator{}
}

// CalculateDiscount calculates the discount amount for a coupon
func (dc *DiscountCalculator) CalculateDiscount(coupon *Coupon, cart *ShoppingCart) float64 {
	subtotal := cart.GetSubtotal()

	switch coupon.Type {
	case CouponTypePercentage:
		return dc.calculatePercentageDiscount(coupon, subtotal)
	case CouponTypeFixedAmount:
		return dc.calculateFixedDiscount(coupon, subtotal)
	case CouponTypeBuyXGetY:
		return dc.calculateBuyXGetYDiscount(coupon, cart)
	case CouponTypeFreeShipping:
		// Free shipping is handled separately, return 0 here
		return 0.0
	default:
		return 0.0
	}
}

// calculatePercentageDiscount calculates percentage-based discount
func (dc *DiscountCalculator) calculatePercentageDiscount(coupon *Coupon, subtotal float64) float64 {
	discount := (subtotal * coupon.DiscountValue) / 100.0
	
	// Apply maximum discount cap if specified
	if coupon.MaxDiscount > 0 && discount > coupon.MaxDiscount {
		discount = coupon.MaxDiscount
	}
	
	return math.Round(discount*100) / 100 // Round to 2 decimal places
}

// calculateFixedDiscount calculates fixed amount discount
func (dc *DiscountCalculator) calculateFixedDiscount(coupon *Coupon, subtotal float64) float64 {
	discount := coupon.DiscountValue
	
	// Don't discount more than the subtotal
	if discount > subtotal {
		discount = subtotal
	}
	
	return math.Round(discount*100) / 100
}

// calculateBuyXGetYDiscount calculates buy X get Y discount
func (dc *DiscountCalculator) calculateBuyXGetYDiscount(coupon *Coupon, cart *ShoppingCart) float64 {
	config := coupon.BuyXGetYConfig
	if config == nil {
		return 0.0
	}

	var applicableItems []*CartItem
	
	// Filter items by category if specified
	if config.ApplicableCategory != "" {
		for _, item := range cart.Items {
			if item.Category == config.ApplicableCategory {
				applicableItems = append(applicableItems, item)
			}
		}
	} else {
		applicableItems = cart.Items
	}

	// Sort items by price (ascending) to maximize discount
	sort.Slice(applicableItems, func(i, j int) bool {
		return applicableItems[i].Price < applicableItems[j].Price
	})

	totalQuantity := 0
	for _, item := range applicableItems {
		totalQuantity += item.Quantity
	}

	// Calculate how many free items we can get
	if totalQuantity < config.BuyQuantity {
		return 0.0
	}

	// Calculate number of times the offer can be applied
	offerCount := totalQuantity / (config.BuyQuantity + config.GetQuantity)
	freeQuantity := offerCount * config.GetQuantity

	// Calculate discount by finding the cheapest items to make free
	discount := 0.0
	remainingFree := freeQuantity
	
	for _, item := range applicableItems {
		if remainingFree <= 0 {
			break
		}
		
		itemsToMakeFree := item.Quantity
		if itemsToMakeFree > remainingFree {
			itemsToMakeFree = remainingFree
		}
		
		discount += item.Price * float64(itemsToMakeFree)
		remainingFree -= itemsToMakeFree
	}

	return math.Round(discount*100) / 100
}

// CalculateTotalDiscount calculates total discount from all applied coupons
func (dc *DiscountCalculator) CalculateTotalDiscount(coupons []*Coupon, cart *ShoppingCart) float64 {
	totalDiscount := 0.0
	
	for _, coupon := range coupons {
		discount := dc.CalculateDiscount(coupon, cart)
		totalDiscount += discount
	}
	
	// Ensure discount doesn't exceed subtotal
	subtotal := cart.GetSubtotal()
	if totalDiscount > subtotal {
		totalDiscount = subtotal
	}
	
	return math.Round(totalDiscount*100) / 100
}

