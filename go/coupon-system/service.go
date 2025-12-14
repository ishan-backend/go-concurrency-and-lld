package main

import (
	"sync"
)

// CouponRecommendationService provides coupon recommendation functionality with concurrency support
type CouponRecommendationService struct {
	repository    CouponRepository
	ruleEvaluator *RuleEvaluator
	mu            sync.RWMutex
}

// NewCouponRecommendationService creates a new coupon recommendation service
func NewCouponRecommendationService(repository CouponRepository) *CouponRecommendationService {
	return &CouponRecommendationService{
		repository:    repository,
		ruleEvaluator: NewRuleEvaluator(),
	}
}

// GetApplicableCoupons returns all coupons applicable to the given order
// This method is thread-safe and uses concurrent rule evaluation for performance
func (s *CouponRecommendationService) GetApplicableCoupons(order *Order) ([]*Coupon, error) {
	if order == nil {
		return nil, &ServiceError{Message: "order cannot be nil"}
	}
	
	s.mu.RLock()
	allCoupons := s.repository.GetAllCoupons()
	s.mu.RUnlock()
	
	// Use channels and goroutines for concurrent coupon evaluation
	type couponResult struct {
		coupon *Coupon
		applicable bool
	}
	
	resultChan := make(chan couponResult, len(allCoupons))
	var wg sync.WaitGroup
	
	// Evaluate each coupon concurrently
	for _, coupon := range allCoupons {
		wg.Add(1)
		go func(c *Coupon) {
			defer wg.Done()
			applicable := s.isCouponApplicable(order, c)
			resultChan <- couponResult{
				coupon: c,
				applicable: applicable,
			}
		}(coupon)
	}
	
	// Close channel when all goroutines complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	
	// Collect applicable coupons
	var applicableCoupons []*Coupon
	for result := range resultChan {
		if result.applicable {
			applicableCoupons = append(applicableCoupons, result.coupon)
		}
	}
	
	return applicableCoupons, nil
}

// GetBestCoupon returns the coupon with the highest discount amount
func (s *CouponRecommendationService) GetBestCoupon(order *Order) (*Coupon, error) {
	applicableCoupons, err := s.GetApplicableCoupons(order)
	if err != nil {
		return nil, err
	}
	
	if len(applicableCoupons) == 0 {
		return nil, &ServiceError{Message: "no applicable coupons found"}
	}
	
	bestCoupon := applicableCoupons[0]
	maxDiscount := bestCoupon.DiscountAmount
	
	for _, coupon := range applicableCoupons[1:] {
		if coupon.DiscountAmount > maxDiscount {
			maxDiscount = coupon.DiscountAmount
			bestCoupon = coupon
		}
	}
	
	return bestCoupon, nil
}

// isCouponApplicable checks if a coupon is applicable to an order
func (s *CouponRecommendationService) isCouponApplicable(order *Order, coupon *Coupon) bool {
	if coupon == nil || len(coupon.Rules) == 0 {
		return false
	}
	
	// Use rule evaluator to check coupon eligibility
	return s.ruleEvaluator.EvaluateRulesFormula(order, coupon)
}

// AddCoupon adds a new coupon to the system (thread-safe)
func (s *CouponRecommendationService) AddCoupon(coupon *Coupon) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	return s.repository.AddCoupon(coupon)
}

// ServiceError represents a service operation error
type ServiceError struct {
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}

