# Word Count Engine

## Problem Statement

Implement a word count engine that processes documents and counts the frequency of each word. The engine should:

1. Count words in a document (case-insensitive)
2. Remove punctuation from words
3. Sort words by frequency (descending), then alphabetically
4. Support processing multiple documents
5. Provide top N words by frequency
6. Lookup frequency of specific words

### Key Features

- **Case-Insensitive**: "Hello" and "hello" are treated as the same word
- **Punctuation Removal**: "world!" and "world" are treated as the same word
- **Frequency Sorting**: Words sorted by count (highest first), then alphabetically
- **Multiple Documents**: Can process multiple documents and combine counts

## Requirements

1. Process a document and count word frequencies
2. Get all word counts sorted by frequency
3. Get top N words by frequency
4. Get frequency of a specific word
5. Reset word counts
6. Handle case-insensitive matching
7. Remove punctuation from words

## Test Cases

### Test Case 1: Basic Word Counting
- **Input**: `"Practice makes perfect. you'll only get Perfect by practice. just practice!"`
- **Expected Output**:
  - practice: 3
  - perfect: 2
  - by: 1
  - get: 1
  - just: 1
  - makes: 1
  - only: 1
  - youll: 1

### Test Case 2: Top Words
- **Input**: Same as Test Case 1, get top 3
- **Expected Output**:
  1. practice: 3
  2. perfect: 2
  3. by: 1 (or other word with count 1, alphabetically first)

### Test Case 3: Word Frequency Lookup
- **Input**: Lookup "practice" in the document from Test Case 1
- **Expected Output**: 3

### Test Case 4: Multiple Documents
- **Input**:
  - Document 1: "The quick brown fox jumps over the lazy dog"
  - Document 2: "The lazy dog was jumped over by the quick brown fox"
  - Document 3: "A quick brown fox is always faster than a lazy dog"
- **Expected Output**:
  - the: 3
  - quick: 3
  - brown: 3
  - fox: 3
  - lazy: 3
  - dog: 3
  - (other words with lower counts)

### Test Case 5: Case-Insensitive and Punctuation
- **Input**: `"Hello, World! HELLO world. hello-WORLD!"`
- **Expected Output**:
  - hello: 3
  - world: 3

## Running the Solutions

### Go

To run the Go solution:

```bash
cd go/stripe/word_count_engine
go run main.go
```

### C++

To compile and run the C++ solution:

```bash
cd go/stripe/word_count_engine
g++ -o main main.cpp
./main
```

