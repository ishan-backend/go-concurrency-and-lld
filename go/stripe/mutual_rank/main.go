package main

import (
	"fmt"
	"sort"
	"strings"
)

// MutualRankCalculator calculates mutual rankings between users
type MutualRankCalculator struct {
	rankings map[string][]string // user -> list of ranked users (in order of preference)
}

// NewMutualRankCalculator creates a new calculator
func NewMutualRankCalculator() *MutualRankCalculator {
	return &MutualRankCalculator{
		rankings: make(map[string][]string),
	}
}

// AddRanking adds a user's ranking of other users
// rankings is a comma-separated string of user IDs in order of preference
func (m *MutualRankCalculator) AddRanking(userID, rankings string) {
	users := strings.Split(rankings, ",")
	// Trim whitespace
	for i := range users {
		users[i] = strings.TrimSpace(users[i])
	}
	m.rankings[userID] = users
}

// GetRank gets the rank of targetUser in userID's ranking (1-indexed)
// Returns 0 if not found
func (m *MutualRankCalculator) GetRank(userID, targetUser string) int {
	rankedUsers, exists := m.rankings[userID]
	if !exists {
		return 0
	}
	
	for i, user := range rankedUsers {
		if user == targetUser {
			return i + 1 // 1-indexed
		}
	}
	return 0
}

// GetMutualRank gets the mutual rank between two users
// Mutual rank is the sum of ranks: rank of userB in userA's list + rank of userA in userB's list
// Lower mutual rank means better mutual preference
func (m *MutualRankCalculator) GetMutualRank(userA, userB string) (int, bool) {
	rankA := m.GetRank(userA, userB)
	rankB := m.GetRank(userB, userA)
	
	if rankA == 0 || rankB == 0 {
		return 0, false // One or both users don't have the other in their ranking
	}
	
	return rankA + rankB, true
}

// FindBestMatches finds the best mutual matches for a given user
// Returns a list of users sorted by mutual rank (best first)
func (m *MutualRankCalculator) FindBestMatches(userID string) []string {
	type match struct {
		userID     string
		mutualRank int
	}
	
	var matches []match
	
	// Get all users that this user has ranked
	rankedUsers, exists := m.rankings[userID]
	if !exists {
		return []string{}
	}
	
	// Check mutual rank with each ranked user
	for _, otherUser := range rankedUsers {
		if mutualRank, ok := m.GetMutualRank(userID, otherUser); ok {
			matches = append(matches, match{
				userID:     otherUser,
				mutualRank: mutualRank,
			})
		}
	}
	
	// Sort by mutual rank (lower is better), then by userID for stability
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].mutualRank != matches[j].mutualRank {
			return matches[i].mutualRank < matches[j].mutualRank
		}
		return matches[i].userID < matches[j].userID
	})
	
	result := make([]string, len(matches))
	for i, m := range matches {
		result[i] = m.userID
	}
	
	return result
}

// GetAllMutualRanks returns all mutual ranks for a given user
func (m *MutualRankCalculator) GetAllMutualRanks(userID string) map[string]int {
	result := make(map[string]int)
	
	rankedUsers, exists := m.rankings[userID]
	if !exists {
		return result
	}
	
	for _, otherUser := range rankedUsers {
		if mutualRank, ok := m.GetMutualRank(userID, otherUser); ok {
			result[otherUser] = mutualRank
		}
	}
	
	return result
}

func main() {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Mutual Rank Calculator")
	fmt.Println(strings.Repeat("=", 60))
	
	calc := NewMutualRankCalculator()
	
	// Example 1: Basic mutual rankings
	fmt.Println("\nExample 1: Basic Mutual Rankings")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	calc.AddRanking("alice", "bob,charlie,dave")
	calc.AddRanking("bob", "alice,charlie,dave")
	calc.AddRanking("charlie", "alice,bob,dave")
	calc.AddRanking("dave", "alice,bob,charlie")
	
	fmt.Println("Rankings:")
	fmt.Println("  alice: bob, charlie, dave")
	fmt.Println("  bob: alice, charlie, dave")
	fmt.Println("  charlie: alice, bob, dave")
	fmt.Println("  dave: alice, bob, charlie")
	fmt.Println()
	
	// Check mutual ranks
	testPairs := [][]string{
		{"alice", "bob"},
		{"alice", "charlie"},
		{"bob", "charlie"},
	}
	
	for _, pair := range testPairs {
		userA, userB := pair[0], pair[1]
		rankA := calc.GetRank(userA, userB)
		rankB := calc.GetRank(userB, userA)
		mutualRank, ok := calc.GetMutualRank(userA, userB)
		
		if ok {
			fmt.Printf("  %s <-> %s: rank(%s->%s)=%d, rank(%s->%s)=%d, mutual rank=%d\n",
				userA, userB, userA, userB, rankA, userB, userA, rankB, mutualRank)
		} else {
			fmt.Printf("  %s <-> %s: No mutual ranking\n", userA, userB)
		}
	}
	
	// Find best matches for alice
	fmt.Println("\nBest matches for alice:")
	bestMatches := calc.FindBestMatches("alice")
	for i, user := range bestMatches {
		mutualRank, _ := calc.GetMutualRank("alice", user)
		fmt.Printf("  %d. %s (mutual rank: %d)\n", i+1, user, mutualRank)
	}
	
	// Example 2: Asymmetric rankings
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 2: Asymmetric Rankings")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	calc2 := NewMutualRankCalculator()
	calc2.AddRanking("user1", "user2,user3,user4")
	calc2.AddRanking("user2", "user4,user1,user3")
	calc2.AddRanking("user3", "user1,user2,user4")
	calc2.AddRanking("user4", "user2,user1,user3")
	
	fmt.Println("Rankings:")
	fmt.Println("  user1: user2, user3, user4")
	fmt.Println("  user2: user4, user1, user3")
	fmt.Println("  user3: user1, user2, user4")
	fmt.Println("  user4: user2, user1, user3")
	fmt.Println()
	
	// Check mutual ranks
	testPairs2 := [][]string{
		{"user1", "user2"},
		{"user1", "user3"},
		{"user2", "user4"},
	}
	
	for _, pair := range testPairs2 {
		userA, userB := pair[0], pair[1]
		rankA := calc2.GetRank(userA, userB)
		rankB := calc2.GetRank(userB, userA)
		mutualRank, ok := calc2.GetMutualRank(userA, userB)
		
		if ok {
			fmt.Printf("  %s <-> %s: rank(%s->%s)=%d, rank(%s->%s)=%d, mutual rank=%d\n",
				userA, userB, userA, userB, rankA, userB, userA, rankB, mutualRank)
		}
	}
	
	// Find best matches for user1
	fmt.Println("\nBest matches for user1:")
	bestMatches2 := calc2.FindBestMatches("user1")
	for i, user := range bestMatches2 {
		mutualRank, _ := calc2.GetMutualRank("user1", user)
		fmt.Printf("  %d. %s (mutual rank: %d)\n", i+1, user, mutualRank)
	}
	
	// Example 3: One-way ranking (no mutual rank)
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 3: One-Way Ranking")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	calc3 := NewMutualRankCalculator()
	calc3.AddRanking("A", "B,C")
	calc3.AddRanking("B", "C,D")
	// A ranks B, but B doesn't rank A
	
	mutualRank, ok := calc3.GetMutualRank("A", "B")
	if ok {
		fmt.Printf("A <-> B: mutual rank = %d\n", mutualRank)
	} else {
		fmt.Println("A <-> B: No mutual ranking (B doesn't rank A)")
	}
}

