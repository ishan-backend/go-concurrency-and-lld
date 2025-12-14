package main

import (
	"fmt"
	"sort"
	"strings"
)

// Pair represents a pair of indices
type Pair struct {
	First  int
	Second int
}

// FindPairs finds all pairs of indices whose values sum to target
type FindPairs struct {
	numbers []int
}

// NewFindPairs creates a new FindPairs instance
func NewFindPairs(numbers []int) *FindPairs {
	return &FindPairs{
		numbers: numbers,
	}
}

// FindPairsWithSum finds all pairs of indices whose values sum to target
// Returns list of pairs (i, j) where i < j and numbers[i] + numbers[j] == target
func (fp *FindPairs) FindPairsWithSum(target int) []Pair {
	var pairs []Pair
	seen := make(map[int][]int) // value -> list of indices
	
	// First pass: build index map
	for i, num := range fp.numbers {
		seen[num] = append(seen[num], i)
	}
	
	// Second pass: find pairs
	usedPairs := make(map[string]bool)
	
	for i, num := range fp.numbers {
		complement := target - num
		
		if indices, exists := seen[complement]; exists {
			for _, j := range indices {
				if i < j { // Ensure i < j to avoid duplicates
					key := fmt.Sprintf("%d,%d", i, j)
					if !usedPairs[key] {
						pairs = append(pairs, Pair{First: i, Second: j})
						usedPairs[key] = true
					}
				}
			}
		}
	}
	
	return pairs
}

// FindPairsWithSumSorted returns pairs sorted by first index, then second index
func (fp *FindPairs) FindPairsWithSumSorted(target int) []Pair {
	pairs := fp.FindPairsWithSum(target)
	
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].First != pairs[j].First {
			return pairs[i].First < pairs[j].First
		}
		return pairs[i].Second < pairs[j].Second
	})
	
	return pairs
}

// CountPairs counts the number of pairs that sum to target
func (fp *FindPairs) CountPairs(target int) int {
	return len(fp.FindPairsWithSum(target))
}

// FindPairsWithSumTwoPointer uses two-pointer technique (requires sorted array)
func (fp *FindPairs) FindPairsWithSumTwoPointer(target int) []Pair {
	// Create a copy and sort with indices
	type indexedNum struct {
		value int
		index int
	}
	
	indexed := make([]indexedNum, len(fp.numbers))
	for i, num := range fp.numbers {
		indexed[i] = indexedNum{value: num, index: i}
	}
	
	sort.Slice(indexed, func(i, j int) bool {
		return indexed[i].value < indexed[j].value
	})
	
	var pairs []Pair
	left := 0
	right := len(indexed) - 1
	
	for left < right {
		sum := indexed[left].value + indexed[right].value
		
		if sum == target {
			// Add pair (ensure smaller index first)
			firstIdx := indexed[left].index
			secondIdx := indexed[right].index
			if firstIdx < secondIdx {
				pairs = append(pairs, Pair{First: firstIdx, Second: secondIdx})
			} else {
				pairs = append(pairs, Pair{First: secondIdx, Second: firstIdx})
			}
			
			// Move pointers to find more pairs
			left++
			right--
			
			// Skip duplicates
			for left < right && indexed[left].value == indexed[left-1].value {
				left++
			}
			for left < right && indexed[right].value == indexed[right+1].value {
				right--
			}
		} else if sum < target {
			left++
		} else {
			right--
		}
	}
	
	return pairs
}

func main() {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Find Pairs - Two Sum Problem")
	fmt.Println(strings.Repeat("=", 60))
	
	// Example 1: Basic two sum
	fmt.Println("\nExample 1: Basic Two Sum")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	numbers1 := []int{2, 7, 11, 15}
	target1 := 9
	fp1 := NewFindPairs(numbers1)
	
	fmt.Printf("Numbers: %v\n", numbers1)
	fmt.Printf("Target: %d\n\n", target1)
	
	pairs1 := fp1.FindPairsWithSumSorted(target1)
	fmt.Println("Pairs (indices):")
	for _, pair := range pairs1 {
		fmt.Printf("  (%d, %d) -> values: [%d, %d] = %d\n",
			pair.First, pair.Second,
			numbers1[pair.First], numbers1[pair.Second],
			numbers1[pair.First]+numbers1[pair.Second])
	}
	
	// Example 2: Multiple pairs
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 2: Multiple Pairs")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	numbers2 := []int{1, 2, 3, 4, 5, 6}
	target2 := 7
	fp2 := NewFindPairs(numbers2)
	
	fmt.Printf("Numbers: %v\n", numbers2)
	fmt.Printf("Target: %d\n\n", target2)
	
	pairs2 := fp2.FindPairsWithSumSorted(target2)
	fmt.Println("Pairs (indices):")
	for _, pair := range pairs2 {
		fmt.Printf("  (%d, %d) -> values: [%d, %d] = %d\n",
			pair.First, pair.Second,
			numbers2[pair.First], numbers2[pair.Second],
			numbers2[pair.First]+numbers2[pair.Second])
	}
	
	fmt.Printf("\nTotal pairs: %d\n", fp2.CountPairs(target2))
	
	// Example 3: Duplicate values
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 3: Duplicate Values")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	numbers3 := []int{3, 3, 4, 5, 6, 3}
	target3 := 6
	fp3 := NewFindPairs(numbers3)
	
	fmt.Printf("Numbers: %v\n", numbers3)
	fmt.Printf("Target: %d\n\n", target3)
	
	pairs3 := fp3.FindPairsWithSumSorted(target3)
	fmt.Println("Pairs (indices):")
	for _, pair := range pairs3 {
		fmt.Printf("  (%d, %d) -> values: [%d, %d] = %d\n",
			pair.First, pair.Second,
			numbers3[pair.First], numbers3[pair.Second],
			numbers3[pair.First]+numbers3[pair.Second])
	}
	
	// Example 4: Two-pointer approach
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 4: Two-Pointer Approach")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	numbers4 := []int{1, 3, 4, 5, 6, 8, 9}
	target4 := 10
	fp4 := NewFindPairs(numbers4)
	
	fmt.Printf("Numbers: %v\n", numbers4)
	fmt.Printf("Target: %d\n\n", target4)
	
	pairs4 := fp4.FindPairsWithSumTwoPointer(target4)
	fmt.Println("Pairs using two-pointer (indices):")
	for _, pair := range pairs4 {
		fmt.Printf("  (%d, %d) -> values: [%d, %d] = %d\n",
			pair.First, pair.Second,
			numbers4[pair.First], numbers4[pair.Second],
			numbers4[pair.First]+numbers4[pair.Second])
	}
	
	// Example 5: No pairs found
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 5: No Pairs Found")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	numbers5 := []int{1, 2, 3, 4}
	target5 := 10
	fp5 := NewFindPairs(numbers5)
	
	fmt.Printf("Numbers: %v\n", numbers5)
	fmt.Printf("Target: %d\n\n", target5)
	
	pairs5 := fp5.FindPairsWithSum(target5)
	if len(pairs5) == 0 {
		fmt.Println("No pairs found")
	} else {
		fmt.Println("Pairs found:")
		for _, pair := range pairs5 {
			fmt.Printf("  (%d, %d)\n", pair.First, pair.Second)
		}
	}
}

