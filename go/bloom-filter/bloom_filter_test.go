package main

import (
	"testing"
)

func TestBloomFilter_AddAndMightContain(t *testing.T) {
	bf, err := NewBloomFilterBuilder().
		WithSize(1000).
		WithNumHashFunctions(3).
		Build()

	if err != nil {
		t.Fatalf("Failed to create Bloom Filter: %v", err)
	}

	// Add elements
	elements := []string{"apple", "banana", "cherry"}
	for _, elem := range elements {
		if err := bf.Add(elem); err != nil {
			t.Errorf("Failed to add element %s: %v", elem, err)
		}
	}

	// Check that added elements are found (no false negatives)
	for _, elem := range elements {
		mightContain, err := bf.MightContain(elem)
		if err != nil {
			t.Errorf("Error checking element %s: %v", elem, err)
		}
		if !mightContain {
			t.Errorf("False negative: element %s was added but not found", elem)
		}
	}
}

func TestBloomFilter_NoFalseNegatives(t *testing.T) {
	bf, err := NewBloomFilterBuilder().
		WithSize(1000).
		WithNumHashFunctions(5).
		Build()

	if err != nil {
		t.Fatalf("Failed to create Bloom Filter: %v", err)
	}

	// Add many elements
	for i := 0; i < 100; i++ {
		elem := string(rune('a' + i%26))
		bf.Add(elem)
	}

	// Verify all added elements are found
	for i := 0; i < 100; i++ {
		elem := string(rune('a' + i%26))
		mightContain, _ := bf.MightContain(elem)
		if !mightContain {
			t.Errorf("False negative: element %s was added but not found", elem)
		}
	}
}

func TestBloomFilter_EmptyElement(t *testing.T) {
	bf, err := NewBloomFilterBuilder().
		WithSize(1000).
		WithNumHashFunctions(3).
		Build()

	if err != nil {
		t.Fatalf("Failed to create Bloom Filter: %v", err)
	}

	// Adding empty element should return error
	err = bf.Add("")
	if err == nil {
		t.Error("Expected error when adding empty element")
	}

	// Checking empty element should return error
	_, err = bf.MightContain("")
	if err == nil {
		t.Error("Expected error when checking empty element")
	}
}

func TestBloomFilter_InvalidConfig(t *testing.T) {
	// Test with zero size - use NewBloomFilter directly since builder uses defaults
	_, err := NewBloomFilter(BloomFilterConfig{
		Size:         0,
		NumHashFuncs: 3,
		HashStrategy: NewFNV1aHashStrategy(),
	})

	if err == nil {
		t.Error("Expected error when creating Bloom Filter with zero size")
	}

	// Test with zero hash functions
	_, err = NewBloomFilter(BloomFilterConfig{
		Size:         1000,
		NumHashFuncs: 0,
		HashStrategy: NewFNV1aHashStrategy(),
	})

	if err == nil {
		t.Error("Expected error when creating Bloom Filter with zero hash functions")
	}
}

func TestBloomFilterBuilder_WithExpectedItems(t *testing.T) {
	bf, err := NewBloomFilterBuilder().
		WithExpectedItems(1000, 0.01).
		Build()

	if err != nil {
		t.Fatalf("Failed to create Bloom Filter: %v", err)
	}

	// Verify that size and hash functions are set
	if bf.Size() == 0 {
		t.Error("Bit array size should be greater than 0")
	}

	if bf.NumHashFunctions() == 0 {
		t.Error("Number of hash functions should be greater than 0")
	}

	// Verify optimal configuration (rough checks)
	if bf.Size() < 100 {
		t.Errorf("Bit array size seems too small: %d", bf.Size())
	}

	if bf.NumHashFunctions() < 1 || bf.NumHashFunctions() > 20 {
		t.Errorf("Number of hash functions seems unreasonable: %d", bf.NumHashFunctions())
	}
}

func TestBloomFilter_Idempotent(t *testing.T) {
	bf, err := NewBloomFilterBuilder().
		WithSize(1000).
		WithNumHashFunctions(3).
		Build()

	if err != nil {
		t.Fatalf("Failed to create Bloom Filter: %v", err)
	}

	elem := "test-element"

	// Add the same element multiple times
	for i := 0; i < 10; i++ {
		if err := bf.Add(elem); err != nil {
			t.Errorf("Failed to add element: %v", err)
		}
	}

	// Should still be found
	mightContain, _ := bf.MightContain(elem)
	if !mightContain {
		t.Error("Element should be found after multiple additions")
	}
}

func TestHashStrategy_FNV1a(t *testing.T) {
	strategy := NewFNV1aHashStrategy()
	
	hash1 := strategy.Hash("test", 0, 1000)
	hash2 := strategy.Hash("test", 0, 1000)
	
	// Same input should produce same hash
	if hash1 != hash2 {
		t.Error("Hash function should be deterministic")
	}

	// Different seeds should produce different hashes
	hash3 := strategy.Hash("test", 1, 1000)
	if hash1 == hash3 {
		t.Error("Different seeds should produce different hashes")
	}

	// Hash should be within bounds
	if hash1 >= 1000 {
		t.Errorf("Hash value %d should be less than bound 1000", hash1)
	}
}

func TestHashStrategy_DJB2(t *testing.T) {
	strategy := NewDJB2HashStrategy()
	
	hash1 := strategy.Hash("test", 0, 1000)
	hash2 := strategy.Hash("test", 0, 1000)
	
	if hash1 != hash2 {
		t.Error("Hash function should be deterministic")
	}

	if hash1 >= 1000 {
		t.Errorf("Hash value %d should be less than bound 1000", hash1)
	}
}

func TestHashStrategy_Murmur(t *testing.T) {
	strategy := NewMurmurHashStrategy()
	
	hash1 := strategy.Hash("test", 0, 1000)
	hash2 := strategy.Hash("test", 0, 1000)
	
	if hash1 != hash2 {
		t.Error("Hash function should be deterministic")
	}

	if hash1 >= 1000 {
		t.Errorf("Hash value %d should be less than bound 1000", hash1)
	}
}

func TestBitSet_SetAndGet(t *testing.T) {
	bs := NewBitSet(1000)

	// Set some bits
	bs.Set(0)
	bs.Set(100)
	bs.Set(999)

	// Verify they are set
	if !bs.Get(0) {
		t.Error("Bit 0 should be set")
	}
	if !bs.Get(100) {
		t.Error("Bit 100 should be set")
	}
	if !bs.Get(999) {
		t.Error("Bit 999 should be set")
	}

	// Verify unset bits are not set
	if bs.Get(1) {
		t.Error("Bit 1 should not be set")
	}
	if bs.Get(500) {
		t.Error("Bit 500 should not be set")
	}
}

func TestBitSet_Bounds(t *testing.T) {
	bs := NewBitSet(100)

	// Setting out of bounds should not panic
	bs.Set(1000)
	bs.Set(2000)

	// Getting out of bounds should return false
	if bs.Get(1000) {
		t.Error("Out of bounds get should return false")
	}
}

func BenchmarkBloomFilter_Add(b *testing.B) {
	bf, _ := NewBloomFilterBuilder().
		WithSize(10000).
		WithNumHashFunctions(5).
		Build()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.Add(string(rune(i)))
	}
}

func BenchmarkBloomFilter_MightContain(b *testing.B) {
	bf, _ := NewBloomFilterBuilder().
		WithSize(10000).
		WithNumHashFunctions(5).
		Build()

	// Pre-populate
	for i := 0; i < 1000; i++ {
		bf.Add(string(rune(i)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.MightContain(string(rune(i % 2000)))
	}
}

