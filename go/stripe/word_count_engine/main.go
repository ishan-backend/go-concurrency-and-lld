package main

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// WordCount represents a word and its count
type WordCount struct {
	Word  string
	Count int
}

// WordCountEngine counts words in a document and returns them sorted by frequency
type WordCountEngine struct {
	wordCounts map[string]int
}

// NewWordCountEngine creates a new word count engine
func NewWordCountEngine() *WordCountEngine {
	return &WordCountEngine{
		wordCounts: make(map[string]int),
	}
}

// ProcessDocument processes a document and counts words
// Words are case-insensitive and punctuation is removed
func (wce *WordCountEngine) ProcessDocument(document string) {
	// Convert to lowercase and split by whitespace
	words := strings.Fields(strings.ToLower(document))
	
	for _, word := range words {
		// Remove punctuation
		cleaned := strings.Map(func(r rune) rune {
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				return r
			}
			return -1 // Remove non-alphanumeric characters
		}, word)
		
		if cleaned != "" {
			wce.wordCounts[cleaned]++
		}
	}
}

// GetWordCounts returns word counts sorted by frequency (descending), then alphabetically
func (wce *WordCountEngine) GetWordCounts() []WordCount {
	var wordCounts []WordCount
	
	for word, count := range wce.wordCounts {
		wordCounts = append(wordCounts, WordCount{
			Word:  word,
			Count: count,
		})
	}
	
	// Sort by count (descending), then by word (ascending)
	sort.Slice(wordCounts, func(i, j int) bool {
		if wordCounts[i].Count != wordCounts[j].Count {
			return wordCounts[i].Count > wordCounts[j].Count
		}
		return wordCounts[i].Word < wordCounts[j].Word
	})
	
	return wordCounts
}

// GetTopWords returns top N words by frequency
func (wce *WordCountEngine) GetTopWords(n int) []WordCount {
	allCounts := wce.GetWordCounts()
	
	if n > len(allCounts) {
		n = len(allCounts)
	}
	
	return allCounts[:n]
}

// GetWordFrequency returns the frequency of a specific word
func (wce *WordCountEngine) GetWordFrequency(word string) int {
	return wce.wordCounts[strings.ToLower(word)]
}

// Reset clears all word counts
func (wce *WordCountEngine) Reset() {
	wce.wordCounts = make(map[string]int)
}

func main() {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Word Count Engine")
	fmt.Println(strings.Repeat("=", 60))
	
	// Example 1: Basic word counting
	fmt.Println("\nExample 1: Basic Word Counting")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	engine := NewWordCountEngine()
	document := "Practice makes perfect. you'll only get Perfect by practice. just practice!"
	engine.ProcessDocument(document)
	
	fmt.Printf("Document: %s\n\n", document)
	
	wordCounts := engine.GetWordCounts()
	fmt.Println("Word Counts:")
	for _, wc := range wordCounts {
		fmt.Printf("  %s: %d\n", wc.Word, wc.Count)
	}
	
	// Example 2: Top words
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 2: Top Words")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	topWords := engine.GetTopWords(3)
	fmt.Println("Top 3 words:")
	for i, wc := range topWords {
		fmt.Printf("  %d. %s: %d\n", i+1, wc.Word, wc.Count)
	}
	
	// Example 3: Word frequency lookup
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 3: Word Frequency Lookup")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	testWords := []string{"practice", "perfect", "makes", "nonexistent"}
	for _, word := range testWords {
		freq := engine.GetWordFrequency(word)
		fmt.Printf("  '%s': %d\n", word, freq)
	}
	
	// Example 4: Multiple documents
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 4: Multiple Documents")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	engine2 := NewWordCountEngine()
	documents := []string{
		"The quick brown fox jumps over the lazy dog",
		"The lazy dog was jumped over by the quick brown fox",
		"A quick brown fox is always faster than a lazy dog",
	}
	
	for i, doc := range documents {
		fmt.Printf("Processing document %d: %s\n", i+1, doc)
		engine2.ProcessDocument(doc)
	}
	
	fmt.Println("\nCombined Word Counts:")
	wordCounts2 := engine2.GetWordCounts()
	for _, wc := range wordCounts2 {
		fmt.Printf("  %s: %d\n", wc.Word, wc.Count)
	}
	
	// Example 5: Case-insensitive and punctuation handling
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 5: Case-Insensitive and Punctuation Handling")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	engine3 := NewWordCountEngine()
	doc3 := "Hello, World! HELLO world. hello-WORLD!"
	engine3.ProcessDocument(doc3)
	
	fmt.Printf("Document: %s\n", doc3)
	fmt.Println("\nWord Counts (case-insensitive, punctuation removed):")
	wordCounts3 := engine3.GetWordCounts()
	for _, wc := range wordCounts3 {
		fmt.Printf("  %s: %d\n", wc.Word, wc.Count)
	}
}

