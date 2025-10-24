package main

import (
	"fmt"
	"sort"
	"strings"
)

// TrieNode represents a node in the Trie
type TrieNode struct {
	children    map[rune]*TrieNode // Child nodes mapped by character
	isEndOfWord bool               // Marks end of a valid word
	word        string             // Store complete word at leaf nodes
	frequency   int                // Word frequency/popularity score
}

// NewTrieNode creates a new TrieNode
func NewTrieNode() *TrieNode {
	return &TrieNode{
		children:    make(map[rune]*TrieNode),
		isEndOfWord: false,
		word:        "",
		frequency:   0,
	}
}

// Trie represents the Trie data structure for autocomplete
type Trie struct {
	root           *TrieNode
	maxSuggestions int // Maximum number of suggestions to return
}

// NewTrie creates a new Trie with default max suggestions
func NewTrie() *Trie {
	return &Trie{
		root:           NewTrieNode(),
		maxSuggestions: 10,
	}
}

// NewTrieWithLimit creates a new Trie with custom suggestion limit
func NewTrieWithLimit(limit int) *Trie {
	return &Trie{
		root:           NewTrieNode(),
		maxSuggestions: limit,
	}
}

// Insert adds a word to the Trie with optional frequency
// Time Complexity: O(m) where m is length of word
func (t *Trie) Insert(word string, frequency int) {
	if word == "" {
		return
	}

	word = strings.ToLower(word)
	node := t.root

	// Traverse/create path for each character
	for _, char := range word {
		if _, exists := node.children[char]; !exists {
			node.children[char] = NewTrieNode()
		}
		node = node.children[char]
	}

	// Mark end of word and store metadata
	node.isEndOfWord = true
	node.word = word
	node.frequency = frequency
}

// Search checks if a word exists in the Trie
// Time Complexity: O(m)
func (t *Trie) Search(word string) bool {
	word = strings.ToLower(word)
	node := t.root

	for _, char := range word {
		if _, exists := node.children[char]; !exists {
			return false
		}
		node = node.children[char]
	}

	return node.isEndOfWord
}

// StartsWith checks if any word starts with the given prefix
// Time Complexity: O(m)
func (t *Trie) StartsWith(prefix string) bool {
	prefix = strings.ToLower(prefix)
	node := t.root

	for _, char := range prefix {
		if _, exists := node.children[char]; !exists {
			return false
		}
		node = node.children[char]
	}

	return true
}

// Suggestion represents an autocomplete suggestion
type Suggestion struct {
	Word      string
	Frequency int
}

// AutoComplete returns suggestions for a given prefix
// Time Complexity: O(m + n) where m is prefix length, n is number of nodes in subtree
func (t *Trie) AutoComplete(prefix string) []Suggestion {
	prefix = strings.ToLower(prefix)
	node := t.root

	// Navigate to the prefix node
	for _, char := range prefix {
		if _, exists := node.children[char]; !exists {
			return []Suggestion{} // No suggestions if prefix doesn't exist
		}
		node = node.children[char]
	}

	// Collect all words with this prefix using DFS
	suggestions := []Suggestion{}
	t.collectSuggestions(node, &suggestions)

	// Sort by frequency (descending) and then alphabetically
	sort.Slice(suggestions, func(i, j int) bool {
		if suggestions[i].Frequency == suggestions[j].Frequency {
			return suggestions[i].Word < suggestions[j].Word
		}
		return suggestions[i].Frequency > suggestions[j].Frequency
	})

	// Limit number of suggestions
	if len(suggestions) > t.maxSuggestions {
		suggestions = suggestions[:t.maxSuggestions]
	}

	return suggestions
}

// collectSuggestions performs DFS to collect all words in subtree
func (t *Trie) collectSuggestions(node *TrieNode, suggestions *[]Suggestion) {
	if node == nil {
		return
	}

	// If this node marks end of word, add to suggestions
	if node.isEndOfWord {
		*suggestions = append(*suggestions, Suggestion{
			Word:      node.word,
			Frequency: node.frequency,
		})
	}

	// Recursively collect from all children
	for _, child := range node.children {
		t.collectSuggestions(child, suggestions)
	}
}

// AutoCompleteTopK returns top K suggestions for a given prefix
func (t *Trie) AutoCompleteTopK(prefix string, k int) []Suggestion {
	oldLimit := t.maxSuggestions
	t.maxSuggestions = k
	suggestions := t.AutoComplete(prefix)
	t.maxSuggestions = oldLimit
	return suggestions
}

// Delete removes a word from the Trie
// Time Complexity: O(m)
func (t *Trie) Delete(word string) bool {
	word = strings.ToLower(word)
	return t.deleteHelper(t.root, word, 0)
}

// deleteHelper recursively deletes a word
func (t *Trie) deleteHelper(node *TrieNode, word string, index int) bool {
	if node == nil {
		return false
	}

	// Reached end of word
	if index == len(word) {
		if !node.isEndOfWord {
			return false // Word doesn't exist
		}
		node.isEndOfWord = false
		node.word = ""
		node.frequency = 0
		// Return true if node has no children (can be deleted)
		return len(node.children) == 0
	}

	char := rune(word[index])
	childNode, exists := node.children[char]
	if !exists {
		return false // Word doesn't exist
	}

	shouldDeleteChild := t.deleteHelper(childNode, word, index+1)

	// Delete child node if it's safe to delete
	if shouldDeleteChild {
		delete(node.children, char)
		// Return true if current node has no children and is not end of another word
		return len(node.children) == 0 && !node.isEndOfWord
	}

	return false
}

// UpdateFrequency updates the frequency of a word
func (t *Trie) UpdateFrequency(word string, newFrequency int) bool {
	word = strings.ToLower(word)
	node := t.root

	for _, char := range word {
		if _, exists := node.children[char]; !exists {
			return false
		}
		node = node.children[char]
	}

	if !node.isEndOfWord {
		return false
	}

	node.frequency = newFrequency
	return true
}

// GetWordCount returns total number of words in the Trie
func (t *Trie) GetWordCount() int {
	return t.countWords(t.root)
}

// countWords recursively counts words in the Trie
func (t *Trie) countWords(node *TrieNode) int {
	if node == nil {
		return 0
	}

	count := 0
	if node.isEndOfWord {
		count = 1
	}

	for _, child := range node.children {
		count += t.countWords(child)
	}

	return count
}

// Display prints the Trie structure (for debugging)
func (t *Trie) Display() {
	fmt.Println("Trie Structure:")
	t.displayHelper(t.root, "", 0)
}

// displayHelper recursively displays the Trie
func (t *Trie) displayHelper(node *TrieNode, prefix string, level int) {
	if node.isEndOfWord {
		fmt.Printf("%s└── [%s] (freq: %d)\n",
			strings.Repeat("  ", level), node.word, node.frequency)
	}

	for char, child := range node.children {
		fmt.Printf("%s├── %c\n", strings.Repeat("  ", level), char)
		t.displayHelper(child, prefix+string(char), level+1)
	}
}

// AutocompleteSystem demonstrates a complete autocomplete system
type AutocompleteSystem struct {
	trie    *Trie
	history map[string]int // Track search history
}

// NewAutocompleteSystem creates a new autocomplete system
func NewAutocompleteSystem() *AutocompleteSystem {
	return &AutocompleteSystem{
		trie:    NewTrie(),
		history: make(map[string]int),
	}
}

// AddWord adds a word to the system
func (as *AutocompleteSystem) AddWord(word string, frequency int) {
	as.trie.Insert(word, frequency)
}

// GetSuggestions returns suggestions and updates history
func (as *AutocompleteSystem) GetSuggestions(prefix string) []string {
	suggestions := as.trie.AutoComplete(prefix)

	// Update history for this prefix
	as.history[prefix]++

	// Convert to string array
	result := make([]string, len(suggestions))
	for i, s := range suggestions {
		result[i] = fmt.Sprintf("%s (score: %d)", s.Word, s.Frequency)
	}
	return result
}

// GetSearchHistory returns search history
func (as *AutocompleteSystem) GetSearchHistory() map[string]int {
	return as.history
}

func main() {
	fmt.Println("=== Trie-based Autocomplete System ===\n")

	// Create autocomplete system
	system := NewAutocompleteSystem()

	// Sample dataset with frequencies
	words := map[string]int{
		"apple":       100,
		"application": 80,
		"apply":       60,
		"app":         150,
		"appetite":    40,
		"banana":      90,
		"band":        70,
		"bandana":     30,
		"can":         120,
		"canada":      95,
		"canal":       50,
		"candy":       85,
	}

	// Insert words
	fmt.Println("Inserting words into Trie...")
	for word, freq := range words {
		system.AddWord(word, freq)
	}

	// Test autocomplete with different prefixes
	testPrefixes := []string{"app", "ban", "can", "appl", "cana"}

	fmt.Println("\n=== Autocomplete Suggestions ===")
	for _, prefix := range testPrefixes {
		suggestions := system.GetSuggestions(prefix)
		fmt.Printf("\nPrefix: '%s'\n", prefix)
		fmt.Println("Suggestions:")
		for i, s := range suggestions {
			fmt.Printf("  %d. %s\n", i+1, s)
		}
	}

	// Test search functionality
	fmt.Println("\n=== Search Tests ===")
	testWords := []string{"app", "application", "apples", "banana"}
	for _, word := range testWords {
		exists := system.trie.Search(word)
		fmt.Printf("Search '%s': %v\n", word, exists)
	}

	// Test prefix checking
	fmt.Println("\n=== Prefix Tests ===")
	testPrefixesCheck := []string{"app", "ban", "xyz", "cand"}
	for _, prefix := range testPrefixesCheck {
		exists := system.trie.StartsWith(prefix)
		fmt.Printf("Prefix '%s' exists: %v\n", prefix, exists)
	}

	// Update frequency
	fmt.Println("\n=== Update Frequency Test ===")
	fmt.Println("Updating 'apple' frequency to 200")
	system.trie.UpdateFrequency("apple", 200)
	suggestions := system.GetSuggestions("app")
	fmt.Println("New suggestions for 'app':")
	for i, s := range suggestions {
		fmt.Printf("  %d. %s\n", i+1, s)
	}

	// Delete word
	fmt.Println("\n=== Delete Test ===")
	fmt.Println("Deleting 'appetite'")
	system.trie.Delete("appetite")
	suggestions = system.GetSuggestions("app")
	fmt.Println("Suggestions for 'app' after deletion:")
	for i, s := range suggestions {
		fmt.Printf("  %d. %s\n", i+1, s)
	}

	// Statistics
	fmt.Println("\n=== System Statistics ===")
	fmt.Printf("Total words in Trie: %d\n", system.trie.GetWordCount())
	fmt.Println("Search history:")
	for prefix, count := range system.GetSearchHistory() {
		fmt.Printf("  '%s': %d searches\n", prefix, count)
	}

	// Top K suggestions
	fmt.Println("\n=== Top 3 Suggestions ===")
	topK := system.trie.AutoCompleteTopK("can", 3)
	fmt.Println("Top 3 suggestions for 'can':")
	for i, s := range topK {
		fmt.Printf("  %d. %s (score: %d)\n", i+1, s.Word, s.Frequency)
	}
}
