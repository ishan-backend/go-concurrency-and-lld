package main

import (
	"container/list"
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrKeyNotFound      = errors.New("key not found")
	ErrCacheFull        = errors.New("cache is full")
	ErrInvalidKey       = errors.New("key cannot be empty")
	ErrInvalidValue     = errors.New("value cannot be nil")
	ErrInvalidConfig    = errors.New("invalid cache configuration")
	ErrContextCancelled = errors.New("context cancelled")
)

// CacheItem represents a cached entry
type CacheItem struct {
	Key        string
	Value      interface{}
	ExpiresAt  time.Time
	CreatedAt  time.Time
	AccessedAt time.Time
	Frequency  int64 // For LFU
	Size       int64 // In bytes
}

func (c *CacheItem) IsExpired() bool {
	if c.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(c.ExpiresAt)
}

// Storage interface defines how cache data is stored
// This allows extensibility to different storage backends (Redis, disk, etc.)
type Storage interface {
	Get(key string) (*CacheItem, bool)
	Put(key string, item *CacheItem)
	Delete(key string)
	Clear()
	Size() int
	Keys() []string
}

// InMemoryStorage is the default in-memory map-based storage implementation
type InMemoryStorage struct {
	items map[string]*CacheItem
}

// NewInMemoryStorage creates a new in-memory storage instance
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		items: make(map[string]*CacheItem),
	}
}

func (s *InMemoryStorage) Get(key string) (*CacheItem, bool) {
	item, exists := s.items[key]
	return item, exists
}

func (s *InMemoryStorage) Put(key string, item *CacheItem) {
	s.items[key] = item
}

func (s *InMemoryStorage) Delete(key string) {
	delete(s.items, key)
}

func (s *InMemoryStorage) Clear() {
	s.items = make(map[string]*CacheItem)
}

func (s *InMemoryStorage) Size() int {
	return len(s.items)
}

func (s *InMemoryStorage) Keys() []string {
	keys := make([]string, 0, len(s.items))
	for k := range s.items {
		keys = append(keys, k)
	}
	return keys
}

// Cache interface defines core operations
// All methods accept context for cancellation and timeout support
type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Put(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
	Size() int
	Stats() CacheStats
}

// CacheStats provides cache metrics
type CacheStats struct {
	Hits      int64
	Misses    int64
	Evictions int64
	Size      int
	MaxSize   int
	HitRate   float64
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	MaxSize         int
	DefaultTTL      time.Duration
	EvictionPolicy  string // "lru", "lfu", "fifo", "ttl"
	CleanupInterval time.Duration
}

// InMemoryCache is the main cache implementation
// Thread-safety: Uses RWMutex for concurrent read/write operations
// Extensibility: Uses Storage interface for storage backend, EvictionPolicy for eviction strategies
type InMemoryCache struct {
	mu             sync.RWMutex
	storage        Storage // Extensible storage backend
	config         CacheConfig
	evictionPolicy EvictionPolicy // Extensible eviction strategy
	stats          CacheStats
	stopCleanup    chan struct{}
}

// CacheSingleton manages the singleton instance of the cache
// Thread-safety: Uses sync.Once to ensure thread-safe initialization
type CacheSingleton struct {
	mu     sync.Mutex // Protects ResetSingleton
	once   sync.Once  // Ensures initialization happens only once
	cache  *InMemoryCache
	config CacheConfig
}

// Global singleton instance
var globalCacheSingleton = &CacheSingleton{}

// GetInstance returns the singleton instance of the cache
// This method uses sync.Once to ensure thread-safe singleton initialization
// The cache is created only once, even if called from multiple goroutines
func GetInstance(config CacheConfig) (*InMemoryCache, error) {
	var err error
	globalCacheSingleton.once.Do(func() {
		globalCacheSingleton.config = config
		globalCacheSingleton.cache, err = newInMemoryCache(config)
	})
	return globalCacheSingleton.cache, err
}

// ResetSingleton resets the singleton instance (useful for testing)
// WARNING: This should only be used in testing scenarios
// Thread-safety: Not thread-safe - should only be called in test setup/teardown when no other goroutines are accessing it
func ResetSingleton() {
	globalCacheSingleton = &CacheSingleton{}
}

// NewInMemoryCache creates a new cache instance (non-singleton)
func NewInMemoryCache(config CacheConfig) (*InMemoryCache, error) {
	return newInMemoryCache(config)
}

// newInMemoryCache is the internal constructor for creating cache instances
// Validates configuration and initializes storage and eviction policy
func newInMemoryCache(config CacheConfig) (*InMemoryCache, error) {
	// Validate configuration
	if config.MaxSize <= 0 {
		return nil, errors.New("MaxSize must be greater than 0")
	}
	if config.CleanupInterval <= 0 {
		return nil, errors.New("CleanupInterval must be greater than 0")
	}

	cache := &InMemoryCache{
		storage:     NewInMemoryStorage(), // Default in-memory storage
		config:      config,
		stats:       CacheStats{MaxSize: config.MaxSize},
		stopCleanup: make(chan struct{}),
	}

	// Initialize eviction policy based on configuration
	// This demonstrates extensibility - new policies can be added by implementing EvictionPolicy interface
	switch config.EvictionPolicy {
	case "lru":
		cache.evictionPolicy = NewLRUPolicy()
	case "lfu":
		cache.evictionPolicy = NewLFUPolicy()
	case "fifo":
		cache.evictionPolicy = NewFIFOPolicy()
	default:
		cache.evictionPolicy = NewLRUPolicy() // Default to LRU
	}

	// Start background goroutine for expired item cleanup
	// This runs periodically based on CleanupInterval
	go cache.cleanupExpired()

	return cache, nil
}

// SetStorage allows changing the storage backend (for extensibility)
// This enables swapping in-memory storage with Redis, disk-based, or other implementations
func (c *InMemoryCache) SetStorage(storage Storage) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.storage = storage
}

// Get retrieves a value from the cache
// Thread-safety: Uses RWMutex with proper lock ordering to prevent race conditions
// Error handling: Validates input, checks context cancellation, handles expired items
func (c *InMemoryCache) Get(ctx context.Context, key string) (interface{}, error) {
	// Validate input
	if key == "" {
		return nil, ErrInvalidKey
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ErrContextCancelled
	default:
	}

	// Use single lock to prevent race conditions between read and write
	// This fixes the previous issue where item could be modified after lock release
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.storage.Get(key)
	if !exists {
		c.stats.Misses++
		return nil, ErrKeyNotFound
	}

	// Check expiration while holding lock
	if item.IsExpired() {
		c.storage.Delete(key)
		c.evictionPolicy.OnDelete(key)
		c.stats.Misses++
		c.stats.Size = c.storage.Size()
		return nil, ErrKeyNotFound
	}

	// Update access metadata and stats while holding lock
	item.AccessedAt = time.Now()
	item.Frequency++
	c.stats.Hits++
	c.evictionPolicy.OnGet(key)

	return item.Value, nil
}

// Put stores a value in the cache
// Thread-safety: Uses mutex to ensure atomic operations
// Error handling: Validates inputs, checks context, handles eviction when cache is full
func (c *InMemoryCache) Put(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// Validate inputs
	if key == "" {
		return ErrInvalidKey
	}
	if value == nil {
		return ErrInvalidValue
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		return ErrContextCancelled
	default:
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Eviction logic: If cache is full and key doesn't exist, evict an item
	// This demonstrates the extensible eviction policy pattern
	if c.storage.Size() >= c.config.MaxSize {
		if _, exists := c.storage.Get(key); !exists {
			evictKey, ok := c.evictionPolicy.Evict()
			if !ok {
				return ErrCacheFull
			}
			c.storage.Delete(evictKey)
			c.stats.Evictions++
		}
	}

	// Create cache item with metadata
	now := time.Now()
	item := &CacheItem{
		Key:        key,
		Value:      value,
		CreatedAt:  now,
		AccessedAt: now,
		Frequency:  1,
	}

	// Set expiration time
	if ttl > 0 {
		item.ExpiresAt = now.Add(ttl)
	} else if c.config.DefaultTTL > 0 {
		item.ExpiresAt = now.Add(c.config.DefaultTTL)
	}

	// Store item and notify eviction policy
	c.storage.Put(key, item)
	c.evictionPolicy.OnPut(key, item)
	c.stats.Size = c.storage.Size()

	return nil
}

// Delete removes a key from the cache
// Thread-safety: Uses mutex for atomic deletion
func (c *InMemoryCache) Delete(ctx context.Context, key string) error {
	if key == "" {
		return ErrInvalidKey
	}

	select {
	case <-ctx.Done():
		return ErrContextCancelled
	default:
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.storage.Delete(key)
	c.evictionPolicy.OnDelete(key)
	c.stats.Size = c.storage.Size()

	return nil
}

// Clear removes all items from the cache
// Thread-safety: Uses mutex for atomic clearing
func (c *InMemoryCache) Clear(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ErrContextCancelled
	default:
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.storage.Clear()
	c.stats.Size = 0

	return nil
}

// Size returns the current number of items in the cache
// Thread-safety: Uses RLock for concurrent reads
func (c *InMemoryCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.storage.Size()
}

func (c *InMemoryCache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	stats := c.stats
	total := stats.Hits + stats.Misses
	if total > 0 {
		stats.HitRate = float64(stats.Hits) / float64(total)
	}

	return stats
}

func (c *InMemoryCache) cleanupExpired() {
	ticker := time.NewTicker(c.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopCleanup:
			return
		case <-ticker.C:
			c.removeExpiredItems()
		}
	}
}

// removeExpiredItems removes all expired items from the cache
// Called periodically by cleanupExpired goroutine
// Thread-safety: Uses mutex to synchronize with other operations
func (c *InMemoryCache) removeExpiredItems() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	keys := c.storage.Keys()
	for _, key := range keys {
		item, exists := c.storage.Get(key)
		if exists && !item.ExpiresAt.IsZero() && now.After(item.ExpiresAt) {
			c.storage.Delete(key)
			c.evictionPolicy.OnDelete(key)
		}
	}
	c.stats.Size = c.storage.Size()
}

func (c *InMemoryCache) Stop() {
	close(c.stopCleanup)
}

// EvictionPolicy defines cache eviction strategies
// This interface enables extensibility - new eviction policies can be added by implementing this interface
// Examples: LRU, LFU, FIFO, TTL-based, Random, etc.
type EvictionPolicy interface {
	OnGet(key string)                    // Called when an item is accessed
	OnPut(key string, item *CacheItem)   // Called when an item is added/updated
	OnDelete(key string)                 // Called when an item is deleted
	Evict() (string, bool)               // Returns the key to evict, or false if none available
}

// LRU (Least Recently Used) Policy
// Implementation: Uses doubly-linked list to track access order
// Thread-safety: Each policy has its own mutex for thread-safe operations
type LRUPolicy struct {
	mu    sync.Mutex
	list  *list.List                    // Most recently used at front, least at back
	items map[string]*list.Element      // Maps keys to list elements for O(1) access
}

func NewLRUPolicy() *LRUPolicy {
	return &LRUPolicy{
		list:  list.New(),
		items: make(map[string]*list.Element),
	}
}

func (p *LRUPolicy) OnGet(key string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if elem, exists := p.items[key]; exists {
		p.list.MoveToFront(elem)
	}
}

func (p *LRUPolicy) OnPut(key string, item *CacheItem) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if elem, exists := p.items[key]; exists {
		p.list.MoveToFront(elem)
		return
	}

	elem := p.list.PushFront(key)
	p.items[key] = elem
}

func (p *LRUPolicy) OnDelete(key string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if elem, exists := p.items[key]; exists {
		p.list.Remove(elem)
		delete(p.items, key)
	}
}

func (p *LRUPolicy) Evict() (string, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	elem := p.list.Back()
	if elem == nil {
		return "", false
	}

	key := elem.Value.(string)
	p.list.Remove(elem)
	delete(p.items, key)

	return key, true
}

// LFU (Least Frequently Used) Policy
// Implementation: Tracks access frequency for each key
// Thread-safety: Uses mutex to protect frequency map
type LFUPolicy struct {
	mu      sync.Mutex
	freqMap map[string]int64 // Maps keys to their access frequency
	minFreq int64            // Tracks minimum frequency (not currently used in simple implementation)
}

func NewLFUPolicy() *LFUPolicy {
	return &LFUPolicy{
		freqMap: make(map[string]int64),
	}
}

func (p *LFUPolicy) OnGet(key string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.freqMap[key]++
}

func (p *LFUPolicy) OnPut(key string, item *CacheItem) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.freqMap[key] = item.Frequency
}

func (p *LFUPolicy) OnDelete(key string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.freqMap, key)
}

func (p *LFUPolicy) Evict() (string, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.freqMap) == 0 {
		return "", false
	}

	var minKey string
	var minFreq int64 = -1

	for key, freq := range p.freqMap {
		if minFreq == -1 || freq < minFreq {
			minFreq = freq
			minKey = key
		}
	}

	delete(p.freqMap, minKey)
	return minKey, true
}

// FIFO (First In First Out) Policy
// Implementation: Uses queue to maintain insertion order
// Thread-safety: Uses mutex to protect queue operations
type FIFOPolicy struct {
	mu    sync.Mutex
	queue []string        // Queue of keys in insertion order
	items map[string]bool // Tracks which keys are in the queue
}

func NewFIFOPolicy() *FIFOPolicy {
	return &FIFOPolicy{
		queue: make([]string, 0),
		items: make(map[string]bool),
	}
}

func (p *FIFOPolicy) OnGet(key string) {
	// FIFO doesn't change on access
}

func (p *FIFOPolicy) OnPut(key string, item *CacheItem) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.items[key] {
		p.queue = append(p.queue, key)
		p.items[key] = true
	}
}

func (p *FIFOPolicy) OnDelete(key string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.items, key)
}

func (p *FIFOPolicy) Evict() (string, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.queue) == 0 {
		return "", false
	}

	key := p.queue[0]
	p.queue = p.queue[1:]
	delete(p.items, key)

	return key, true
}

// // ShardedCache implements distributed caching across shards
// type ShardedCache struct {
// 	shards    []*InMemoryCache
// 	shardMask uint32
// }

// func NewShardedCache(shardCount int, config CacheConfig) *ShardedCache {
// 	shards := make([]*InMemoryCache, shardCount)
// 	for i := 0; i < shardCount; i++ {
// 		shards[i], _ = NewInMemoryCache(config)
// 	}

// 	return &ShardedCache{
// 		shards:    shards,
// 		shardMask: uint32(shardCount - 1),
// 	}
// }

// func (s *ShardedCache) getShard(key string) *InMemoryCache {
// 	hash := fnv32(key)
// 	return s.shards[hash&s.shardMask]
// }

// func (s *ShardedCache) Get(ctx context.Context, key string) (interface{}, error) {
// 	return s.getShard(key).Get(ctx, key)
// }

// func (s *ShardedCache) Put(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
// 	return s.getShard(key).Put(ctx, key, value, ttl)
// }

// func (s *ShardedCache) Delete(ctx context.Context, key string) error {
// 	return s.getShard(key).Delete(ctx, key)
// }

// func (s *ShardedCache) Clear(ctx context.Context) error {
// 	for _, shard := range s.shards {
// 		if err := shard.Clear(ctx); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (s *ShardedCache) Size() int {
// 	total := 0
// 	for _, shard := range s.shards {
// 		total += shard.Size()
// 	}
// 	return total
// }

// func (s *ShardedCache) Stats() CacheStats {
// 	var combined CacheStats
// 	for _, shard := range s.shards {
// 		stats := shard.Stats()
// 		combined.Hits += stats.Hits
// 		combined.Misses += stats.Misses
// 		combined.Evictions += stats.Evictions
// 		combined.Size += stats.Size
// 	}

// 	total := combined.Hits + combined.Misses
// 	if total > 0 {
// 		combined.HitRate = float64(combined.Hits) / float64(total)
// 	}

// 	return combined
// }

// // FNV-32 hash function
// func fnv32(key string) uint32 {
// 	hash := uint32(2166136261)
// 	const prime32 = uint32(16777619)
// 	for i := 0; i < len(key); i++ {
// 		hash *= prime32
// 		hash ^= uint32(key[i])
// 	}
// 	return hash
// }

// Example main function for testing
func main() {
	ctx := context.Background()

	// Configure cache
	config := CacheConfig{
		MaxSize:         5,
		DefaultTTL:      10 * time.Second,
		EvictionPolicy:  "lru",
		CleanupInterval: 5 * time.Second,
	}

	println("=== Testing Singleton Pattern with sync.Once ===\n")

	// Test 1: Get singleton instance multiple times
	println("Test 1: Singleton Pattern - Multiple GetInstance calls")

	// First call - should create the instance
	cache1, err := GetInstance(config)
	if err != nil {
		panic(err)
	}
	defer cache1.Stop()
	println("  First GetInstance() call - cache created")

	// Second call - should return the same instance
	cache2, err := GetInstance(config)
	if err != nil {
		panic(err)
	}
	println("  Second GetInstance() call - same instance returned")

	// Third call - should return the same instance
	cache3, err := GetInstance(config)
	if err != nil {
		panic(err)
	}
	println("  Third GetInstance() call - same instance returned")

	// Verify they are the same instance
	if cache1 == cache2 && cache2 == cache3 {
		println("  ✓ All instances are the same (singleton working)")
	} else {
		println("  ✗ Instances are different (singleton not working)")
	}
	println()

	// Test 2: Demonstrate singleton behavior with data
	println("Test 2: Singleton Data Persistence")

	// Add data to cache1
	cache1.Put(ctx, "singleton:test", "Hello Singleton", 0)
	println("  Added data to cache1")

	// Get data from cache2 (should be the same instance)
	if val, err := cache2.Get(ctx, "singleton:test"); err == nil {
		println("  Retrieved from cache2:", val.(string))
	}

	// Get data from cache3 (should be the same instance)
	if val, err := cache3.Get(ctx, "singleton:test"); err == nil {
		println("  Retrieved from cache3:", val.(string))
	}
	println()

	// Test 3: Compare with non-singleton cache
	println("Test 3: Comparison with Non-Singleton Cache")

	// Create a non-singleton cache
	nonSingletonCache, err := NewInMemoryCache(config)
	if err != nil {
		panic(err)
	}
	defer nonSingletonCache.Stop()

	// Add data to non-singleton cache
	nonSingletonCache.Put(ctx, "nonsingleton:test", "Hello Non-Singleton", 0)
	println("  Added data to non-singleton cache")

	// Try to get from singleton cache (should not exist)
	if _, err := cache1.Get(ctx, "nonsingleton:test"); err != nil {
		println("  Data not found in singleton cache (expected):", err.Error())
	}

	// Try to get from non-singleton cache (should exist)
	if val, err := nonSingletonCache.Get(ctx, "nonsingleton:test"); err == nil {
		println("  Data found in non-singleton cache:", val.(string))
	}
	println()

	// Continue with original tests using singleton cache
	cache := cache1

	println("=== Testing Cache Operations ===\n")

	// Test 1: Basic Put and Get
	println("Test 1: Basic Put and Get")
	cache.Put(ctx, "user:1", "Alice", 0)
	cache.Put(ctx, "user:2", "Bob", 0)
	cache.Put(ctx, "user:3", "Charlie", 0)

	if val, err := cache.Get(ctx, "user:1"); err == nil {
		println("  user:1 =", val.(string))
	}

	if val, err := cache.Get(ctx, "user:2"); err == nil {
		println("  user:2 =", val.(string))
	}

	println("  Cache size:", cache.Size())
	println()

	// Test 2: Cache Eviction (max size = 5)
	println("Test 2: Testing LRU Eviction (max size = 5)")
	cache.Put(ctx, "user:4", "David", 0)
	cache.Put(ctx, "user:5", "Eve", 0)
	println("  Added 2 more items, size:", cache.Size())

	// Access user:1 to make it recently used
	cache.Get(ctx, "user:1")
	println("  Accessed user:1 to make it recently used")

	// Add one more - should evict user:2 (least recently used)
	cache.Put(ctx, "user:6", "Frank", 0)
	println("  Added user:6, size:", cache.Size())

	// Try to get evicted item
	if _, err := cache.Get(ctx, "user:2"); err != nil {
		println("  user:2 was evicted (expected):", err.Error())
	}

	// user:1 should still exist
	if val, err := cache.Get(ctx, "user:1"); err == nil {
		println("  user:1 still exists:", val.(string))
	}
	println()

	// Test 3: TTL Expiration
	println("Test 3: Testing TTL Expiration")
	cache.Put(ctx, "temp:1", "Temporary Data", 2*time.Second)

	if val, err := cache.Get(ctx, "temp:1"); err == nil {
		println("  temp:1 =", val.(string), "(just added)")
	}

	println("  Waiting 3 seconds for expiration...")
	time.Sleep(3 * time.Second)

	if _, err := cache.Get(ctx, "temp:1"); err != nil {
		println("  temp:1 expired:", err.Error())
	}
	println()

	// Test 4: Cache Stats
	println("Test 4: Cache Statistics")

	// Generate some hits and misses
	cache.Get(ctx, "user:1")
	cache.Get(ctx, "user:3")
	cache.Get(ctx, "user:999") // miss
	cache.Get(ctx, "user:888") // miss

	stats := cache.Stats()
	println("  Total Hits:", stats.Hits)
	println("  Total Misses:", stats.Misses)
	println("  Total Evictions:", stats.Evictions)
	println("  Current Size:", stats.Size)
	println("  Max Size:", stats.MaxSize)
	println("  Hit Rate:", int(stats.HitRate*100), "%")
	println()

	// Test 5: Delete Operation
	println("Test 5: Testing Delete")
	cache.Put(ctx, "delete:me", "Will be deleted", 0)
	println("  Added delete:me")

	cache.Delete(ctx, "delete:me")
	println("  Deleted delete:me")

	if _, err := cache.Get(ctx, "delete:me"); err != nil {
		println("  delete:me not found (expected):", err.Error())
	}
	println()

	// Test 6: Testing different eviction policies
	println("Test 6: Testing LFU Policy")
	lfuConfig := CacheConfig{
		MaxSize:         3,
		DefaultTTL:      0,
		EvictionPolicy:  "lfu",
		CleanupInterval: 5 * time.Second,
	}

	lfuCache, _ := NewInMemoryCache(lfuConfig)
	defer lfuCache.Stop()

	lfuCache.Put(ctx, "a", "A", 0)
	lfuCache.Put(ctx, "b", "B", 0)
	lfuCache.Put(ctx, "c", "C", 0)

	// Access 'a' multiple times
	lfuCache.Get(ctx, "a")
	lfuCache.Get(ctx, "a")
	lfuCache.Get(ctx, "a")

	// Access 'b' once
	lfuCache.Get(ctx, "b")

	// Add new item - should evict 'c' (least frequently used)
	lfuCache.Put(ctx, "d", "D", 0)

	if _, err := lfuCache.Get(ctx, "c"); err != nil {
		println("  'c' was evicted (least frequently used)")
	}

	if val, err := lfuCache.Get(ctx, "a"); err == nil {
		println("  'a' still exists (most frequently used):", val.(string))
	}
	println()

	// Test 7: Thread-safe singleton initialization
	println("Test 7: Thread-Safe Singleton Initialization")

	// Reset singleton for this test
	ResetSingleton()

	// Create multiple goroutines that try to get the singleton instance
	const numGoroutines = 10
	results := make(chan *InMemoryCache, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			cache, _ := GetInstance(config)
			results <- cache
		}()
	}

	// Collect all results
	var instances []*InMemoryCache
	for i := 0; i < numGoroutines; i++ {
		instances = append(instances, <-results)
	}

	// Verify all instances are the same
	allSame := true
	firstInstance := instances[0]
	for i := 1; i < len(instances); i++ {
		if instances[i] != firstInstance {
			allSame = false
			break
		}
	}

	if allSame {
		println("  ✓ All", numGoroutines, "goroutines got the same singleton instance")
	} else {
		println("  ✗ Different instances were created (thread-safety issue)")
	}

	// Clean up
	firstInstance.Stop()
	println()

	println("=== All Tests Completed ===")
}
