package main

import (
	"errors"
	"time"
)

var (
	ErrCouponNotFound          = errors.New("coupon not found")
	ErrCouponExpired           = errors.New("coupon has expired")
	ErrCouponNotStarted        = errors.New("coupon is not yet valid")
	ErrCouponUsageLimitReached = errors.New("coupon usage limit reached")
	ErrMinPurchaseNotMet       = errors.New("minimum purchase amount not met")
	ErrCouponAlreadyApplied    = errors.New("coupon already applied to cart")
	ErrUserNotEligible         = errors.New("user is not eligible for this coupon")
	ErrInvalidCouponType       = errors.New("invalid coupon type")
)

// CouponValidator validates coupon eligibility
type CouponValidator struct {
	coupons map[string]*Coupon
}

// NewCouponValidator creates a new coupon validator
func NewCouponValidator() *CouponValidator {
	return &CouponValidator{
		coupons: make(map[string]*Coupon),
	}
}

// RegisterCoupon registers a coupon in the validator
func (cv *CouponValidator) RegisterCoupon(coupon *Coupon) {
	cv.coupons[coupon.CouponID] = coupon
}

// GetCouponByCode retrieves a coupon by its code
func (cv *CouponValidator) GetCouponByCode(code string) (*Coupon, error) {
	for _, coupon := range cv.coupons {
		if coupon.Code == code {
			return coupon, nil
		}
	}
	return nil, ErrCouponNotFound
}

// ValidateCoupon validates if a coupon can be applied to a cart
func (cv *CouponValidator) ValidateCoupon(coupon *Coupon, cart *ShoppingCart, userID string) error {
	if coupon == nil {
		return ErrCouponNotFound
	}

	// Check if coupon is already applied
	for _, appliedCoupon := range cart.AppliedCoupons {
		if appliedCoupon.CouponID == coupon.CouponID {
			return ErrCouponAlreadyApplied
		}
	}

	// Check validity period
	now := time.Now()
	if now.Before(coupon.ValidFrom) {
		return ErrCouponNotStarted
	}
	if now.After(coupon.ValidUntil) {
		return ErrCouponExpired
	}

	// Check usage limit
	if coupon.UsageLimit > 0 && coupon.UsedCount >= coupon.UsageLimit {
		return ErrCouponUsageLimitReached
	}

	// Check user eligibility
	if coupon.UserSpecific {
		allowed := false
		for _, allowedUser := range coupon.AllowedUsers {
			if allowedUser == userID {
				allowed = true
				break
			}
		}
		if !allowed {
			return ErrUserNotEligible
		}
	}

	// Check minimum purchase amount
	subtotal := cart.GetSubtotal()
	if subtotal < coupon.MinPurchase {
		return ErrMinPurchaseNotMet
	}

	// Validate coupon-specific rules
	switch coupon.Type {
	case CouponTypeBuyXGetY:
		if coupon.BuyXGetYConfig == nil {
			return ErrInvalidCouponType
		}
		return cv.validateBuyXGetY(coupon, cart)
	case CouponTypePercentage, CouponTypeFixedAmount, CouponTypeFreeShipping:
		// These are validated during discount calculation
		return nil
	default:
		return ErrInvalidCouponType
	}
}

// validateBuyXGetY validates buy X get Y coupon eligibility
func (cv *CouponValidator) validateBuyXGetY(coupon *Coupon, cart *ShoppingCart) error {
	config := coupon.BuyXGetYConfig

	// Check minimum purchase if specified
	if config.MinPurchase > 0 {
		subtotal := cart.GetSubtotal()
		if subtotal < config.MinPurchase {
			return ErrMinPurchaseNotMet
		}
	}

	// Check if cart has enough items in the applicable category
	if config.ApplicableCategory != "" {
		totalQuantity := 0
		for _, item := range cart.Items {
			if item.Category == config.ApplicableCategory {
				totalQuantity += item.Quantity
			}
		}
		if totalQuantity < config.BuyQuantity {
			return errors.New("insufficient items in required category")
		}
	} else {
		// Check total quantity across all items
		totalQuantity := 0
		for _, item := range cart.Items {
			totalQuantity += item.Quantity
		}
		if totalQuantity < config.BuyQuantity {
			return errors.New("insufficient items in cart")
		}
	}

	return nil
}

// IsCouponApplicable checks if a coupon is applicable to a cart without checking if it's already applied
// This is useful for getting all applicable coupons
func (cv *CouponValidator) IsCouponApplicable(coupon *Coupon, cart *ShoppingCart, userID string) error {
	if coupon == nil {
		return ErrCouponNotFound
	}

	// Check validity period
	now := time.Now()
	if now.Before(coupon.ValidFrom) {
		return ErrCouponNotStarted
	}
	if now.After(coupon.ValidUntil) {
		return ErrCouponExpired
	}

	// Check usage limit
	if coupon.UsageLimit > 0 && coupon.UsedCount >= coupon.UsageLimit {
		return ErrCouponUsageLimitReached
	}

	// Check user eligibility
	if coupon.UserSpecific {
		allowed := false
		for _, allowedUser := range coupon.AllowedUsers {
			if allowedUser == userID {
				allowed = true
				break
			}
		}
		if !allowed {
			return ErrUserNotEligible
		}
	}

	// Check minimum purchase amount
	subtotal := cart.GetSubtotal()
	if subtotal < coupon.MinPurchase {
		return ErrMinPurchaseNotMet
	}

	// Validate coupon-specific rules
	switch coupon.Type {
	case CouponTypeBuyXGetY:
		if coupon.BuyXGetYConfig == nil {
			return ErrInvalidCouponType
		}
		return cv.validateBuyXGetY(coupon, cart)
	case CouponTypePercentage, CouponTypeFixedAmount, CouponTypeFreeShipping:
		// These are validated during discount calculation
		return nil
	default:
		return ErrInvalidCouponType
	}
}

// GetApplicableCoupons returns all coupons that are applicable to the given cart
func (cv *CouponValidator) GetApplicableCoupons(cart *ShoppingCart, userID string) []*Coupon {
	applicableCoupons := make([]*Coupon, 0)

	for _, coupon := range cv.coupons {
		err := cv.IsCouponApplicable(coupon, cart, userID)
		if err == nil {
			applicableCoupons = append(applicableCoupons, coupon)
		}
	}

	return applicableCoupons
}
