package main

import (
	"fmt"
	"strings"
)

// Record represents a database record (map with string keys and integer values)
type Record map[string]int

// RecordComparator compares records based on a key and direction
type RecordComparator struct {
	key       string
	direction string // "asc" or "desc"
}

// NewRecordComparator creates a new comparator
func NewRecordComparator(key, direction string) *RecordComparator {
	return &RecordComparator{
		key:       key,
		direction: direction,
	}
}

// Compare compares two records
// Returns -1 if rec1 comes before rec2, 0 if equal, 1 if rec1 comes after rec2
func (rc *RecordComparator) Compare(rec1, rec2 Record) int {
	val1 := rc.getValue(rec1)
	val2 := rc.getValue(rec2)
	
	if val1 < val2 {
		if rc.direction == "asc" {
			return -1
		}
		return 1
	} else if val1 > val2 {
		if rc.direction == "asc" {
			return 1
		}
		return -1
	}
	return 0
}

// getValue gets the value for the key, defaulting to 0 if not present
func (rc *RecordComparator) getValue(rec Record) int {
	if val, exists := rec[rc.key]; exists {
		return val
	}
	return 0
}

// MinByKey returns the record with the minimum value for the specified key
// Records without the key are considered to have value 0
func MinByKey(key string, records []Record) Record {
	return FirstByKey(key, "asc", records)
}

// FirstByKey returns the first record according to key and direction
// direction must be "asc" (minimum) or "desc" (maximum)
// Records without the key are considered to have value 0
func FirstByKey(key, direction string, records []Record) Record {
	if len(records) == 0 {
		return nil
	}
	
	comparator := NewRecordComparator(key, direction)
	result := records[0]
	
	for i := 1; i < len(records); i++ {
		if comparator.Compare(records[i], result) < 0 {
			result = records[i]
		}
	}
	
	return result
}

func main() {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Analytical Database - Part 1: min_by_key")
	fmt.Println(strings.Repeat("=", 60))
	
	// Example 1: Basic min_by_key
	fmt.Println("\nExample 1: Basic min_by_key")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	records1 := []Record{
		{"a": 1, "b": 2},
		{"a": 2},
	}
	result1 := MinByKey("a", records1)
	fmt.Printf("Records: %v\n", records1)
	fmt.Printf("min_by_key(\"a\"): %v\n", result1)
	
	// Example 2: Different order
	fmt.Println("\nExample 2: Different order")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	records2 := []Record{
		{"a": 2},
		{"a": 1, "b": 2},
	}
	result2 := MinByKey("a", records2)
	fmt.Printf("Records: %v\n", records2)
	fmt.Printf("min_by_key(\"a\"): %v\n", result2)
	
	// Example 3: Missing key (treated as 0)
	fmt.Println("\nExample 3: Missing key (treated as 0)")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	records3 := []Record{
		{"a": 1, "b": 2},
		{"a": 2},
	}
	result3 := MinByKey("b", records3)
	fmt.Printf("Records: %v\n", records3)
	fmt.Printf("min_by_key(\"b\"): %v (missing key treated as 0)\n", result3)
	
	// Example 4: Empty record
	fmt.Println("\nExample 4: Empty record")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	records4 := []Record{{}}
	result4 := MinByKey("a", records4)
	fmt.Printf("Records: %v\n", records4)
	fmt.Printf("min_by_key(\"a\"): %v\n", result4)
	
	// Example 5: Negative values
	fmt.Println("\nExample 5: Negative values")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	records5 := []Record{
		{"a": -1},
		{"b": -1},
	}
	result5 := MinByKey("b", records5)
	fmt.Printf("Records: %v\n", records5)
	fmt.Printf("min_by_key(\"b\"): %v\n", result5)
	
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Analytical Database - Part 2: first_by_key")
	fmt.Println(strings.Repeat("=", 60))
	
	// Example 6: first_by_key ascending
	fmt.Println("\nExample 6: first_by_key ascending")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	records6 := []Record{{"a": 1}}
	result6 := FirstByKey("a", "asc", records6)
	fmt.Printf("Records: %v\n", records6)
	fmt.Printf("first_by_key(\"a\", \"asc\"): %v\n", result6)
	
	// Example 7: first_by_key with missing key
	fmt.Println("\nExample 7: first_by_key with missing key")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	records7 := []Record{
		{"b": 1},
		{"b": -2},
		{"a": 10},
	}
	result7 := FirstByKey("a", "asc", records7)
	fmt.Printf("Records: %v\n", records7)
	fmt.Printf("first_by_key(\"a\", \"asc\"): %v (missing key = 0, so b:1 or b:-2)\n", result7)
	
	// Example 8: first_by_key descending
	fmt.Println("\nExample 8: first_by_key descending")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	records8 := []Record{
		{"b": 1},
		{"b": -2},
		{"a": 10},
	}
	result8 := FirstByKey("a", "desc", records8)
	fmt.Printf("Records: %v\n", records8)
	fmt.Printf("first_by_key(\"a\", \"desc\"): %v\n", result8)
	
	// Example 9: first_by_key on different key
	fmt.Println("\nExample 9: first_by_key on different key")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	records9 := []Record{
		{"b": 1},
		{"b": -2},
		{"a": 10},
	}
	result9a := FirstByKey("b", "asc", records9)
	result9b := FirstByKey("b", "desc", records9)
	fmt.Printf("Records: %v\n", records9)
	fmt.Printf("first_by_key(\"b\", \"asc\"): %v\n", result9a)
	fmt.Printf("first_by_key(\"b\", \"desc\"): %v\n", result9b)
	
	// Example 10: first_by_key with empty records
	fmt.Println("\nExample 10: first_by_key with empty records")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	records10 := []Record{
		{},
		{"a": 10, "b": -10},
		{},
		{"a": 3, "c": 3},
	}
	result10 := FirstByKey("a", "desc", records10)
	fmt.Printf("Records: %v\n", records10)
	fmt.Printf("first_by_key(\"a\", \"desc\"): %v\n", result10)
	
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Analytical Database - Part 3: RecordComparator")
	fmt.Println(strings.Repeat("=", 60))
	
	// Example 11: Comparator usage
	fmt.Println("\nExample 11: Comparator usage")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	cmp := NewRecordComparator("a", "asc")
	rec1 := Record{"a": 1}
	rec2 := Record{"a": 2}
	rec3 := Record{"a": 1}
	
	fmt.Printf("Compare {%v} vs {%v}: %d\n", rec1, rec2, cmp.Compare(rec1, rec2))
	fmt.Printf("Compare {%v} vs {%v}: %d\n", rec2, rec1, cmp.Compare(rec2, rec1))
	fmt.Printf("Compare {%v} vs {%v}: %d\n", rec1, rec3, cmp.Compare(rec1, rec3))
	
	// Example 12: Comparator with descending order
	fmt.Println("\nExample 12: Comparator with descending order")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	cmpDesc := NewRecordComparator("a", "desc")
	fmt.Printf("Compare {%v} vs {%v} (desc): %d\n", rec1, rec2, cmpDesc.Compare(rec1, rec2))
	fmt.Printf("Compare {%v} vs {%v} (desc): %d\n", rec2, rec1, cmpDesc.Compare(rec2, rec1))
}

