#include <iostream>
#include <string>
#include <vector>
#include <sstream>
#include <map>
#include <queue>
#include <climits>
#include <algorithm>

// Route represents a shipping route
struct Route {
    std::string origin;
    std::string destination;
    std::string carrier;
    int cost;
    
    Route(const std::string& o, const std::string& d, const std::string& c, int co)
        : origin(o), destination(d), carrier(c), cost(co) {}
};

// Path represents a shipping path with intermediate steps
struct Path {
    std::vector<Route> routes;
    int totalCost;
    
    Path() : totalCost(0) {}
};

// ShippingCostCalculator handles shipping cost calculations
class ShippingCostCalculator {
private:
    // origin -> destination -> route
    std::map<std::string, std::map<std::string, Route*>> directRoutes;
    
public:
    // Parse input string and build route graph
    bool parseInput(const std::string& inputString) {
        // Parse input string: "US,UK,UPS,5:US,CA,FedEx,3:CA,UK,DHL,7"
        std::vector<std::string> segments;
        std::stringstream ss(inputString);
        std::string segment;
        
        while (std::getline(ss, segment, ':')) {
            segments.push_back(segment);
        }
        
        for (const auto& seg : segments) {
            std::vector<std::string> parts;
            std::stringstream segStream(seg);
            std::string part;
            
            while (std::getline(segStream, part, ',')) {
                // Trim whitespace
                part.erase(0, part.find_first_not_of(" \t"));
                part.erase(part.find_last_not_of(" \t") + 1);
                parts.push_back(part);
            }
            
            if (parts.size() != 4) {
                std::cerr << "Invalid segment format: " << seg << std::endl;
                return false;
            }
            
            std::string origin = parts[0];
            std::string destination = parts[1];
            std::string carrier = parts[2];
            int cost;
            try {
                cost = std::stoi(parts[3]);
            } catch (...) {
                std::cerr << "Invalid cost in segment " << seg << std::endl;
                return false;
            }
            
            Route* route = new Route(origin, destination, carrier, cost);
            
            // Store route in the graph
            directRoutes[origin][destination] = route;
        }
        
        return true;
    }
    
    // Find direct cost from origin to destination
    bool findDirectCost(const std::string& origin, const std::string& destination,
                       int& cost, std::string& carrier) {
        if (directRoutes.find(origin) != directRoutes.end()) {
            auto& routes = directRoutes[origin];
            if (routes.find(destination) != routes.end()) {
                cost = routes[destination]->cost;
                carrier = routes[destination]->carrier;
                return true;
            }
        }
        return false;
    }
    
    // Find cheapest path from origin to destination using BFS
    Path* findPathWithIntermediate(const std::string& origin, const std::string& destination) {
        struct QueueItem {
            std::string location;
            std::vector<Route> path;
            int cost;
            std::map<std::string, bool> visited;
        };
        
        std::queue<QueueItem> q;
        QueueItem start;
        start.location = origin;
        start.cost = 0;
        start.visited[origin] = true;
        q.push(start);
        
        Path* bestPath = nullptr;
        int minCost = INT_MAX;
        
        while (!q.empty()) {
            QueueItem current = q.front();
            q.pop();
            
            // If we reached the destination, check if this is the best path
            if (current.location == destination) {
                if (current.cost < minCost) {
                    minCost = current.cost;
                    if (bestPath) delete bestPath;
                    bestPath = new Path();
                    bestPath->routes = current.path;
                    bestPath->totalCost = current.cost;
                }
                continue;
            }
            
            // Explore all possible routes from current location
            if (directRoutes.find(current.location) != directRoutes.end()) {
                auto& routes = directRoutes[current.location];
                for (auto& pair : routes) {
                    const std::string& dest = pair.first;
                    Route* route = pair.second;
                    
                    // Skip if already visited (avoid cycles)
                    if (current.visited.find(dest) != current.visited.end()) {
                        continue;
                    }
                    
                    // Create new visited set
                    std::map<std::string, bool> newVisited = current.visited;
                    newVisited[dest] = true;
                    
                    // Create new path
                    std::vector<Route> newPath = current.path;
                    newPath.push_back(*route);
                    
                    QueueItem next;
                    next.location = dest;
                    next.path = newPath;
                    next.cost = current.cost + route->cost;
                    next.visited = newVisited;
                    q.push(next);
                }
            }
        }
        
        return bestPath;
    }
    
    // Destructor to clean up allocated routes
    ~ShippingCostCalculator() {
        for (auto& originPair : directRoutes) {
            for (auto& destPair : originPair.second) {
                delete destPair.second;
            }
        }
    }
};

int main() {
    std::string inputString = "US,UK,UPS,5:US,CA,FedEx,3:CA,UK,DHL,7";
    
    std::cout << std::string(61, '=') << std::endl;
    std::cout << "Shipping Cost Calculator" << std::endl;
    std::cout << std::string(61, '=') << std::endl;
    std::cout << "Input: " << inputString << std::endl << std::endl;
    
    ShippingCostCalculator calc;
    if (!calc.parseInput(inputString)) {
        std::cerr << "Error parsing input" << std::endl;
        return 1;
    }
    
    // Part 1: Direct route cost
    std::cout << "Part 1: Direct Route Costs" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    std::vector<std::pair<std::string, std::string>> testCases = {
        {"US", "UK"},
        {"US", "CA"},
        {"CA", "UK"},
        {"UK", "US"}  // Should not exist
    };
    
    for (const auto& tc : testCases) {
        int cost;
        std::string carrier;
        if (calc.findDirectCost(tc.first, tc.second, cost, carrier)) {
            std::cout << "  " << tc.first << " -> " << tc.second 
                      << ": Cost = $" << cost << " (Carrier: " << carrier << ")" << std::endl;
        } else {
            std::cout << "  " << tc.first << " -> " << tc.second 
                      << ": No direct route found" << std::endl;
        }
    }
    
    // Part 2: Route with intermediate steps
    std::cout << "\nPart 2: Routes with Intermediate Steps" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    std::vector<std::pair<std::string, std::string>> paths = {
        {"US", "UK"},
        {"US", "CA"},
        {"CA", "UK"}
    };
    
    for (const auto& p : paths) {
        std::cout << "\nFinding cheapest path from " << p.first << " to " << p.second << ":" << std::endl;
        
        Path* bestPath = calc.findPathWithIntermediate(p.first, p.second);
        if (bestPath == nullptr) {
            std::cout << "  Error: No path found" << std::endl;
            continue;
        }
        
        std::cout << "  Cheapest Path (Cost: $" << bestPath->totalCost << "):" << std::endl;
        for (size_t i = 0; i < bestPath->routes.size(); i++) {
            const Route& route = bestPath->routes[i];
            std::cout << "    Step " << (i + 1) << ": " << route.origin << " -> " 
                      << route.destination << " via " << route.carrier 
                      << " (Cost: $" << route.cost << ")" << std::endl;
        }
        
        delete bestPath;
    }
    
    // Example: US -> CA -> UK
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example: US -> UK (comparing direct vs intermediate)" << std::endl;
    std::cout << std::string(60, '=') << std::endl;
    
    int directCost;
    std::string carrier;
    if (calc.findDirectCost("US", "UK", directCost, carrier)) {
        std::cout << "Direct route: $" << directCost << " (Carrier: " << carrier << ")" << std::endl;
    } else {
        std::cout << "Direct route: Not available" << std::endl;
    }
    
    Path* intermediatePath = calc.findPathWithIntermediate("US", "UK");
    if (intermediatePath != nullptr) {
        if (intermediatePath->routes.size() > 1) {
            std::cout << "\nIntermediate route (US -> CA -> UK): $" 
                      << intermediatePath->totalCost << std::endl;
            for (const auto& route : intermediatePath->routes) {
                std::cout << "  " << route.origin << " -> " << route.destination 
                          << ": $" << route.cost << std::endl;
            }
            std::cout << "  Total: $" << intermediatePath->totalCost << std::endl;
        }
        delete intermediatePath;
    }
    
    return 0;
}

