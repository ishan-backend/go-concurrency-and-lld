# Apache Lucene - Low Level Design Document

## Table of Contents
1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Core Data Structures](#core-data-structures)
4. [Algorithms](#algorithms)
5. [Components](#components)
6. [Design Patterns](#design-patterns)
7. [Performance Optimizations](#performance-optimizations)

---

## Overview

Apache Lucene is a high-performance, full-featured text search engine library written in Java. It provides powerful indexing and search capabilities through an inverted index data structure. Lucene is designed to handle large-scale text search efficiently.

### Key Features
- **Full-text search**: Fast, scalable text search across documents
- **Field-based indexing**: Different fields can have different indexing strategies
- **Ranking**: Relevance scoring using TF-IDF, BM25, and custom scoring
- **Faceted search**: Support for filtering and faceting
- **Distributed search**: Can be distributed across multiple machines
- **Pluggable analyzers**: Customizable text analysis

---

## Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        Application Layer                      │
└───────────────────────┬─────────────────────────────────────┘
                        │
        ┌───────────────┴───────────────┐
        │                               │
┌───────▼────────┐            ┌────────▼────────┐
│  IndexWriter   │            │   IndexReader   │
│  (Indexing)    │            │   (Searching)   │
└───────┬────────┘            └────────┬────────┘
        │                               │
        │                               │
┌───────▼──────────────────────────────▼────────┐
│            Index Files (Segments)              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐     │
│  │ Segment1 │  │ Segment2 │  │ SegmentN │     │
│  └──────────┘  └──────────┘  └──────────┘     │
└─────────────────────────────────────────────────┘
```

### Core Components Flow

1. **Indexing Flow**: Document → Analyzer → Tokenizer → IndexWriter → Segments
2. **Search Flow**: Query → QueryParser → IndexReader → Searcher → Results

---

## Core Data Structures

### 1. Inverted Index

The fundamental data structure in Lucene is the **inverted index**. Instead of mapping documents to terms, it maps terms to documents.

#### Structure:
```
Term → Posting List
"lucene" → [doc1, doc3, doc7]
"search" → [doc1, doc2, doc5]
"index"  → [doc2, doc3, doc7]
```

#### Implementation Details:
- **Term Dictionary**: Maps terms to their posting lists
- **Posting Lists**: Lists of documents containing the term, with positions and frequencies
- **Term Frequency (TF)**: Number of times a term appears in a document
- **Document Frequency (DF)**: Number of documents containing the term

### 2. Document Store (.fdt, .fdx files)

Stores the original document fields for retrieval.

#### Structure:
```
Document Store:
├── Field 1: "title" → "Apache Lucene"
├── Field 2: "content" → "Lucene is a search library..."
└── Field 3: "author" → "John Doe"
```

#### File Format:
- **.fdt**: Field data file (stores actual field values)
- **.fdx**: Field index file (points to field data locations)

### 3. Term Dictionary (.tim, .tip files)

Maps terms to their posting lists efficiently.

#### Structure:
```
Term Dictionary (FST - Finite State Transducer):
┌─────────┐
│ "apple" │ → Posting List Pointer
│ "apply" │ → Posting List Pointer
│ "april" │ → Posting List Pointer
└─────────┘
```

#### Implementation:
- Uses **FST (Finite State Transducer)** for compact representation
- Prefix compression for similar terms
- **.tim**: Term info file (stores term metadata)
- **.tip**: Term index file (index into .tim file)

### 4. Posting Lists (.doc, .pos, .pay files)

Stores document IDs, positions, and payloads for each term.

#### Structure:
```
Posting List for term "lucene":
┌─────────────────────────────────────────┐
│ DocID: 1, TF: 3, Positions: [5, 12, 45] │
│ DocID: 3, TF: 2, Positions: [8, 20]    │
│ DocID: 7, TF: 1, Positions: [15]       │
└─────────────────────────────────────────┘
```

#### File Format:
- **.doc**: Document IDs and term frequencies (delta-encoded)
- **.pos**: Term positions within documents (delta-encoded)
- **.pay**: Payloads (optional metadata per term occurrence)

#### Encoding:
- **Delta Encoding**: Store differences between consecutive values
- **Variable-byte encoding**: Compress integers efficiently
- **Skip Lists**: For fast intersection of posting lists

### 5. Norms (.nvd, .nvm files)

Stores normalization factors for scoring.

#### Purpose:
- Field length normalization
- Boost factors
- Used in TF-IDF and BM25 scoring

#### Structure:
```
Norms:
Doc1 → 1.0 (normalized length)
Doc2 → 0.8
Doc3 → 1.2
```

### 6. Term Vectors (.tvx, .tvd, .tvf files)

Stores per-document term frequency information.

#### Use Cases:
- Highlighting
- More Like This queries
- Term frequency analysis

### 7. Segment Files

Lucene divides the index into **segments** (immutable, append-only).

#### Segment Structure:
```
Segment:
├── .si (Segment Info)
├── .cfs (Compound File - optional)
├── .fdt, .fdx (Field Store)
├── .tim, .tip (Term Dictionary)
├── .doc, .pos, .pay (Posting Lists)
├── .nvd, .nvm (Norms)
└── .tvx, .tvd, .tvf (Term Vectors)
```

#### Segment Properties:
- **Immutable**: Once written, segments are never modified
- **Merge**: Old segments are merged into new ones
- **Deletion**: Uses .del file to mark deleted documents

### 8. Skip Lists

Optimization for fast posting list intersection.

#### Structure:
```
Posting List with Skip Pointers:
[1, 3, 7, 15, 20, 25, 30, ...]
 ↑     ↑     ↑
Skip  Skip  Skip
```

#### Algorithm:
- Skip every K documents
- Allows binary search-like skipping
- Reduces comparisons during intersection

---

## Algorithms

### 1. Indexing Algorithm

#### Process Flow:

```
1. Document Input
   ↓
2. Analyzer (Tokenization, Lowercasing, Stemming, Stop words)
   ↓
3. Build In-Memory Index (RAM Buffer)
   ↓
4. When buffer full → Flush to Disk Segment
   ↓
5. Merge Segments (Periodically)
```

#### Detailed Steps:

**Step 1: Document Analysis**
```java
Document doc = new Document();
doc.add(new TextField("content", "Apache Lucene is great"));
Analyzer analyzer = new StandardAnalyzer();
TokenStream stream = analyzer.tokenStream("content", doc);
```

**Step 2: Token Processing**
- Tokenize: Split into tokens
- Normalize: Lowercase, remove accents
- Filter: Remove stop words, apply stemming

**Step 3: In-Memory Indexing**
- Build inverted index in RAM
- Use hash maps for term → posting list mapping
- Accumulate term frequencies and positions

**Step 4: Flush to Segment**
- When RAM buffer exceeds threshold
- Write segment files (.tim, .doc, .fdt, etc.)
- Use compression (LZ4, Deflate)

**Step 5: Segment Merging**
- Merge multiple small segments into larger ones
- Reduce number of segments
- Optimize for query performance

### 2. Query Processing Algorithm

#### Types of Queries:

**Term Query:**
```
Query: "lucene"
→ Lookup term in dictionary
→ Retrieve posting list
→ Return matching documents
```

**Boolean Query:**
```
Query: "lucene AND search"
→ Get posting lists for "lucene" and "search"
→ Intersect posting lists (AND)
→ Union posting lists (OR)
→ Return results
```

**Phrase Query:**
```
Query: "apache lucene"
→ Get posting lists for both terms
→ Check if terms appear in sequence
→ Use position information
```

#### Query Processing Steps:

1. **Query Parsing**: Parse query string into query tree
2. **Term Extraction**: Extract terms from query
3. **Dictionary Lookup**: Find terms in term dictionary
4. **Posting List Retrieval**: Load posting lists
5. **Query Execution**: Execute query operations (AND, OR, NOT)
6. **Scoring**: Calculate relevance scores
7. **Ranking**: Sort by score
8. **Result Collection**: Collect top-K results

### 3. Posting List Intersection Algorithm

#### Algorithm: Two-Pointer Technique

```python
def intersect(posting_list1, posting_list2):
    result = []
    i, j = 0, 0
    
    while i < len(posting_list1) and j < len(posting_list2):
        if posting_list1[i] == posting_list2[j]:
            result.append(posting_list1[i])
            i += 1
            j += 1
        elif posting_list1[i] < posting_list2[j]:
            i += 1
        else:
            j += 1
    
    return result
```

#### Optimization with Skip Lists:

```python
def intersect_with_skips(posting_list1, posting_list2):
    # Use skip pointers to jump ahead
    # Reduces comparisons for sparse intersections
    result = []
    i, j = 0, 0
    
    while i < len(posting_list1) and j < len(posting_list2):
        if posting_list1[i] == posting_list2[j]:
            result.append(posting_list1[i])
            i += 1
            j += 1
        elif posting_list1[i] < posting_list2[j]:
            i = skip_to(posting_list1, i, posting_list2[j])
        else:
            j = skip_to(posting_list2, j, posting_list1[i])
    
    return result
```

### 4. Ranking Algorithms

#### TF-IDF (Term Frequency - Inverse Document Frequency)

**Formula:**
```
score(q, d) = Σ tf(t, d) × idf(t) × boost(t)

where:
- tf(t, d) = term frequency of t in document d
- idf(t) = log(N / df(t))
  - N = total number of documents
  - df(t) = document frequency of term t
- boost(t) = boost factor for term t
```

**Implementation:**
```java
float score = 0.0f;
for (Term term : queryTerms) {
    float tf = getTermFrequency(term, doc);
    float idf = Math.log(numDocs / getDocumentFrequency(term));
    float norm = getNorm(doc);
    score += tf * idf * boost * norm;
}
```

#### BM25 (Best Matching 25)

**Formula:**
```
score(q, d) = Σ idf(t) × (tf(t, d) × (k1 + 1)) / (tf(t, d) + k1 × (1 - b + b × |d|/avgdl))

where:
- k1 = 1.2 (tuning parameter)
- b = 0.75 (tuning parameter)
- |d| = document length
- avgdl = average document length
```

**Advantages over TF-IDF:**
- Better handling of document length normalization
- Saturation function prevents over-weighting high frequencies
- More robust for modern search applications

### 5. Segment Merging Algorithm

#### Merge Policy: Tiered Merge Policy

**Strategy:**
1. Segments are organized into tiers by size
2. Merge segments of similar size
3. Prefer merging smaller segments first

**Algorithm:**
```
1. Identify merge candidates (segments to merge)
2. Calculate merge cost (I/O operations)
3. Select segments with lowest merge cost
4. Create new merged segment
5. Update segment list
6. Delete old segments
```

**Example:**
```
Before Merge:
Segment1: 10MB
Segment2: 10MB
Segment3: 10MB
Segment4: 50MB

After Merge:
Segment1+2+3: 30MB (merged)
Segment4: 50MB
```

### 6. Compression Algorithms

#### Variable-Byte Encoding

**Purpose**: Compress integers efficiently

**Algorithm:**
- Encode integers using variable number of bytes
- Use continuation bit to indicate more bytes
- Smaller numbers use fewer bytes

**Example:**
```
Number: 5 → [00000101] (1 byte)
Number: 300 → [10000010] [01001100] (2 bytes)
```

#### Delta Encoding

**Purpose**: Compress sequences of sorted integers

**Algorithm:**
- Store differences between consecutive values
- First value stored as-is
- Subsequent values stored as deltas

**Example:**
```
Original: [1, 3, 7, 15, 20]
Encoded:  [1, 2, 4, 8, 5]
          (1, 3-1, 7-3, 15-7, 20-15)
```

#### LZ4 Compression

**Purpose**: Fast compression for stored fields

**Characteristics:**
- Fast compression/decompression
- Good compression ratio
- Used for .fdt files

### 7. FST (Finite State Transducer) Construction

**Purpose**: Compact representation of term dictionary

**Properties:**
- Prefix compression
- Suffix compression
- Shared paths for common prefixes/suffixes

**Construction Algorithm:**
```
1. Sort all terms lexicographically
2. Build FST bottom-up
3. Share common prefixes
4. Share common suffixes
5. Store output values (posting list pointers)
```

**Example:**
```
Terms: ["apple", "apply", "april"]
FST Structure:
        a
        |
        p
       / \
      p   r
     / \   \
    l   l   i
   / \   \   \
  e   y   y   l
```

---

## Components

### 1. Analyzer

**Purpose**: Convert text into tokens

**Components:**
- **Tokenizer**: Splits text into tokens
- **TokenFilter**: Modifies tokens (lowercase, stemming, etc.)

**Example:**
```java
Analyzer analyzer = new StandardAnalyzer();
TokenStream stream = analyzer.tokenStream("field", "Apache Lucene");
// Output: ["apache", "lucene"]
```

**Common Analyzers:**
- **StandardAnalyzer**: Standard tokenization, lowercase, stop words
- **WhitespaceAnalyzer**: Split on whitespace only
- **KeywordAnalyzer**: No tokenization
- **Custom Analyzer**: User-defined tokenization pipeline

### 2. IndexWriter

**Purpose**: Write documents to index

**Key Methods:**
```java
IndexWriter writer = new IndexWriter(directory, config);
writer.addDocument(doc);      // Add document
writer.updateDocument(term, doc); // Update document
writer.deleteDocuments(term); // Delete documents
writer.commit();              // Commit changes
writer.close();               // Close writer
```

**Configuration:**
- **RAM Buffer Size**: Controls when to flush
- **Merge Policy**: How to merge segments
- **Similarity**: Scoring algorithm (TF-IDF, BM25)

### 3. IndexReader

**Purpose**: Read from index

**Types:**
- **DirectoryReader**: Reads from directory
- **MultiReader**: Reads from multiple indexes

**Key Methods:**
```java
IndexReader reader = DirectoryReader.open(directory);
int numDocs = reader.numDocs();
TermFreqVector vector = reader.getTermFreqVector(docId, field);
reader.close();
```

### 4. IndexSearcher

**Purpose**: Execute queries and return results

**Key Methods:**
```java
IndexSearcher searcher = new IndexSearcher(reader);
Query query = new TermQuery(new Term("field", "value"));
TopDocs results = searcher.search(query, 10);
for (ScoreDoc scoreDoc : results.scoreDocs) {
    Document doc = searcher.doc(scoreDoc.doc);
}
```

### 5. Query Classes

**TermQuery**: Single term search
```java
Query q = new TermQuery(new Term("content", "lucene"));
```

**BooleanQuery**: Combine multiple queries
```java
BooleanQuery.Builder builder = new BooleanQuery.Builder();
builder.add(new TermQuery(new Term("content", "lucene")), Occur.MUST);
builder.add(new TermQuery(new Term("content", "search")), Occur.MUST);
Query q = builder.build();
```

**PhraseQuery**: Search for phrase
```java
PhraseQuery q = new PhraseQuery.Builder()
    .add(new Term("content", "apache"))
    .add(new Term("content", "lucene"))
    .build();
```

**WildcardQuery**: Pattern matching
```java
Query q = new WildcardQuery(new Term("content", "luc*"));
```

### 6. Directory

**Purpose**: Abstract storage layer

**Implementations:**
- **FSDirectory**: File system storage
- **RAMDirectory**: In-memory storage
- **MMapDirectory**: Memory-mapped files

---

## Design Patterns

### 1. Strategy Pattern

**Usage**: Different scoring algorithms, analyzers, merge policies

```java
// Scoring strategy
Similarity similarity = new BM25Similarity();
// or
Similarity similarity = new ClassicSimilarity(); // TF-IDF
```

### 2. Factory Pattern

**Usage**: Create analyzers, queries, readers

```java
Analyzer analyzer = AnalyzerFactory.createStandardAnalyzer();
```

### 3. Builder Pattern

**Usage**: Construct complex queries, index configurations

```java
BooleanQuery query = new BooleanQuery.Builder()
    .add(mustClause, Occur.MUST)
    .add(shouldClause, Occur.SHOULD)
    .build();
```

### 4. Iterator Pattern

**Usage**: Iterate through posting lists, documents, terms

```java
TermsEnum termsEnum = terms.iterator();
while (termsEnum.next() != null) {
    Term term = new Term(field, termsEnum.term());
}
```

### 5. Template Method Pattern

**Usage**: Query execution, document processing pipeline

### 6. Observer Pattern

**Usage**: IndexWriter events, merge notifications

---

## Performance Optimizations

### 1. Caching Strategies

**Term Dictionary Cache:**
- Cache frequently accessed terms
- Reduce disk I/O

**Field Cache:**
- Cache field values in memory
- Fast access for sorting/faceting

**Filter Cache:**
- Cache filter results
- Reuse for multiple queries

### 2. Memory Management

**RAM Buffer:**
- Accumulate documents in memory before flushing
- Reduces disk I/O operations
- Configurable size

**Memory-Mapped Files:**
- Use OS virtual memory
- Efficient for large indexes
- Automatic paging

### 3. I/O Optimizations

**Sequential Access:**
- Posting lists stored sequentially
- Optimize for sequential reads

**Compression:**
- Compress stored fields
- Reduce disk space
- Trade-off: CPU vs. I/O

**Batching:**
- Batch multiple operations
- Reduce system calls

### 4. Query Optimizations

**Query Rewriting:**
- Simplify queries
- Remove redundant clauses

**Early Termination:**
- Stop after finding top-K results
- Don't process entire posting lists

**Boolean Query Optimization:**
- Process smallest posting lists first
- Use skip lists for intersection

### 5. Index Optimizations

**Segment Merging:**
- Merge small segments
- Reduce number of segments
- Improve query performance

**Deletion Handling:**
- Mark deletions, don't rewrite segments
- Apply deletions during merge

**Field Selection:**
- Store only necessary fields
- Use field selectors

### 6. Concurrency

**Reader-Writer Lock:**
- Multiple readers, single writer
- IndexReader is thread-safe
- IndexWriter requires synchronization

**Near Real-Time Search:**
- Use IndexReader.reopen() for updates
- Don't block searches during indexing

---

## Data Flow Examples

### Indexing Flow

```
1. Document Input
   Document: {id: 1, title: "Lucene Guide", content: "Apache Lucene..."}
   
2. Analysis
   Analyzer → Tokens: ["apache", "lucene", "guide"]
   
3. Index Building (In-Memory)
   Term Dictionary:
   - "apache" → [doc1]
   - "lucene" → [doc1]
   - "guide" → [doc1]
   
4. Flush to Segment
   Write .tim, .doc, .fdt files
   
5. Merge (if needed)
   Merge segments into larger segment
```

### Query Flow

```
1. Query Input
   Query: "lucene AND search"
   
2. Query Parsing
   BooleanQuery:
   - MUST: TermQuery("lucene")
   - MUST: TermQuery("search")
   
3. Dictionary Lookup
   - "lucene" → Posting List: [1, 3, 7]
   - "search" → Posting List: [1, 2, 5]
   
4. Posting List Intersection
   Intersect([1,3,7], [1,2,5]) → [1]
   
5. Scoring
   Calculate BM25 score for doc1
   
6. Ranking
   Sort by score
   
7. Result
   Return top-K documents
```

---

## Scalability Considerations

### 1. Horizontal Scaling

**Sharding:**
- Split index across multiple machines
- Route queries to appropriate shards
- Merge results from all shards

**Replication:**
- Replicate index for high availability
- Load balance queries across replicas

### 2. Vertical Scaling

**Memory:**
- Increase RAM buffer size
- Cache more data in memory

**Disk:**
- Use SSDs for better I/O performance
- RAID for redundancy

### 3. Distributed Search

**Elasticsearch Architecture:**
- Built on Lucene
- Adds distributed capabilities
- Sharding, replication, cluster management

---

## Summary

Apache Lucene is a sophisticated search engine library built on:

1. **Inverted Index**: Core data structure mapping terms to documents
2. **Segments**: Immutable, append-only index files
3. **Compression**: Variable-byte encoding, delta encoding, LZ4
4. **Algorithms**: TF-IDF, BM25, posting list intersection, segment merging
5. **Components**: Analyzer, IndexWriter, IndexReader, IndexSearcher
6. **Optimizations**: Caching, memory management, query optimization

The design prioritizes:
- **Performance**: Fast indexing and searching
- **Scalability**: Handle large document collections
- **Flexibility**: Pluggable components and algorithms
- **Efficiency**: Compact storage, minimal I/O

This architecture enables Lucene to power search in applications ranging from small-scale to enterprise-level systems like Elasticsearch and Solr.

