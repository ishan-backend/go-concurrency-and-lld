# LFU Cache Implementation

A Go library for handling LFU (Least Frequently Used) cache operations in **O(1)** time complexity:

- **Get**: Retrieve a value by key
- **Set**: Insert or update a key-value pair
- **GetLFU**: Get the least frequently used key

## Paper Reference

This implementation is based on the paper:
**"An O(1) Algorithm for Implementing the LFU Cache Eviction Scheme"** by K. Shah, A. Mitra, and D. Matani

[Paper Link](https://github.com/teivah/go-lfu/blob/main/res/a-constant-algorithm-for-implementing-the-lfu-cache-eviction-scheme.pdf)

## Algorithm Overview

The paper describes an elegant O(1) solution using a **doubly linked list of frequency nodes**, where each node represents a frequency level and contains all keys accessed that many times.

### Key Data Structures

1. **Frequency Node (`freqNode`)**: 
   - Represents a frequency level (e.g., frequency 1, 2, 3...)
   - Contains a map of keys (`items`) that have been accessed exactly `value` times
   - Linked via `prev` and `next` pointers in a doubly linked list
   - Nodes are ordered by frequency (lowest to highest)

2. **Cache Item (`item`)**:
   - Stores the actual value (`data`)
   - Maintains a pointer to its parent frequency node (`parent`)

3. **Cache (`Cache`)**:
   - `byKey`: Hash map for O(1) key lookup → item mapping
   - `freq`: Sentinel node for the doubly linked list of frequency nodes
   - `maxLength`: Maximum cache capacity

### How O(1) Operations Are Achieved

#### **Get Operation** (O(1))
```go
func (c *Cache[K, V]) Get(k K) (V, bool)
```

1. **Lookup**: O(1) hash map lookup via `byKey[k]`
2. **Increment Frequency**: 
   - Move key from current frequency node to next frequency node (freq+1)
   - If next node doesn't exist or has wrong frequency, create it
   - Remove key from old frequency node
   - If old frequency node becomes empty, delete it
3. **Return**: O(1) return of the value

**Time Complexity**: O(1) - All operations are hash map lookups/insertions or linked list pointer updates

#### **Set Operation** (O(1))
```go
func (c *Cache[K, V]) Set(k K, v V)
```

1. **Check Existence**: O(1) hash map lookup
   - If key exists, update value and return (no frequency change)
2. **Eviction** (if needed): O(1) removal from lowest frequency node
3. **Insert**: 
   - Add key to frequency 1 node (create if doesn't exist)
   - Insert into `byKey` hash map

**Time Complexity**: O(1) - Hash map operations and linked list insertions

#### **GetLFU Operation** (O(1))
```go
func (c *Cache[K, _]) GetLFU() (K, bool)
```

1. **Find Minimum Frequency**: The first node after sentinel (`freq.next`) always has the minimum frequency
2. **Return Any Key**: Return any key from that node's items map

**Time Complexity**: O(1) - Direct access to first frequency node

### Code Walkthrough

#### Frequency Node Structure
```7:12:go/lfu-cache/main.go
type freqNode[K comparable, V any] struct {
	value int
	items map[K]struct{}
	prev  *freqNode[K, V]
	next  *freqNode[K, V]
}
```

- `value`: The frequency count (1, 2, 3, ...)
- `items`: Set of keys accessed exactly `value` times
- `prev`/`next`: Doubly linked list pointers

#### Get Operation - Frequency Increment
```63:74:go/lfu-cache/main.go
	freq := tmp.parent
	nextFreq := freq.next
	if nextFreq == c.freq || nextFreq.value != freq.value+1 {
		nextFreq = getNewNode(freq.value+1, freq, nextFreq)
	}
	nextFreq.items[k] = struct{}{}
	tmp.parent = nextFreq
	delete(freq.items, k)
	if len(freq.items) == 0 {
		deleteNode(freq)
	}
	return tmp.data, true
```

When a key is accessed:
1. Check if next frequency node exists and has correct value
2. If not, create a new frequency node
3. Move key to the new frequency node
4. Remove from old node and clean up if empty

#### Set Operation - Insert at Frequency 1
```93:99:go/lfu-cache/main.go
	freq := c.freq.next
	if freq.value != 1 {
		freq = getNewNode(1, c.freq, freq)
	}

	freq.items[k] = struct{}{}
	c.byKey[k] = newItem(v, freq)
```

New items always start at frequency 1:
1. Check if frequency 1 node exists
2. Create it if needed (right after sentinel)
3. Add key to frequency 1 node
4. Store in hash map

#### Eviction - Remove LFU Item
```118:134:go/lfu-cache/main.go
func (c *Cache[K, _]) deleteLFU() {
	var key K
	head := c.freq.next

	for k := range head.items {
		key = k
		break
	}

	if len(head.items) == 1 {
		c.freq.next = c.freq.next.next
	} else {
		delete(head.items, key)
	}

	delete(c.byKey, key)
}
```

Eviction is O(1):
1. Get any key from the first frequency node (lowest frequency)
2. Remove from frequency node (delete node if it becomes empty)
3. Remove from hash map

### Why This Works

The brilliance of this approach is:
- **Frequency nodes are ordered**: The doubly linked list maintains frequency order, so `freq.next` always points to the minimum frequency
- **No traversal needed**: We never need to search for minimum frequency - it's always at the head
- **Lazy node creation**: Frequency nodes are created only when needed (when a key moves to that frequency)
- **Hash map for O(1) access**: Direct key lookup without searching

### Thread Safety

The implementation uses `sync.RWMutex` for thread-safe operations:
- `GetLFU` uses `RLock` (read-only)
- `Get` and `Set` use `Lock` (read-write)

### Example Usage

```go
cache := NewCache[int, string](2) // Capacity of 2

cache.Set(1, "one")
cache.Set(2, "two")

value, found := cache.Get(1) // Returns "one", true
lfuKey, found := cache.GetLFU() // Returns 2 (least frequently used)

cache.Set(3, "three") // Evicts key 2 (LFU)
```

### Testing

Run tests with:
```bash
go test -v
```

The test suite includes:
- Basic cache operations (Get, Set, GetLFU)
- Eviction behavior
- Frequency tracking
- Race condition testing
