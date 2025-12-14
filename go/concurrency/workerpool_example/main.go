package main

import (
	"concurrency/workerpool"
	"fmt"
)

func main() {
	// Example 1: Compare strings within a single set to find duplicates
	// fmt.Println("=== Example 1: Finding matching strings in a single set ===")
	// strings := []string{
	// 	"hello",
	// 	"world",
	// 	"hello",
	// 	"go",
	// 	"world",
	// 	"hello",
	// 	"python",
	// 	"go",
	// }

	// // Create comparator with limited workers (e.g., 4 workers for millions of strings)
	// comparator := workerpool.NewStringComparator(4, strings)

	// // Generate all pairwise comparison tasks
	// comparator.GenerateComparisonTasks()

	// // Wait for all comparisons to complete
	// comparator.WaitForCompletion()

	// // Print results
	// comparator.PrintResults()

	// Example 2: Compare two files (represented as string arrays)
	fmt.Println("\n=== Example 2: Comparing two files ===")
	file1Strings := []string{
		"file1_line1",
		"file1_line2",
		"file1_line3",
		"common_line",
	}

	file2Strings := []string{
		"file2_line1",
		"common_line",
		"file2_line3",
		"common_line",
	}

	// Create comparator for file comparison
	fileComparator := workerpool.NewStringComparator(4, nil)

	// Generate comparison tasks between two files
	fileComparator.GenerateComparisonTasksForFiles(file1Strings, file2Strings)

	// Wait for completion
	fileComparator.WaitForCompletion()

	// Print results
	fileComparator.PrintResults()

	// Example 3: Access results programmatically
	fmt.Println("\n=== Example 3: Accessing results programmatically ===")
	results := fileComparator.GetResults()
	matchingSets := results.GetMatchingSets()

	for str, indices := range matchingSets {
		fmt.Printf("String %q appears at indices: %v\n", str, indices)
	}
}
