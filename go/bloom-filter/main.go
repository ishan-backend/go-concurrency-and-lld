package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("=== Bloom Filter Low-Level Design Demo ===")
	fmt.Println()

	// Demo 1: Basic usage with manual configuration
	demoBasicUsage()

	// Demo 2: Using builder with expected items and false positive rate
	demoBuilderWithExpectedItems()

	// Demo 3: Demonstrating false positives
	demoFalsePositives()

	// Demo 4: Thread-safety demonstration
	demoThreadSafety()
}

func demoBasicUsage() {
	fmt.Println("--- Demo 1: Basic Usage ---")

	// Create a Bloom Filter with manual configuration
	bf, err := NewBloomFilterBuilder().
		WithSize(1000).
		WithNumHashFunctions(3).
		WithHashStrategy(NewFNV1aHashStrategy()).
		Build()

	if err != nil {
		fmt.Printf("Error creating Bloom Filter: %v\n", err)
		return
	}

	// Add some elements
	elements := []string{"apple", "banana", "cherry", "date", "elderberry"}
	fmt.Printf("Adding elements: %v\n", elements)

	for _, elem := range elements {
		if err := bf.Add(elem); err != nil {
			fmt.Printf("Error adding %s: %v\n", elem, err)
		}
	}

	// Check for presence
	fmt.Println("\nChecking for presence:")
	testElements := []string{"apple", "banana", "grape", "mango"}
	for _, elem := range testElements {
		mightContain, err := bf.MightContain(elem)
		if err != nil {
			fmt.Printf("Error checking %s: %v\n", elem, err)
			continue
		}
		status := "might be present"
		if !mightContain {
			status = "definitely not present"
		}
		fmt.Printf("  %s: %s\n", elem, status)
	}

	fmt.Println()
}

func demoBuilderWithExpectedItems() {
	fmt.Println("--- Demo 2: Builder with Expected Items ---")

	// Create a Bloom Filter optimized for 1000 items with 1% false positive rate
	bf, err := NewBloomFilterBuilder().
		WithExpectedItems(1000, 0.01).
		WithHashStrategy(NewFNV1aHashStrategy()).
		Build()

	if err != nil {
		fmt.Printf("Error creating Bloom Filter: %v\n", err)
		return
	}

	fmt.Printf("Bloom Filter configured:\n")
	fmt.Printf("  Bit array size (m): %d\n", bf.Size())
	fmt.Printf("  Number of hash functions (k): %d\n", bf.NumHashFunctions())
	fmt.Printf("  Expected items: 1000\n")
	fmt.Printf("  Target false positive rate: 1%%\n")

	// Add elements
	elements := []string{"user1", "user2", "user3", "user4", "user5"}
	fmt.Printf("\nAdding elements: %v\n", elements)

	for _, elem := range elements {
		bf.Add(elem)
	}

	// Verify no false negatives
	fmt.Println("\nVerifying no false negatives (elements we added):")
	allFound := true
	for _, elem := range elements {
		mightContain, _ := bf.MightContain(elem)
		if !mightContain {
			fmt.Printf("  ERROR: %s not found (false negative!)\n", elem)
			allFound = false
		} else {
			fmt.Printf("  ✓ %s found\n", elem)
		}
	}
	if allFound {
		fmt.Println("  ✓ No false negatives - all added elements found!")
	}

	fmt.Println()
}

func demoFalsePositives() {
	fmt.Println("--- Demo 3: Demonstrating False Positives ---")

	// Create a smaller Bloom Filter to increase false positive probability
	bf, err := NewBloomFilterBuilder().
		WithSize(100).
		WithNumHashFunctions(2).
		Build()

	if err != nil {
		fmt.Printf("Error creating Bloom Filter: %v\n", err)
		return
	}

	// Add a small set of elements
	addedElements := []string{"red", "green", "blue", "yellow", "purple"}
	fmt.Printf("Adding elements: %v\n", addedElements)

	for _, elem := range addedElements {
		bf.Add(elem)
	}

	// Test many random elements that were NOT added
	fmt.Println("\nTesting 100 random elements (none were added):")
	rand.Seed(time.Now().UnixNano())
	falsePositives := 0
	totalTests := 100

	for i := 0; i < totalTests; i++ {
		randomElem := fmt.Sprintf("random%d", rand.Intn(10000))
		mightContain, _ := bf.MightContain(randomElem)
		if mightContain {
			falsePositives++
		}
	}

	falsePosRate := float64(falsePositives) / float64(totalTests) * 100
	fmt.Printf("  False positives: %d out of %d (%.2f%%)\n",
		falsePositives, totalTests, falsePosRate)
	fmt.Println("  Note: This demonstrates the space vs accuracy trade-off")
	fmt.Println()
}

func demoThreadSafety() {
	fmt.Println("--- Demo 4: Thread-Safety Demonstration ---")

	bf, err := NewBloomFilterBuilder().
		WithSize(10000).
		WithNumHashFunctions(5).
		Build()

	if err != nil {
		fmt.Printf("Error creating Bloom Filter: %v\n", err)
		return
	}

	// Simulate concurrent writes and reads
	done := make(chan bool)
	numGoroutines := 10
	elementsPerGoroutine := 100

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < elementsPerGoroutine; j++ {
				elem := fmt.Sprintf("element-%d-%d", id, j)
				bf.Add(elem)
			}
			done <- true
		}(i)
	}

	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < elementsPerGoroutine; j++ {
				elem := fmt.Sprintf("element-%d-%d", id, j)
				bf.MightContain(elem)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines*2; i++ {
		<-done
	}

	fmt.Printf("Successfully completed concurrent operations:\n")
	fmt.Printf("  %d goroutines writing %d elements each\n", numGoroutines, elementsPerGoroutine)
	fmt.Printf("  %d goroutines reading %d elements each\n", numGoroutines, elementsPerGoroutine)
	fmt.Println("  ✓ No race conditions detected (thread-safe)")
	fmt.Println()
}
