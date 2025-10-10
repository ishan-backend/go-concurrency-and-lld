package multithreaded_paranthesis_parser

import (
	"runtime"
	"sync"
)

/*
https://www.lldcoding.com/multi-threaded-parentheses-validator
Problem Statement:
You are given a string s of length N consisting only of the characters '(' and ')'. A parentheses string is considered valid if:

Every '(' has a matching ')' to its right.

No prefix of the string contains more ')' than '('.

Write a function:

boolean isValidParallel(String s, int P)
that determines whether s is valid, using P worker threads. Each thread can only read a disjoint substring of s.
*/

/*
Solution:

In a multi-threaded version, the twist is: how to partition the work (divide the string among threads), and then combine the results so that across thread boundaries you still correctly detect mismatches (for example, an opening in thread A may be closed in thread B, or vice versa).

So the design needs:
	•	Per-chunk / per-thread processing (local validation).
	•	A way to combine / merge partial results so that the global string is validated.
	•	Handling boundary conditions (leftover opens / unmatched closes crossing thread boundaries).

This is similar to parallel or distributed bracket-checking algorithms. (One known approach: each chunk returns two values: minimum prefix nesting (i.e., how deep it dips below zero) and net change in nesting; then you merge results)

Coordinator
 ├─ splits into chunks
 ├─ launches goroutines running ChunkProcessor
 ├─ collects PartialResults
 ├─ merges them via Merger into final result
 └─ checks validity

ChunkProcessor (per-chunk)
BracketMatcher (utility)
Merger (merge logic)

*/

func ValidateMultithreaded(s string) bool {
	runes := []rune(s) // rune = int32 used to represent a single character
	n := len(runes)

	noOfThreads := runtime.NumCPU()
	chunkSizeForEveryThread := (n + noOfThreads - 1) / noOfThreads // for even distribution of runes across threads
	type resultItem struct {
		idx    int
		result PartialResult
	}
	threadResults := make([]resultItem, noOfThreads)

	// process
	var wg sync.WaitGroup
	for i := 0; i < noOfThreads; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			start := i * chunkSizeForEveryThread
			end := (i + 1) * chunkSizeForEveryThread
			if end > n {
				end = n
			}

			if start > end {
				return
			}

			chunk := runes[start:end]
			processedResult := ProcessChunk(chunk)
			threadResults[i] = resultItem{
				idx:    i,
				result: processedResult,
			}
		}(i)
	}

	wg.Wait()

	// merge
	finalResult := threadResults[0].result
	for i := 1; i < noOfThreads; i++ {
		finalResult = Merge(finalResult, threadResults[i].result)
	}

	return finalResult.net == 0
}

type PartialResult struct {
	net            int // net change in nesting -> can be +ve/-ve
	minPrefToMatch int // for left result to match right result
	opens          []rune
	closes         []rune
}

func ProcessChunk(s []rune) PartialResult {
	net := 0 // open-close
	var opens []rune
	var closes []rune
	minPrefToMatch := 0

	for _, ch := range s {
		if isOpen(ch) {
			net++
			opens = append(opens, ch)
		} else {
			net--
			if net < 0 {
				minPrefToMatch = net
			}
			closes = append(closes, ch)
		}
	}

	return PartialResult{
		net:            net,
		minPrefToMatch: minPrefToMatch,
		opens:          opens,
		closes:         closes,
	}
}

func Merge(left, right PartialResult) PartialResult {
	finalNet := left.net + right.net // net is additive
	minPrefToMatch := min(left.minPrefToMatch, left.net+right.minPrefToMatch)

	return PartialResult{
		net:            finalNet,
		minPrefToMatch: minPrefToMatch,
	}
}

func isOpen(ch rune) bool {
	return ch == '(' || ch == '{' || ch == '['
}
