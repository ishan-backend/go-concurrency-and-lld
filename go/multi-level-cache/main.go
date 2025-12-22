package main

import (
	"errors"
	"fmt"
	"sync"
)

// DataStore interface for storing key-value pairs
type DataStore interface {
	Get(key string) (string, error)
	Put(key, value string) error
	Delete(key string) error
	Size() int
}

// InMemoryDataStore implements DataStore using in-memory map
type InMemoryDataStore struct {
	mu    sync.Mutex
	Store map[string]string
}

func NewInMemoryDataStore() DataStore {
	return &InMemoryDataStore{
		Store: make(map[string]string),
	}
}

func (i *InMemoryDataStore) Get(key string) (string, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if value, ok := i.Store[key]; ok {
		return value, nil
	}

	return "", errors.New("key not found")
}

func (i *InMemoryDataStore) Put(key, value string) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.Store[key] = value
	return nil
}

func (i *InMemoryDataStore) Delete(key string) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	delete(i.Store, key)
	return nil
}

func (i *InMemoryDataStore) Size() int {
	i.mu.Lock()
	defer i.mu.Unlock()

	return len(i.Store)
}

// EvictionPolicy interface for cache eviction strategies
type EvictionPolicy interface {
	OnGet(key string) (string, error)  // Returns evicted key if eviction needed
	OnPut(key string) (string, error)  // Returns evicted key if eviction needed
	GetEvictionCallback() func(string) // Returns callback to delete from data store
	SetEvictionCallback(func(string))
}

// LRUEvictionPolicy implements LRU (Least Recently Used) eviction policy
type LRUEvictionPolicy struct {
	mu               sync.Mutex
	Capacity         int
	accessOrder      []string       // Most recently used at the end
	keyPosition      map[string]int // Maps key to its index in accessOrder
	evictionCallback func(string)   // Callback to delete from data store
}

func NewLRUEvictionPolicy(capacity int) EvictionPolicy {
	return &LRUEvictionPolicy{
		Capacity:    capacity,
		accessOrder: make([]string, 0),
		keyPosition: make(map[string]int),
	}
}

func (l *LRUEvictionPolicy) GetEvictionCallback() func(string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.evictionCallback
}

func (l *LRUEvictionPolicy) SetEvictionCallback(callback func(string)) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.evictionCallback = callback
}

func (l *LRUEvictionPolicy) OnGet(key string) (string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Move key to end (most recently used)
	l.moveToEnd(key)
	return "", nil
}

func (l *LRUEvictionPolicy) OnPut(key string) (string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Check if key already exists
	if _, exists := l.keyPosition[key]; exists {
		// Key exists, move to end
		l.moveToEnd(key)
		return "", nil
	}

	// New key, add to end
	l.accessOrder = append(l.accessOrder, key)
	l.keyPosition[key] = len(l.accessOrder) - 1

	// Check if eviction is needed
	if len(l.accessOrder) > l.Capacity {
		// Evict least recently used (first element)
		evictedKey := l.accessOrder[0]
		l.accessOrder = l.accessOrder[1:]

		// Update positions
		for i := range l.accessOrder {
			l.keyPosition[l.accessOrder[i]] = i
		}
		delete(l.keyPosition, evictedKey)

		return evictedKey, nil
	}

	return "", nil
}

func (l *LRUEvictionPolicy) moveToEnd(key string) {
	if pos, exists := l.keyPosition[key]; exists {
		// Remove from current position
		l.accessOrder = append(l.accessOrder[:pos], l.accessOrder[pos+1:]...)
		// Add to end
		l.accessOrder = append(l.accessOrder, key)
		// Update positions
		for i := range l.accessOrder {
			l.keyPosition[l.accessOrder[i]] = i
		}
	}
}

// CacheLayerInterface defines operations for a single cache layer
type CacheLayerInterface interface {
	Get(key string) (string, error)
	GetValue(key string) (string, error) // Get value without updating LRU
	Put(key, value string) error
	GetLayerNumber() int
	GetReadTime() int
	GetWriteTime() int
	Contains(key string) bool
}

// CacheLayer represents a single level in the multi-level cache
type CacheLayer struct {
	DataStore      DataStore
	EvictionPolicy EvictionPolicy
	LayerNumber    int
	ReadTime       int
	WriteTime      int
	mu             sync.Mutex
}

func NewCacheLayer(layerNumber int, dataStore DataStore, evictionPolicy EvictionPolicy, readTime int, writeTime int) CacheLayerInterface {
	layer := &CacheLayer{
		LayerNumber:    layerNumber,
		DataStore:      dataStore,
		EvictionPolicy: evictionPolicy,
		ReadTime:       readTime,
		WriteTime:      writeTime,
	}

	// Set up eviction callback
	evictionPolicy.SetEvictionCallback(func(key string) {
		dataStore.Delete(key)
	})

	return layer
}

func (c *CacheLayer) Get(key string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	value, err := c.DataStore.Get(key)
	if err != nil {
		return "", err
	}

	// Update eviction policy on access
	_, _ = c.EvictionPolicy.OnGet(key)
	return value, nil
}

func (c *CacheLayer) GetValue(key string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.DataStore.Get(key)
}

func (c *CacheLayer) Put(key, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if eviction is needed
	evictedKey, err := c.EvictionPolicy.OnPut(key)
	if err == nil && evictedKey != "" {
		// Evict the key from data store
		callback := c.EvictionPolicy.GetEvictionCallback()
		if callback != nil {
			callback(evictedKey)
		}
	}

	// Put the new key-value pair
	return c.DataStore.Put(key, value)
}

func (c *CacheLayer) Contains(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, err := c.DataStore.Get(key)
	return err == nil
}

func (c *CacheLayer) GetLayerNumber() int {
	return c.LayerNumber
}

func (c *CacheLayer) GetReadTime() int {
	return c.ReadTime
}

func (c *CacheLayer) GetWriteTime() int {
	return c.WriteTime
}

// MultiLayerCacheInterface defines operations for multi-level cache
type MultiLayerCacheInterface interface {
	AddLayer(layerNumber int, dataStore DataStore, evictionPolicy EvictionPolicy, readTime int, writeTime int) error
	Get(key string) (string, int, error)
	Put(key, value string) (int, error)
}

// MultiLayerCache implements a multi-level cache system
// L1 is the top layer (fastest, smallest), Ln is the bottom layer (slowest, largest)
type MultiLayerCache struct {
	mu              sync.Mutex
	Layers          []CacheLayerInterface
	LastLayerNumber int
}

func NewMultiLayerCache() *MultiLayerCache {
	return &MultiLayerCache{
		Layers:          make([]CacheLayerInterface, 0),
		LastLayerNumber: -1,
	}
}

func (m *MultiLayerCache) AddLayer(layerNumber int, dataStore DataStore, evictionPolicy EvictionPolicy, readTime int, writeTime int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Validate layer number (should be sequential starting from 0)
	if layerNumber != m.LastLayerNumber+1 {
		return fmt.Errorf("layer number must be %d, got %d", m.LastLayerNumber+1, layerNumber)
	}

	layer := NewCacheLayer(layerNumber, dataStore, evictionPolicy, readTime, writeTime)
	m.Layers = append(m.Layers, layer)
	m.LastLayerNumber = layerNumber

	return nil
}

// Get retrieves a key from the cache
// It searches from L1 to Ln, and promotes the value to all higher layers if found
func (m *MultiLayerCache) Get(key string) (string, int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	readTimeTotal := 0
	writeTimeTotal := 0
	value := ""
	foundLayerIndex := -1

	// Search from L1 (index 0) to Ln (last index)
	for i, layer := range m.Layers {
		readTimeTotal += layer.GetReadTime()

		val, err := layer.Get(key)
		if err == nil && val != "" {
			value = val
			foundLayerIndex = i
			break
		}
	}

	if foundLayerIndex == -1 {
		return "", readTimeTotal, errors.New("key not found in any layer")
	}

	// Promote value to all higher layers (L1 to L(found-1))
	for i := 0; i < foundLayerIndex; i++ {
		writeTimeTotal += m.Layers[i].GetWriteTime()
		if err := m.Layers[i].Put(key, value); err != nil {
			return "", readTimeTotal + writeTimeTotal, fmt.Errorf("error promoting to layer %d: %w", i, err)
		}
	}

	return value, readTimeTotal + writeTimeTotal, nil
}

// Put writes a key-value pair to the cache
// It writes from Ln (last layer) to L1, stopping early if value already matches
func (m *MultiLayerCache) Put(key, value string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	totalTime := 0
	n := len(m.Layers)

	// Write from Ln (last index) to L1 (index 0)
	for i := n - 1; i >= 0; i-- {
		layer := m.Layers[i]

		// Read to check current value (without updating LRU)
		totalTime += layer.GetReadTime()
		currentVal, err := layer.GetValue(key)

		// If value matches, stop writing to higher layers
		if err == nil && currentVal == value {
			break
		}

		// Write the value
		totalTime += layer.GetWriteTime()
		if err := layer.Put(key, value); err != nil {
			return totalTime, fmt.Errorf("error writing to layer %d: %w", i, err)
		}
	}

	return totalTime, nil
}

func main() {
	fmt.Println("=== Multi-Level Cache System Demo ===")
	fmt.Println()

	// Create multi-level cache with 3 layers
	multiLayerCache := NewMultiLayerCache()

	// L1: Fastest, smallest capacity (2 items), Read: 1ms, Write: 2ms
	dataStoreL1 := NewInMemoryDataStore()
	evictionPolicyL1 := NewLRUEvictionPolicy(2)
	if err := multiLayerCache.AddLayer(0, dataStoreL1, evictionPolicyL1, 1, 2); err != nil {
		fmt.Printf("Error adding L1: %v\n", err)
		return
	}

	// L2: Medium speed, medium capacity (4 items), Read: 3ms, Write: 6ms
	dataStoreL2 := NewInMemoryDataStore()
	evictionPolicyL2 := NewLRUEvictionPolicy(4)
	if err := multiLayerCache.AddLayer(1, dataStoreL2, evictionPolicyL2, 3, 6); err != nil {
		fmt.Printf("Error adding L2: %v\n", err)
		return
	}

	// L3: Slowest, larger capacity (3 items), Read: 5ms, Write: 10ms
	dataStoreL3 := NewInMemoryDataStore()
	evictionPolicyL3 := NewLRUEvictionPolicy(3)
	if err := multiLayerCache.AddLayer(2, dataStoreL3, evictionPolicyL3, 5, 10); err != nil {
		fmt.Printf("Error adding L3: %v\n", err)
		return
	}

	fmt.Printf("Created cache with %d layers\n\n", len(multiLayerCache.Layers))

	// Test WRITE operations
	fmt.Println("--- WRITE Operations ---")

	time, err := multiLayerCache.Put("A", "1")
	if err != nil {
		fmt.Printf("Error writing A: %v\n", err)
	} else {
		fmt.Printf("PUT(A, 1) - Total time: %dms\n", time)
	}

	time, err = multiLayerCache.Put("B", "2")
	if err != nil {
		fmt.Printf("Error writing B: %v\n", err)
	} else {
		fmt.Printf("PUT(B, 2) - Total time: %dms\n", time)
	}

	time, err = multiLayerCache.Put("C", "3")
	if err != nil {
		fmt.Printf("Error writing C: %v\n", err)
	} else {
		fmt.Printf("PUT(C, 3) - Total time: %dms\n", time)
	}

	time, err = multiLayerCache.Put("D", "4")
	if err != nil {
		fmt.Printf("Error writing D: %v\n", err)
	} else {
		fmt.Printf("PUT(D, 4) - Total time: %dms\n", time)
	}

	fmt.Println()

	// Test READ operations
	fmt.Println("--- READ Operations ---")

	val, time, err := multiLayerCache.Get("A")
	if err != nil {
		fmt.Printf("Error reading A: %v\n", err)
	} else {
		fmt.Printf("GET(A) - Value: %s, Total time: %dms\n", val, time)
	}

	val, time, err = multiLayerCache.Get("B")
	if err != nil {
		fmt.Printf("Error reading B: %v\n", err)
	} else {
		fmt.Printf("GET(B) - Value: %s, Total time: %dms\n", val, time)
	}

	val, time, err = multiLayerCache.Get("C")
	if err != nil {
		fmt.Printf("Error reading C: %v\n", err)
	} else {
		fmt.Printf("GET(C) - Value: %s, Total time: %dms\n", val, time)
	}

	// Test reading a key that should trigger eviction in L1
	fmt.Println("\n--- Testing Eviction ---")
	val, time, err = multiLayerCache.Get("D")
	if err != nil {
		fmt.Printf("Error reading D: %v\n", err)
	} else {
		fmt.Printf("GET(D) - Value: %s, Total time: %dms\n", val, time)
	}

	// Try to read A again (should be promoted back to L1)
	fmt.Println()
	val, time, err = multiLayerCache.Get("A")
	if err != nil {
		fmt.Printf("Error reading A: %v\n", err)
	} else {
		fmt.Printf("GET(A) again - Value: %s, Total time: %dms\n", val, time)
	}

	fmt.Println("\n=== Demo Complete ===")
}
