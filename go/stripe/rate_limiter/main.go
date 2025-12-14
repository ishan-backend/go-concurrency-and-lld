package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// RateLimiter implements rate limiting using sliding window algorithm
type RateLimiter struct {
	mu              sync.RWMutex
	requests        map[string][]time.Time // userID -> list of request timestamps
	maxRequests     int
	windowDuration  time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxRequests int, windowDuration time.Duration) *RateLimiter {
	return &RateLimiter{
		requests:       make(map[string][]time.Time),
		maxRequests:    maxRequests,
		windowDuration: windowDuration,
	}
}

// IsAllowed checks if a request is allowed for the given userID
// Returns true if allowed, false if rate limited
func (r *RateLimiter) IsAllowed(userID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	now := time.Now()
	windowStart := now.Add(-r.windowDuration)
	
	// Get or create request history for user
	requests, exists := r.requests[userID]
	if !exists {
		r.requests[userID] = []time.Time{now}
		return true
	}
	
	// Remove requests outside the window
	validRequests := make([]time.Time, 0)
	for _, reqTime := range requests {
		if reqTime.After(windowStart) {
			validRequests = append(validRequests, reqTime)
		}
	}
	
	// Check if we're at the limit
	if len(validRequests) >= r.maxRequests {
		return false
	}
	
	// Add current request
	validRequests = append(validRequests, now)
	r.requests[userID] = validRequests
	return true
}

// GetRemainingRequests returns the number of remaining requests for a user
func (r *RateLimiter) GetRemainingRequests(userID string) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	now := time.Now()
	windowStart := now.Add(-r.windowDuration)
	
	requests, exists := r.requests[userID]
	if !exists {
		return r.maxRequests
	}
	
	// Count valid requests
	count := 0
	for _, reqTime := range requests {
		if reqTime.After(windowStart) {
			count++
		}
	}
	
	return r.maxRequests - count
}

// Reset clears all rate limit data (useful for testing)
func (r *RateLimiter) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.requests = make(map[string][]time.Time)
}

// TokenBucketRateLimiter implements rate limiting using token bucket algorithm
type TokenBucketRateLimiter struct {
	mu              sync.RWMutex
	buckets         map[string]*TokenBucket
	maxTokens       int
	refillRate      float64 // tokens per second
	refillInterval  time.Duration
}

// TokenBucket represents a token bucket for a user
type TokenBucket struct {
	tokens     float64
	lastRefill time.Time
}

// NewTokenBucketRateLimiter creates a new token bucket rate limiter
func NewTokenBucketRateLimiter(maxTokens int, refillRate float64, refillInterval time.Duration) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		buckets:        make(map[string]*TokenBucket),
		maxTokens:      maxTokens,
		refillRate:     refillRate,
		refillInterval: refillInterval,
	}
}

// IsAllowed checks if a request is allowed using token bucket
func (t *TokenBucketRateLimiter) IsAllowed(userID string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	
	now := time.Now()
	
	// Get or create bucket for user
	bucket, exists := t.buckets[userID]
	if !exists {
		t.buckets[userID] = &TokenBucket{
			tokens:     float64(t.maxTokens),
			lastRefill: now,
		}
		return true
	}
	
	// Refill tokens based on time elapsed
	elapsed := now.Sub(bucket.lastRefill).Seconds()
	tokensToAdd := elapsed * t.refillRate
	bucket.tokens = min(bucket.tokens+tokensToAdd, float64(t.maxTokens))
	bucket.lastRefill = now
	
	// Check if we have tokens
	if bucket.tokens >= 1.0 {
		bucket.tokens -= 1.0
		return true
	}
	
	return false
}

// GetRemainingTokens returns the number of remaining tokens for a user
func (t *TokenBucketRateLimiter) GetRemainingTokens(userID string) int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	
	bucket, exists := t.buckets[userID]
	if !exists {
		return t.maxTokens
	}
	
	now := time.Now()
	elapsed := now.Sub(bucket.lastRefill).Seconds()
	tokensToAdd := elapsed * t.refillRate
	currentTokens := min(bucket.tokens+tokensToAdd, float64(t.maxTokens))
	
	return int(currentTokens)
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Rate Limiter - Sliding Window")
	fmt.Println(strings.Repeat("=", 60))
	
	// Example 1: Sliding Window Rate Limiter
	fmt.Println("\nExample 1: Sliding Window Rate Limiter")
	fmt.Println("-" + strings.Repeat("-", 60))
	fmt.Println("Limit: 5 requests per 10 seconds")
	
	limiter := NewRateLimiter(5, 10*time.Second)
	
	// Simulate requests
	userID := "user1"
	for i := 0; i < 7; i++ {
		allowed := limiter.IsAllowed(userID)
		remaining := limiter.GetRemainingRequests(userID)
		fmt.Printf("Request %d: Allowed=%v, Remaining=%d\n", i+1, allowed, remaining)
		time.Sleep(500 * time.Millisecond)
	}
	
	// Example 2: Rate limit exceeded
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 2: Rate Limit Exceeded")
	fmt.Println("-" + strings.Repeat("-", 60))
	fmt.Println("Limit: 3 requests per 5 seconds")
	
	limiter2 := NewRateLimiter(3, 5*time.Second)
	userID2 := "user2"
	
	for i := 0; i < 5; i++ {
		allowed := limiter2.IsAllowed(userID2)
		remaining := limiter2.GetRemainingRequests(userID2)
		fmt.Printf("Request %d: Allowed=%v, Remaining=%d\n", i+1, allowed, remaining)
	}
	
	// Example 3: Multiple users
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 3: Multiple Users")
	fmt.Println("-" + strings.Repeat("-", 60))
	fmt.Println("Limit: 2 requests per 3 seconds per user")
	
	limiter3 := NewRateLimiter(2, 3*time.Second)
	
	users := []string{"alice", "bob", "alice", "bob", "alice"}
	for i, uid := range users {
		allowed := limiter3.IsAllowed(uid)
		remaining := limiter3.GetRemainingRequests(uid)
		fmt.Printf("User %s, Request %d: Allowed=%v, Remaining=%d\n", uid, i+1, allowed, remaining)
	}
	
	// Example 4: Token Bucket Rate Limiter
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Rate Limiter - Token Bucket")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("\nExample 4: Token Bucket Rate Limiter")
	fmt.Println("-" + strings.Repeat("-", 60))
	fmt.Println("Max tokens: 5, Refill rate: 1 token per second")
	
	tokenLimiter := NewTokenBucketRateLimiter(5, 1.0, time.Second)
	userID3 := "user3"
	
	for i := 0; i < 7; i++ {
		allowed := tokenLimiter.IsAllowed(userID3)
		remaining := tokenLimiter.GetRemainingTokens(userID3)
		fmt.Printf("Request %d: Allowed=%v, Remaining tokens=%d\n", i+1, allowed, remaining)
		if i < 4 {
			time.Sleep(200 * time.Millisecond)
		} else {
			time.Sleep(1 * time.Second) // Wait for refill
		}
	}
	
	// Example 5: Token bucket with slow requests
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 5: Token Bucket - Slow Requests")
	fmt.Println("-" + strings.Repeat("-", 60))
	fmt.Println("Max tokens: 3, Refill rate: 0.5 tokens per second")
	
	tokenLimiter2 := NewTokenBucketRateLimiter(3, 0.5, time.Second)
	userID4 := "user4"
	
	fmt.Println("Making 3 requests quickly, then waiting for refill...")
	for i := 0; i < 3; i++ {
		allowed := tokenLimiter2.IsAllowed(userID4)
		remaining := tokenLimiter2.GetRemainingTokens(userID4)
		fmt.Printf("Request %d: Allowed=%v, Remaining tokens=%d\n", i+1, allowed, remaining)
	}
	
	fmt.Println("\nWaiting 3 seconds for token refill...")
	time.Sleep(3 * time.Second)
	
	allowed := tokenLimiter2.IsAllowed(userID4)
	remaining := tokenLimiter2.GetRemainingTokens(userID4)
	fmt.Printf("Request after wait: Allowed=%v, Remaining tokens=%d\n", allowed, remaining)
}

