package main

// HashStrategy defines the contract for hash functions used in the Bloom Filter.
// Different implementations can provide different hashing algorithms.
type HashStrategy interface {
	// Hash computes a hash value for the input string using a seed.
	// The seed allows us to simulate multiple independent hash functions.
	// Returns a value in the range [0, bound)
	Hash(input string, seed uint64, bound uint64) uint64
}

// FNV1aHashStrategy implements the FNV-1a (64-bit) hashing algorithm.
// It's fast, simple, and has good dispersion for small strings.
type FNV1aHashStrategy struct{}

const (
	fnvOffsetBasis uint64 = 14695981039346656037
	fnvPrime       uint64 = 1099511628211
)

func NewFNV1aHashStrategy() *FNV1aHashStrategy {
	return &FNV1aHashStrategy{}
}

func (h *FNV1aHashStrategy) Hash(input string, seed uint64, bound uint64) uint64 {
	hash := fnvOffsetBasis ^ seed
	for i := 0; i < len(input); i++ {
		hash ^= uint64(input[i])
		hash *= fnvPrime
	}
	return hash % bound
}

// DJB2HashStrategy implements the DJB2 hash function.
// Known for simplicity and relatively low collision rate.
type DJB2HashStrategy struct{}

func NewDJB2HashStrategy() *DJB2HashStrategy {
	return &DJB2HashStrategy{}
}

func (h *DJB2HashStrategy) Hash(input string, seed uint64, bound uint64) uint64 {
	hash := uint64(5381) + seed
	for i := 0; i < len(input); i++ {
		hash = ((hash << 5) + hash) + uint64(input[i])
	}
	return hash % bound
}

// MurmurHashStrategy implements a simplified MurmurHash-like algorithm.
// Provides good distribution and is commonly used in Bloom Filters.
type MurmurHashStrategy struct{}

func NewMurmurHashStrategy() *MurmurHashStrategy {
	return &MurmurHashStrategy{}
}

func (h *MurmurHashStrategy) Hash(input string, seed uint64, bound uint64) uint64 {
	const (
		m = uint64(0xc6a4a7935bd1e995)
		r = 47
	)

	hash := seed ^ (uint64(len(input)) * m)

	for i := 0; i < len(input); i++ {
		k := uint64(input[i])
		k *= m
		k ^= k >> r
		k *= m
		hash ^= k
		hash *= m
	}

	hash ^= hash >> r
	hash *= m
	hash ^= hash >> r

	return hash % bound
}

