package main

import (
	"fmt"
	"strings"
)

// compressWord compresses a single word using run-length encoding
// e.g., "stripe" -> "s4e" (s + 4 middle chars + e)
func compressWord(word string) string {
	if len(word) <= 2 {
		return word
	}
	// Keep first and last char, replace middle with count
	middle := len(word) - 2
	return string(word[0]) + fmt.Sprintf("%d", middle) + string(word[len(word)-1])
}

// compressMinorParts compresses minor parts (separated by '.')
// minor_parts indicates how many parts from the end to compress
func compressMinorParts(major string, minorParts int) string {
	parts := strings.Split(major, ".")

	// Determine starting index for compression
	// We compress from the end, minorParts at a time
	startIdx := len(parts) - minorParts
	if startIdx < 0 {
		startIdx = 0
	}

	// Compress parts from startIdx onwards
	for i := startIdx; i < len(parts); i++ {
		parts[i] = compressWord(parts[i])
	}

	return strings.Join(parts, ".")
}

// compressMajorParts compresses major parts (separated by '/')
// minor_parts indicates how many minor parts to compress in each major part
func compressString(str string, minorParts int) string {
	majorParts := strings.Split(str, "/")

	for i := range majorParts {
		majorParts[i] = compressMinorParts(majorParts[i], minorParts)
	}

	return strings.Join(majorParts, "/")
}

// performCompression performs compression in two phases
func performCompression(str string, minorParts int) string {
	// Phase 1: Compress the first major part
	fmt.Printf("Original string: %s\n", str)
	fmt.Printf("Minor parts to compress: %d\n\n", minorParts)

	phase1 := compressString(str, minorParts)
	fmt.Printf("After Phase 1 (compress first major part): %s\n", phase1)

	// Phase 2: Compress all major parts
	majorParts := strings.Split(phase1, "/")
	for i := range majorParts {
		majorParts[i] = compressMinorParts(majorParts[i], minorParts)
	}
	phase2 := strings.Join(majorParts, "/")
	fmt.Printf("After Phase 2 (compress all major parts): %s\n\n", phase2)

	return phase2
}

func main() {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Example 1:")
	fmt.Println(strings.Repeat("=", 60))
	str1 := "stripe.com/payments/checkout/customer.john.doe"
	result1 := performCompression(str1, 2)
	fmt.Printf("Final result: %s\n\n", result1)

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Example 2:")
	fmt.Println(strings.Repeat("=", 60))
	str2 := "www.api.stripe.com/checkout"
	result2 := performCompression(str2, 3)
	fmt.Printf("Final result: %s\n\n", result2)

	// Additional test case
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Test Case 3 (single char words):")
	fmt.Println(strings.Repeat("=", 60))
	str3 := "a.b.c/d.e.f"
	result3 := performCompression(str3, 2)
	fmt.Printf("Final result: %s\n", result3)
}
