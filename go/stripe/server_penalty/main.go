package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// ServerEvent represents an event for a server
type ServerEvent struct {
	ServerID    string
	EventType   string // "response_time", "error", etc.
	Value       float64
	Timestamp   time.Time
}

// ServerPenaltyCalculator calculates penalties for servers
type ServerPenaltyCalculator struct {
	events []ServerEvent
}

// NewServerPenaltyCalculator creates a new calculator
func NewServerPenaltyCalculator() *ServerPenaltyCalculator {
	return &ServerPenaltyCalculator{
		events: make([]ServerEvent, 0),
	}
}

// AddEvent adds an event for a server
func (s *ServerPenaltyCalculator) AddEvent(serverID, eventType string, value float64, timestamp time.Time) {
	s.events = append(s.events, ServerEvent{
		ServerID:  serverID,
		EventType: eventType,
		Value:     value,
		Timestamp: timestamp,
	})
}

// CalculatePenalty calculates penalty for a server based on events
// Penalty rules:
// - Response time > threshold: penalty = (response_time - threshold) * multiplier
// - Errors: penalty = error_count * error_penalty
// - Total penalty = sum of all penalties
func (s *ServerPenaltyCalculator) CalculatePenalty(serverID string, responseTimeThreshold float64, responseTimeMultiplier float64, errorPenalty float64) float64 {
	totalPenalty := 0.0
	
	for _, event := range s.events {
		if event.ServerID != serverID {
			continue
		}
		
		switch event.EventType {
		case "response_time":
			if event.Value > responseTimeThreshold {
				penalty := (event.Value - responseTimeThreshold) * responseTimeMultiplier
				totalPenalty += penalty
			}
		case "error":
			totalPenalty += errorPenalty
		}
	}
	
	return totalPenalty
}

// GetServerPenalties returns penalties for all servers
func (s *ServerPenaltyCalculator) GetServerPenalties(responseTimeThreshold float64, responseTimeMultiplier float64, errorPenalty float64) map[string]float64 {
	penalties := make(map[string]float64)
	serverSet := make(map[string]bool)
	
	// Collect all server IDs
	for _, event := range s.events {
		serverSet[event.ServerID] = true
	}
	
	// Calculate penalty for each server
	for serverID := range serverSet {
		penalties[serverID] = s.CalculatePenalty(serverID, responseTimeThreshold, responseTimeMultiplier, errorPenalty)
	}
	
	return penalties
}

// GetServersByPenalty returns servers sorted by penalty (highest first)
func (s *ServerPenaltyCalculator) GetServersByPenalty(responseTimeThreshold float64, responseTimeMultiplier float64, errorPenalty float64) []string {
	penalties := s.GetServerPenalties(responseTimeThreshold, responseTimeMultiplier, errorPenalty)
	
	type serverPenalty struct {
		serverID string
		penalty  float64
	}
	
	var serverPenalties []serverPenalty
	for serverID, penalty := range penalties {
		serverPenalties = append(serverPenalties, serverPenalty{
			serverID: serverID,
			penalty:  penalty,
		})
	}
	
	// Sort by penalty (highest first), then by server ID
	sort.Slice(serverPenalties, func(i, j int) bool {
		if serverPenalties[i].penalty != serverPenalties[j].penalty {
			return serverPenalties[i].penalty > serverPenalties[j].penalty
		}
		return serverPenalties[i].serverID < serverPenalties[j].serverID
	})
	
	result := make([]string, len(serverPenalties))
	for i, sp := range serverPenalties {
		result[i] = sp.serverID
	}
	
	return result
}

// GetEventsForServer returns all events for a specific server
func (s *ServerPenaltyCalculator) GetEventsForServer(serverID string) []ServerEvent {
	var result []ServerEvent
	for _, event := range s.events {
		if event.ServerID == serverID {
			result = append(result, event)
		}
	}
	return result
}

func main() {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Server Penalty Calculator")
	fmt.Println(strings.Repeat("=", 60))
	
	calc := NewServerPenaltyCalculator()
	
	// Example 1: Response time penalties
	fmt.Println("\nExample 1: Response Time Penalties")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	baseTime := time.Now()
	calc.AddEvent("server1", "response_time", 150.0, baseTime) // 50ms over threshold
	calc.AddEvent("server1", "response_time", 120.0, baseTime.Add(time.Second))
	calc.AddEvent("server2", "response_time", 200.0, baseTime) // 100ms over threshold
	calc.AddEvent("server2", "response_time", 110.0, baseTime.Add(time.Second))
	calc.AddEvent("server3", "response_time", 90.0, baseTime) // Under threshold
	
	responseTimeThreshold := 100.0  // 100ms threshold
	responseTimeMultiplier := 0.5   // $0.50 per ms over threshold
	errorPenalty := 10.0            // $10 per error
	
	fmt.Printf("Response Time Threshold: %.1f ms\n", responseTimeThreshold)
	fmt.Printf("Response Time Multiplier: $%.2f per ms over threshold\n", responseTimeMultiplier)
	fmt.Printf("Error Penalty: $%.2f per error\n\n", errorPenalty)
	
	penalties := calc.GetServerPenalties(responseTimeThreshold, responseTimeMultiplier, errorPenalty)
	
	fmt.Println("Server Penalties:")
	for serverID, penalty := range penalties {
		fmt.Printf("  %s: $%.2f\n", serverID, penalty)
	}
	
	// Server1: (150-100)*0.5 + (120-100)*0.5 = 25 + 10 = 35
	// Server2: (200-100)*0.5 + (110-100)*0.5 = 50 + 5 = 55
	// Server3: 0 (all under threshold)
	
	// Example 2: Combined response time and errors
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 2: Combined Response Time and Errors")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	calc2 := NewServerPenaltyCalculator()
	calc2.AddEvent("serverA", "response_time", 150.0, baseTime)
	calc2.AddEvent("serverA", "error", 1.0, baseTime.Add(time.Second))
	calc2.AddEvent("serverA", "error", 1.0, baseTime.Add(2*time.Second))
	calc2.AddEvent("serverB", "response_time", 120.0, baseTime)
	calc2.AddEvent("serverB", "response_time", 130.0, baseTime.Add(time.Second))
	calc2.AddEvent("serverC", "error", 1.0, baseTime)
	calc2.AddEvent("serverC", "error", 1.0, baseTime.Add(time.Second))
	calc2.AddEvent("serverC", "error", 1.0, baseTime.Add(2*time.Second))
	
	penalties2 := calc2.GetServerPenalties(responseTimeThreshold, responseTimeMultiplier, errorPenalty)
	
	fmt.Println("Server Penalties:")
	for serverID, penalty := range penalties2 {
		events := calc2.GetEventsForServer(serverID)
		fmt.Printf("  %s: $%.2f", serverID, penalty)
		fmt.Printf(" (events: %d)", len(events))
		fmt.Println()
	}
	
	// ServerA: (150-100)*0.5 + 2*10 = 25 + 20 = 45
	// ServerB: (120-100)*0.5 + (130-100)*0.5 = 10 + 15 = 25
	// ServerC: 3*10 = 30
	
	// Example 3: Servers sorted by penalty
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 3: Servers Sorted by Penalty (Highest First)")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	sortedServers := calc2.GetServersByPenalty(responseTimeThreshold, responseTimeMultiplier, errorPenalty)
	
	fmt.Println("Servers sorted by penalty:")
	for i, serverID := range sortedServers {
		penalty := penalties2[serverID]
		fmt.Printf("  %d. %s: $%.2f\n", i+1, serverID, penalty)
	}
	
	// Example 4: Time-windowed penalties (events within last N minutes)
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 4: Recent Events Only (Last 5 minutes)")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	calc3 := NewServerPenaltyCalculator()
	now := time.Now()
	calc3.AddEvent("serverX", "response_time", 150.0, now.Add(-2*time.Minute))      // Recent
	calc3.AddEvent("serverX", "response_time", 200.0, now.Add(-10*time.Minute))   // Old
	calc3.AddEvent("serverY", "error", 1.0, now.Add(-1*time.Minute))                // Recent
	calc3.AddEvent("serverY", "error", 1.0, now.Add(-6*time.Minute))                // Old
	
	// Filter events within last 5 minutes
	timeWindow := 5 * time.Minute
	cutoffTime := now.Add(-timeWindow)
	
	penalties3 := make(map[string]float64)
	serverSet := make(map[string]bool)
	
	for _, event := range calc3.events {
		if event.Timestamp.After(cutoffTime) {
			serverSet[event.ServerID] = true
		}
	}
	
	for serverID := range serverSet {
		penalty := 0.0
		for _, event := range calc3.events {
			if event.ServerID == serverID && event.Timestamp.After(cutoffTime) {
				switch event.EventType {
				case "response_time":
					if event.Value > responseTimeThreshold {
						penalty += (event.Value - responseTimeThreshold) * responseTimeMultiplier
					}
				case "error":
					penalty += errorPenalty
				}
			}
		}
		penalties3[serverID] = penalty
	}
	
	fmt.Println("Server Penalties (last 5 minutes only):")
	for serverID, penalty := range penalties3 {
		fmt.Printf("  %s: $%.2f\n", serverID, penalty)
	}
}

