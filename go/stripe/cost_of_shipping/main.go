package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Route represents a shipping route
type Route struct {
	Origin      string
	Destination string
	Carrier     string
	Cost        int
}

// ShippingCostCalculator handles shipping cost calculations
type ShippingCostCalculator struct {
	directRoutes map[string]map[string]*Route // origin -> destination -> route
}

// NewShippingCostCalculator creates a new calculator from input string
func NewShippingCostCalculator(inputString string) (*ShippingCostCalculator, error) {
	calc := &ShippingCostCalculator{
		directRoutes: make(map[string]map[string]*Route),
	}

	// Parse input string: "US,UK,UPS,5:US,CA,FedEx,3:CA,UK,DHL,7"
	segments := strings.Split(inputString, ":")

	for _, segment := range segments {
		parts := strings.Split(segment, ",")
		if len(parts) != 4 {
			return nil, fmt.Errorf("invalid segment format: %s", segment)
		}

		origin := strings.TrimSpace(parts[0])
		destination := strings.TrimSpace(parts[1])
		carrier := strings.TrimSpace(parts[2])
		cost, err := strconv.Atoi(strings.TrimSpace(parts[3]))
		if err != nil {
			return nil, fmt.Errorf("invalid cost in segment %s: %v", segment, err)
		}

		route := &Route{
			Origin:      origin,
			Destination: destination,
			Carrier:     carrier,
			Cost:        cost,
		}

		// Store route in the graph
		if calc.directRoutes[origin] == nil {
			calc.directRoutes[origin] = make(map[string]*Route)
		}
		calc.directRoutes[origin][destination] = route
	}

	return calc, nil
}

// FindDirectCost finds the cost to ship directly from origin to destination
// Returns cost and carrier if route exists, otherwise returns error
func (calc *ShippingCostCalculator) FindDirectCost(origin, destination string) (int, string, error) {
	if routes, exists := calc.directRoutes[origin]; exists {
		if route, found := routes[destination]; found {
			return route.Cost, route.Carrier, nil
		}
	}
	return 0, "", fmt.Errorf("no direct route found from %s to %s", origin, destination)
}

// Path represents a shipping path with intermediate steps
type Path struct {
	Route     []*Route
	TotalCost int
}

// FindPathWithIntermediate finds the cheapest path from origin to destination
// allowing intermediate steps. Returns the path and total cost.
func (calc *ShippingCostCalculator) FindPathWithIntermediate(origin, destination string) (*Path, error) {
	// Use BFS to find the shortest path (cheapest cost)
	type queueItem struct {
		location string
		path     []*Route
		cost     int
		visited  map[string]bool
	}

	queue := []queueItem{{
		location: origin,
		path:     []*Route{},
		cost:     0,
		visited:  map[string]bool{origin: true},
	}}

	var bestPath *Path
	minCost := -1

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// If we reached the destination, check if this is the best path
		if current.location == destination {
			if minCost == -1 || current.cost < minCost {
				minCost = current.cost
				bestPath = &Path{
					Route:     current.path,
					TotalCost: current.cost,
				}
			}
			continue
		}

		// Explore all possible routes from current location
		if routes, exists := calc.directRoutes[current.location]; exists {
			for dest, route := range routes {
				// Skip if already visited (avoid cycles)
				if current.visited[dest] {
					continue
				}

				// Create new visited set
				newVisited := make(map[string]bool)
				for k, v := range current.visited {
					newVisited[k] = v
				}
				newVisited[dest] = true

				// Create new path
				newPath := make([]*Route, len(current.path))
				copy(newPath, current.path)
				newPath = append(newPath, route)

				queue = append(queue, queueItem{
					location: dest,
					path:     newPath,
					cost:     current.cost + route.Cost,
					visited:  newVisited,
				})
			}
		}
	}

	if bestPath == nil {
		return nil, fmt.Errorf("no path found from %s to %s", origin, destination)
	}

	return bestPath, nil
}

// FindAllPathsWithIntermediate finds all possible paths from origin to destination
// Returns all paths sorted by cost (cheapest first)
func (calc *ShippingCostCalculator) FindAllPathsWithIntermediate(origin, destination string) ([]*Path, error) {
	var allPaths []*Path

	var dfs func(current string, path []*Route, cost int, visited map[string]bool)
	dfs = func(current string, path []*Route, cost int, visited map[string]bool) {
		if current == destination {
			newPath := make([]*Route, len(path))
			copy(newPath, path)
			allPaths = append(allPaths, &Path{
				Route:     newPath,
				TotalCost: cost,
			})
			return
		}

		if routes, exists := calc.directRoutes[current]; exists {
			for dest, route := range routes {
				if visited[dest] {
					continue
				}

				newVisited := make(map[string]bool)
				for k, v := range visited {
					newVisited[k] = v
				}
				newVisited[dest] = true

				newPath := make([]*Route, len(path))
				copy(newPath, path)
				newPath = append(newPath, route)

				dfs(dest, newPath, cost+route.Cost, newVisited)
			}
		}
	}

	dfs(origin, []*Route{}, 0, map[string]bool{origin: true})

	if len(allPaths) == 0 {
		return nil, fmt.Errorf("no path found from %s to %s", origin, destination)
	}

	// Sort by cost (simple bubble sort for small datasets)
	for i := 0; i < len(allPaths)-1; i++ {
		for j := i + 1; j < len(allPaths); j++ {
			if allPaths[i].TotalCost > allPaths[j].TotalCost {
				allPaths[i], allPaths[j] = allPaths[j], allPaths[i]
			}
		}
	}

	return allPaths, nil
}

func main() {
	inputString := "US,UK,UPS,5:US,CA,FedEx,3:CA,UK,DHL,7"

	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println("Shipping Cost Calculator")
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Printf("Input: %s\n\n", inputString)

	calc, err := NewShippingCostCalculator(inputString)
	if err != nil {
		fmt.Printf("Error creating calculator: %v\n", err)
		return
	}

	// Part 1: Direct route cost
	fmt.Println("Part 1: Direct Route Costs")
	fmt.Println("-" + strings.Repeat("-", 60))

	testCases := []struct {
		origin      string
		destination string
	}{
		{"US", "UK"},
		{"US", "CA"},
		{"CA", "UK"},
		{"UK", "US"}, // Should not exist
	}

	for _, tc := range testCases {
		cost, carrier, err := calc.FindDirectCost(tc.origin, tc.destination)
		if err != nil {
			fmt.Printf("  %s -> %s: %v\n", tc.origin, tc.destination, err)
		} else {
			fmt.Printf("  %s -> %s: Cost = $%d (Carrier: %s)\n", tc.origin, tc.destination, cost, carrier)
		}
	}

	// Part 2: Route with intermediate steps
	fmt.Println("\nPart 2: Routes with Intermediate Steps")
	fmt.Println("-" + strings.Repeat("-", 60))

	paths := []struct {
		origin      string
		destination string
	}{
		{"US", "UK"},
		{"US", "CA"},
		{"CA", "UK"},
	}

	for _, p := range paths {
		fmt.Printf("\nFinding cheapest path from %s to %s:\n", p.origin, p.destination)

		// Find cheapest path
		bestPath, err := calc.FindPathWithIntermediate(p.origin, p.destination)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
			continue
		}

		fmt.Printf("  Cheapest Path (Cost: $%d):\n", bestPath.TotalCost)
		for i, route := range bestPath.Route {
			fmt.Printf("    Step %d: %s -> %s via %s (Cost: $%d)\n",
				i+1, route.Origin, route.Destination, route.Carrier, route.Cost)
		}

		// Find all paths
		allPaths, err := calc.FindAllPathsWithIntermediate(p.origin, p.destination)
		if err == nil && len(allPaths) > 1 {
			fmt.Printf("\n  All possible paths:\n")
			for idx, path := range allPaths {
				fmt.Printf("    Path %d (Cost: $%d): ", idx+1, path.TotalCost)
				steps := []string{}
				for _, route := range path.Route {
					steps = append(steps, fmt.Sprintf("%s->%s", route.Origin, route.Destination))
				}
				fmt.Printf("%s\n", strings.Join(steps, " -> "))
			}
		}
	}

	// Example: US -> CA -> UK
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example: US -> UK (comparing direct vs intermediate)")
	fmt.Println(strings.Repeat("=", 60))

	directCost, carrier, err := calc.FindDirectCost("US", "UK")
	if err == nil {
		fmt.Printf("Direct route: $%d (Carrier: %s)\n", directCost, carrier)
	} else {
		fmt.Printf("Direct route: Not available\n")
	}

	// Find all paths to show intermediate option
	allPaths, err := calc.FindAllPathsWithIntermediate("US", "UK")
	if err != nil {
		fmt.Printf("Error finding paths: %v\n", err)
	} else {
		fmt.Printf("\nAll available paths:\n")
		for idx, path := range allPaths {
			if len(path.Route) == 1 {
				fmt.Printf("  Path %d (Direct): $%d\n", idx+1, path.TotalCost)
				fmt.Printf("    %s -> %s via %s\n",
					path.Route[0].Origin, path.Route[0].Destination, path.Route[0].Carrier)
			} else {
				fmt.Printf("  Path %d (Via intermediate): $%d\n", idx+1, path.TotalCost)
				for i, route := range path.Route {
					fmt.Printf("    Step %d: %s -> %s via %s (Cost: $%d)\n",
						i+1, route.Origin, route.Destination, route.Carrier, route.Cost)
				}
			}
		}

		// Show the intermediate path specifically
		for _, path := range allPaths {
			if len(path.Route) > 1 {
				fmt.Printf("\nIntermediate route (US -> CA -> UK): $%d\n", path.TotalCost)
				fmt.Printf("  US -> CA: $%d\n", path.Route[0].Cost)
				fmt.Printf("  CA -> UK: $%d\n", path.Route[1].Cost)
				fmt.Printf("  Total: $%d\n", path.TotalCost)
				break
			}
		}
	}
}
