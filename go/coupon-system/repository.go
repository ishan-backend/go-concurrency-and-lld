package main

import (
	"sync"
)

// CouponRepository manages the storage and retrieval of coupons with thread safety
type CouponRepository interface {
	AddCoupon(coupon *Coupon) error
	GetCoupon(couponID string) (*Coupon, error)
	GetAllCoupons() []*Coupon
	RemoveCoupon(couponID string) error
}

// InMemoryCouponRepository is an in-memory implementation of CouponRepository
type InMemoryCouponRepository struct {
	coupons map[string]*Coupon
	mu      sync.RWMutex
}

// NewInMemoryCouponRepository creates a new in-memory coupon repository
func NewInMemoryCouponRepository() *InMemoryCouponRepository {
	return &InMemoryCouponRepository{
		coupons: make(map[string]*Coupon),
	}
}

// AddCoupon adds a coupon to the repository (thread-safe)
func (r *InMemoryCouponRepository) AddCoupon(coupon *Coupon) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if coupon == nil {
		return &RepositoryError{Message: "coupon cannot be nil"}
	}
	
	if coupon.CouponID == "" {
		return &RepositoryError{Message: "coupon ID cannot be empty"}
	}
	
	r.coupons[coupon.CouponID] = coupon
	return nil
}

// GetCoupon retrieves a coupon by ID (thread-safe)
func (r *InMemoryCouponRepository) GetCoupon(couponID string) (*Coupon, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	coupon, exists := r.coupons[couponID]
	if !exists {
		return nil, &RepositoryError{Message: "coupon not found: " + couponID}
	}
	
	return coupon, nil
}

// GetAllCoupons retrieves all coupons (thread-safe)
func (r *InMemoryCouponRepository) GetAllCoupons() []*Coupon {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	coupons := make([]*Coupon, 0, len(r.coupons))
	for _, coupon := range r.coupons {
		coupons = append(coupons, coupon)
	}
	
	return coupons
}

// RemoveCoupon removes a coupon from the repository (thread-safe)
func (r *InMemoryCouponRepository) RemoveCoupon(couponID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.coupons[couponID]; !exists {
		return &RepositoryError{Message: "coupon not found: " + couponID}
	}
	
	delete(r.coupons, couponID)
	return nil
}

// RepositoryError represents a repository operation error
type RepositoryError struct {
	Message string
}

func (e *RepositoryError) Error() string {
	return e.Message
}

