package main

import (
	"container/list"
	"fmt"
	"hash/fnv"
)

// Entry represents a key-value pair in the HashMap
type Entry struct {
	key   string
	value interface{}
	hash  uint32
}

// HashMap structure with collision handling using chaining with container/list
type HashMap struct {
	buckets    []*list.List // Array of linked lists for chaining
	capacity   int
	size       int
	loadFactor float64
}

const (
	defaultCapacity   = 16
	defaultLoadFactor = 0.75
)

// NewHashMap creates a new HashMap with default capacity
func NewHashMap() *HashMap {
	buckets := make([]*list.List, defaultCapacity)
	for i := range buckets {
		buckets[i] = list.New()
	}
	return &HashMap{
		buckets:    buckets,
		capacity:   defaultCapacity,
		size:       0,
		loadFactor: defaultLoadFactor,
	}
}

// NewHashMapWithCapacity creates a new HashMap with specified capacity
func NewHashMapWithCapacity(capacity int) *HashMap {
	buckets := make([]*list.List, capacity)
	for i := range buckets {
		buckets[i] = list.New()
	}
	return &HashMap{
		buckets:    buckets,
		capacity:   capacity,
		size:       0,
		loadFactor: defaultLoadFactor,
	}
}

// hash computes the hash value for a given key
func (hm *HashMap) hash(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

// getBucketIndex returns the bucket index for a given hash
func (hm *HashMap) getBucketIndex(hash uint32) int {
	return int(hash % uint32(hm.capacity))
}

// Put inserts or updates a key-value pair in the HashMap
// Handles collision using chaining with container/list
func (hm *HashMap) Put(key string, value interface{}) {
	hash := hm.hash(key)
	bucketIndex := hm.getBucketIndex(hash)
	bucket := hm.buckets[bucketIndex]

	// Check if key already exists (update value)
	for e := bucket.Front(); e != nil; e = e.Next() {
		entry := e.Value.(Entry)
		if entry.key == key {
			// Update existing entry
			e.Value = Entry{
				key:   key,
				value: value,
				hash:  hash,
			}
			return
		}
	}

	// Add new entry to the end of the list
	newEntry := Entry{
		key:   key,
		value: value,
		hash:  hash,
	}
	bucket.PushBack(newEntry)
	hm.size++

	// Check if rehashing is needed
	if float64(hm.size)/float64(hm.capacity) > hm.loadFactor {
		hm.rehash()
	}
}

// Get retrieves the value associated with a key
func (hm *HashMap) Get(key string) (interface{}, bool) {
	hash := hm.hash(key)
	bucketIndex := hm.getBucketIndex(hash)
	bucket := hm.buckets[bucketIndex]

	// Search for key in the list
	for e := bucket.Front(); e != nil; e = e.Next() {
		entry := e.Value.(Entry)
		if entry.key == key {
			return entry.value, true
		}
	}

	return nil, false
}

// Remove deletes a key-value pair from the HashMap
func (hm *HashMap) Remove(key string) (interface{}, bool) {
	hash := hm.hash(key)
	bucketIndex := hm.getBucketIndex(hash)
	bucket := hm.buckets[bucketIndex]

	// Find and remove entry from the list
	for e := bucket.Front(); e != nil; e = e.Next() {
		entry := e.Value.(Entry)
		if entry.key == key {
			value := entry.value
			bucket.Remove(e)
			hm.size--
			return value, true
		}
	}

	return nil, false
}

// ContainsKey checks if a key exists in the HashMap
func (hm *HashMap) ContainsKey(key string) bool {
	_, exists := hm.Get(key)
	return exists
}

// Size returns the number of entries in the HashMap
func (hm *HashMap) Size() int {
	return hm.size
}

// IsEmpty checks if the HashMap is empty
func (hm *HashMap) IsEmpty() bool {
	return hm.size == 0
}

// rehash resizes and redistributes all entries
func (hm *HashMap) rehash() {
	oldBuckets := hm.buckets
	oldCapacity := hm.capacity

	// Double the capacity
	hm.capacity = hm.capacity * 2
	hm.buckets = make([]*list.List, hm.capacity)
	for i := range hm.buckets {
		hm.buckets[i] = list.New()
	}
	hm.size = 0

	// Redistribute all entries
	for i := 0; i < oldCapacity; i++ {
		bucket := oldBuckets[i]
		for e := bucket.Front(); e != nil; e = e.Next() {
			entry := e.Value.(Entry)
			hm.Put(entry.key, entry.value)
		}
	}
}

// Display prints the HashMap contents (for debugging)
func (hm *HashMap) Display() {
	fmt.Println("HashMap Contents:")
	for i := 0; i < hm.capacity; i++ {
		bucket := hm.buckets[i]
		if bucket.Len() > 0 {
			fmt.Printf("Bucket %d: ", i)
			for e := bucket.Front(); e != nil; e = e.Next() {
				entry := e.Value.(Entry)
				fmt.Printf("[%s=%v] ", entry.key, entry.value)
			}
			fmt.Println()
		}
	}
	fmt.Printf("Size: %d, Capacity: %d\n", hm.size, hm.capacity)
}

// Keys returns all keys in the HashMap
func (hm *HashMap) Keys() []string {
	keys := make([]string, 0, hm.size)
	for i := 0; i < hm.capacity; i++ {
		bucket := hm.buckets[i]
		for e := bucket.Front(); e != nil; e = e.Next() {
			entry := e.Value.(Entry)
			keys = append(keys, entry.key)
		}
	}
	return keys
}

// Values returns all values in the HashMap
func (hm *HashMap) Values() []interface{} {
	values := make([]interface{}, 0, hm.size)
	for i := 0; i < hm.capacity; i++ {
		bucket := hm.buckets[i]
		for e := bucket.Front(); e != nil; e = e.Next() {
			entry := e.Value.(Entry)
			values = append(values, entry.value)
		}
	}
	return values
}

// Clear removes all entries from the HashMap
func (hm *HashMap) Clear() {
	hm.buckets = make([]*list.List, defaultCapacity)
	for i := range hm.buckets {
		hm.buckets[i] = list.New()
	}
	hm.capacity = defaultCapacity
	hm.size = 0
}

func main() {
	// Create new HashMap
	hm := NewHashMap()

	// Test basic operations
	fmt.Println("=== Testing Basic Operations ===")
	hm.Put("apple", 100)
	hm.Put("banana", 200)
	hm.Put("orange", 300)
	hm.Put("grape", 400)

	if val, ok := hm.Get("apple"); ok {
		fmt.Printf("Get 'apple': %v\n", val)
	}
	if val, ok := hm.Get("banana"); ok {
		fmt.Printf("Get 'banana': %v\n", val)
	}

	// Test collision handling
	fmt.Println("\n=== Testing Collision Handling ===")
	hm.Put("mango", 500)
	hm.Put("peach", 600)
	hm.Display()

	// Test update
	fmt.Println("\n=== Testing Update ===")
	hm.Put("apple", 150)
	if val, ok := hm.Get("apple"); ok {
		fmt.Printf("Get 'apple' after update: %v\n", val)
	}

	// Test remove
	fmt.Println("\n=== Testing Remove ===")
	if val, ok := hm.Remove("banana"); ok {
		fmt.Printf("Removed 'banana': %v\n", val)
	}
	if _, ok := hm.Get("banana"); !ok {
		fmt.Println("Get 'banana' after remove: not found")
	}

	// Test rehashing
	fmt.Println("\n=== Testing Rehashing ===")
	for i := 0; i < 20; i++ {
		hm.Put(fmt.Sprintf("key%d", i), i*10)
	}
	hm.Display()

	// Test additional methods
	fmt.Println("\n=== Testing Additional Methods ===")
	fmt.Printf("All keys: %v\n", hm.Keys())
	fmt.Printf("Contains 'apple': %v\n", hm.ContainsKey("apple"))
	fmt.Printf("Contains 'banana': %v\n", hm.ContainsKey("banana"))

	fmt.Println("\n=== Final Stats ===")
	fmt.Printf("Size: %d\n", hm.Size())
	fmt.Printf("Is Empty: %v\n", hm.IsEmpty())
}
