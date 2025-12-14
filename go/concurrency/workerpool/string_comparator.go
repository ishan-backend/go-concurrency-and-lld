package workerpool

import (
	"fmt"
	"sync"
)

// StringComparisonTask represents a task to compare two strings
type StringComparisonTask struct {
	str1      string
	str2      string
	index1    int
	index2    int
	resultMap *ResultCollector
}

// Run executes the string comparison
func (s *StringComparisonTask) Run() {
	if s.str1 == s.str2 {
		s.resultMap.AddMatch(s.index1, s.index2, s.str1)
	}
}

// ResultCollector collects and organizes matching strings
type ResultCollector struct {
	mu            sync.RWMutex
	matches       map[string][]int // string -> list of indices that have this string
	indexToString map[int]string   // index -> string value
}

// NewResultCollector creates a new result collector
func NewResultCollector() *ResultCollector {
	return &ResultCollector{
		matches:       make(map[string][]int),
		indexToString: make(map[int]string),
	}
}

// AddMatch records that two indices have matching strings
func (r *ResultCollector) AddMatch(index1, index2 int, str string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Store the string for each index
	r.indexToString[index1] = str
	r.indexToString[index2] = str

	// Add indices to the matches set for this string
	if _, exists := r.matches[str]; !exists {
		r.matches[str] = make([]int, 0)
	}

	// Use a map to track which indices we've already added (for O(1) lookup)
	indexSet := make(map[int]bool)
	for _, idx := range r.matches[str] {
		indexSet[idx] = true
	}

	// Add index1 if not already present
	if !indexSet[index1] {
		r.matches[str] = append(r.matches[str], index1)
		indexSet[index1] = true
	}

	// Add index2 if not already present
	if !indexSet[index2] {
		r.matches[str] = append(r.matches[str], index2)
	}
}

// GetMatchingSets returns all sets of matching strings
// Each set contains all indices that have the same string value
func (r *ResultCollector) GetMatchingSets() map[string][]int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string][]int)
	for str, indices := range r.matches {
		// Only include sets with 2 or more matches
		if len(indices) >= 2 {
			result[str] = make([]int, len(indices))
			copy(result[str], indices)
		}
	}
	return result
}

// PrintResults prints all matching sets in a readable format
func (r *ResultCollector) PrintResults() {
	sets := r.GetMatchingSets()
	fmt.Printf("\n=== Matching String Sets ===\n")
	fmt.Printf("Total unique matching strings: %d\n\n", len(sets))

	for str, indices := range sets {
		fmt.Printf("String: %q\n", str)
		fmt.Printf("  Found at indices: %v\n", indices)
		fmt.Printf("  Count: %d\n\n", len(indices))
	}
}

// StringComparator manages string comparison using the worker pool
type StringComparator struct {
	pool        *WorkerPool
	collector   *ResultCollector
	strings     []string
	comparisons int
}

// NewStringComparator creates a new string comparator with default buffer size (1000)
func NewStringComparator(workerCount int, strings []string) *StringComparator {
	return NewStringComparatorWithBuffer(workerCount, 1000, strings)
}

// NewStringComparatorWithBuffer creates a new string comparator with custom buffer size.
// For millions of string comparisons, use a larger buffer (e.g., 5000-10000) to improve throughput.
func NewStringComparatorWithBuffer(workerCount int, bufferSize int, strings []string) *StringComparator {
	return &StringComparator{
		pool:      NewWorkerPoolWithBuffer(workerCount, bufferSize),
		collector: NewResultCollector(),
		strings:   strings,
	}
}

// GenerateComparisonTasks generates all pairwise comparison tasks
// This creates tasks for comparing every string with every other string
func (sc *StringComparator) GenerateComparisonTasks() {
	count := 0
	for i := 0; i < len(sc.strings); i++ {
		for j := i + 1; j < len(sc.strings); j++ {
			task := &StringComparisonTask{
				str1:      sc.strings[i],
				str2:      sc.strings[j],
				index1:    i,
				index2:    j,
				resultMap: sc.collector,
			}
			sc.pool.Add(task)
			count++
		}
	}
	sc.comparisons = count
	fmt.Printf("Generated %d comparison tasks\n", count)
}

// GenerateComparisonTasksForFiles generates comparison tasks for two file string arrays
// This compares each string from file1 with each string from file2
func (sc *StringComparator) GenerateComparisonTasksForFiles(file1Strings []string, file2Strings []string) {
	count := 0
	for i, str1 := range file1Strings {
		for j, str2 := range file2Strings {
			task := &StringComparisonTask{
				str1:      str1,
				str2:      str2,
				index1:    i,
				index2:    len(file1Strings) + j, // Offset indices for file2
				resultMap: sc.collector,
			}
			sc.pool.Add(task)
			count++
		}
	}
	sc.comparisons = count
	fmt.Printf("Generated %d comparison tasks (file1: %d strings, file2: %d strings)\n",
		count, len(file1Strings), len(file2Strings))
}

// WaitForCompletion waits for all comparison tasks to complete
func (sc *StringComparator) WaitForCompletion() {
	sc.pool.Wait()
	fmt.Printf("All %d comparisons completed\n", sc.comparisons)
}

// GetResults returns the result collector
func (sc *StringComparator) GetResults() *ResultCollector {
	return sc.collector
}

// PrintResults prints all matching sets
func (sc *StringComparator) PrintResults() {
	sc.collector.PrintResults()
}
