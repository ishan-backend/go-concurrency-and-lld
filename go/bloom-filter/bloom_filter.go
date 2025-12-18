package main

import (
	"errors"
	"math"
	"sync"
)

var (
	ErrInvalidSize          = errors.New("bit array size must be greater than 0")
	ErrInvalidHashFunctions = errors.New("number of hash functions must be greater than 0")
	ErrInvalidElement       = errors.New("element cannot be empty")
)

// BloomFilter is a probabilistic data structure that tests membership in a set.
// It can have false positives but never false negatives.
type BloomFilter struct {
	bitArray      *BitSet           // The bit array representing the filter
	numHashFuncs  int               // Number of hash functions (k)
	hashStrategy  HashStrategy      // Hash strategy to use
	size          uint64            // Size of the bit array (m)
	mu            sync.RWMutex      // Read-write mutex for thread-safety
}

// BitSet is a thread-safe bit array implementation.
type BitSet struct {
	bits  []uint64 // Using uint64 slices for efficient storage
	size  uint64
	mu    sync.RWMutex
}

// NewBitSet creates a new BitSet with the specified size.
func NewBitSet(size uint64) *BitSet {
	// Each uint64 holds 64 bits, so we need ceil(size/64) uint64s
	numWords := (size + 63) / 64
	return &BitSet{
		bits: make([]uint64, numWords),
		size: size,
	}
}

// Set sets the bit at the given index to 1.
func (bs *BitSet) Set(index uint64) {
	if index >= bs.size {
		return
	}
	bs.mu.Lock()
	defer bs.mu.Unlock()
	wordIndex := index / 64
	bitIndex := index % 64
	bs.bits[wordIndex] |= 1 << bitIndex
}

// Get returns true if the bit at the given index is set to 1.
func (bs *BitSet) Get(index uint64) bool {
	if index >= bs.size {
		return false
	}
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	wordIndex := index / 64
	bitIndex := index % 64
	return (bs.bits[wordIndex] & (1 << bitIndex)) != 0
}

// Size returns the size of the bit array.
func (bs *BitSet) Size() uint64 {
	return bs.size
}

// BloomFilterConfig holds configuration for creating a BloomFilter.
type BloomFilterConfig struct {
	Size          uint64        // Size of bit array (m)
	NumHashFuncs  int           // Number of hash functions (k)
	HashStrategy  HashStrategy  // Hash strategy to use
}

// NewBloomFilter creates a new BloomFilter with the given configuration.
func NewBloomFilter(config BloomFilterConfig) (*BloomFilter, error) {
	if config.Size == 0 {
		return nil, ErrInvalidSize
	}
	if config.NumHashFuncs <= 0 {
		return nil, ErrInvalidHashFunctions
	}
	if config.HashStrategy == nil {
		config.HashStrategy = NewFNV1aHashStrategy()
	}

	return &BloomFilter{
		bitArray:     NewBitSet(config.Size),
		numHashFuncs: config.NumHashFuncs,
		hashStrategy: config.HashStrategy,
		size:         config.Size,
	}, nil
}

// Add adds an element to the Bloom Filter.
// This operation is idempotent - adding the same element multiple times has no effect.
func (bf *BloomFilter) Add(element string) error {
	if element == "" {
		return ErrInvalidElement
	}

	bf.mu.Lock()
	defer bf.mu.Unlock()

	// Apply k hash functions and set corresponding bits
	for i := 0; i < bf.numHashFuncs; i++ {
		seed := uint64(i)
		index := bf.hashStrategy.Hash(element, seed, bf.size)
		bf.bitArray.Set(index)
	}

	return nil
}

// MightContain checks if an element might be in the Bloom Filter.
// Returns true if the element might be present (possible false positive).
// Returns false if the element is definitely not present (no false negatives).
func (bf *BloomFilter) MightContain(element string) (bool, error) {
	if element == "" {
		return false, ErrInvalidElement
	}

	bf.mu.RLock()
	defer bf.mu.RUnlock()

	// Check all k hash function positions
	for i := 0; i < bf.numHashFuncs; i++ {
		seed := uint64(i)
		index := bf.hashStrategy.Hash(element, seed, bf.size)
		if !bf.bitArray.Get(index) {
			// If any bit is not set, element is definitely not present
			return false, nil
		}
	}

	// All bits are set, element might be present
	return true, nil
}

// Size returns the size of the bit array.
func (bf *BloomFilter) Size() uint64 {
	return bf.size
}

// NumHashFunctions returns the number of hash functions used.
func (bf *BloomFilter) NumHashFunctions() int {
	return bf.numHashFuncs
}

// BloomFilterBuilder provides a fluent API for building BloomFilter instances.
type BloomFilterBuilder struct {
	size          uint64
	numHashFuncs  int
	hashStrategy  HashStrategy
	expectedItems uint64
	falsePosRate  float64
}

// NewBloomFilterBuilder creates a new builder for BloomFilter.
func NewBloomFilterBuilder() *BloomFilterBuilder {
	return &BloomFilterBuilder{
		hashStrategy: NewFNV1aHashStrategy(),
	}
}

// WithSize sets the size of the bit array (m).
func (b *BloomFilterBuilder) WithSize(size uint64) *BloomFilterBuilder {
	b.size = size
	return b
}

// WithNumHashFunctions sets the number of hash functions (k).
func (b *BloomFilterBuilder) WithNumHashFunctions(k int) *BloomFilterBuilder {
	b.numHashFuncs = k
	return b
}

// WithHashStrategy sets the hash strategy to use.
func (b *BloomFilterBuilder) WithHashStrategy(strategy HashStrategy) *BloomFilterBuilder {
	b.hashStrategy = strategy
	return b
}

// WithExpectedItems sets the expected number of items and calculates optimal m and k.
// This is a convenience method that computes optimal values based on expected items
// and desired false positive rate.
func (b *BloomFilterBuilder) WithExpectedItems(expectedItems uint64, falsePositiveRate float64) *BloomFilterBuilder {
	b.expectedItems = expectedItems
	b.falsePosRate = falsePositiveRate

	// Calculate optimal m (bit array size)
	// m = -n * ln(p) / (ln(2)^2)
	// where n = expectedItems, p = falsePositiveRate
	if expectedItems > 0 && falsePositiveRate > 0 && falsePositiveRate < 1 {
		m := -float64(expectedItems) * math.Log(falsePositiveRate) / (math.Log(2) * math.Log(2))
		b.size = uint64(math.Ceil(m))

		// Calculate optimal k (number of hash functions)
		// k = (m/n) * ln(2)
		if b.size > 0 {
			k := (float64(b.size) / float64(expectedItems)) * math.Log(2)
			b.numHashFuncs = int(math.Ceil(k))
			if b.numHashFuncs < 1 {
				b.numHashFuncs = 1
			}
		}
	}

	return b
}

// Build creates a BloomFilter instance with the configured parameters.
func (b *BloomFilterBuilder) Build() (*BloomFilter, error) {
	// If size and numHashFuncs are not set, use defaults or calculate from expected items
	if b.size == 0 {
		if b.expectedItems > 0 && b.falsePosRate > 0 {
			// Recalculate if not already done
			b.WithExpectedItems(b.expectedItems, b.falsePosRate)
		} else {
			// Default values
			b.size = 1000
			if b.numHashFuncs == 0 {
				b.numHashFuncs = 3
			}
		}
	}

	if b.numHashFuncs == 0 {
		b.numHashFuncs = 3
	}

	if b.hashStrategy == nil {
		b.hashStrategy = NewFNV1aHashStrategy()
	}

	return NewBloomFilter(BloomFilterConfig{
		Size:         b.size,
		NumHashFuncs: b.numHashFuncs,
		HashStrategy: b.hashStrategy,
	})
}

