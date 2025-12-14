#include <iostream>
#include <string>
#include <vector>
#include <map>
#include <algorithm>
#include <ctime>
#include <chrono>

struct ServerEvent {
    std::string serverID;
    std::string eventType; // "response_time", "error", etc.
    double value;
    std::chrono::system_clock::time_point timestamp;
    
    ServerEvent(const std::string& sid, const std::string& et, double v, 
                const std::chrono::system_clock::time_point& ts)
        : serverID(sid), eventType(et), value(v), timestamp(ts) {}
};

class ServerPenaltyCalculator {
private:
    std::vector<ServerEvent> events;
    
public:
    // Add an event for a server
    void addEvent(const std::string& serverID, const std::string& eventType, 
                  double value, const std::chrono::system_clock::time_point& timestamp) {
        events.push_back(ServerEvent(serverID, eventType, value, timestamp));
    }
    
    // Calculate penalty for a server based on events
    double calculatePenalty(const std::string& serverID, 
                           double responseTimeThreshold,
                           double responseTimeMultiplier,
                           double errorPenalty) {
        double totalPenalty = 0.0;
        
        for (const auto& event : events) {
            if (event.serverID != serverID) {
                continue;
            }
            
            if (event.eventType == "response_time") {
                if (event.value > responseTimeThreshold) {
                    double penalty = (event.value - responseTimeThreshold) * responseTimeMultiplier;
                    totalPenalty += penalty;
                }
            } else if (event.eventType == "error") {
                totalPenalty += errorPenalty;
            }
        }
        
        return totalPenalty;
    }
    
    // Get penalties for all servers
    std::map<std::string, double> getServerPenalties(double responseTimeThreshold,
                                                     double responseTimeMultiplier,
                                                     double errorPenalty) {
        std::map<std::string, double> penalties;
        std::map<std::string, bool> serverSet;
        
        // Collect all server IDs
        for (const auto& event : events) {
            serverSet[event.serverID] = true;
        }
        
        // Calculate penalty for each server
        for (const auto& pair : serverSet) {
            penalties[pair.first] = calculatePenalty(pair.first, responseTimeThreshold, 
                                                    responseTimeMultiplier, errorPenalty);
        }
        
        return penalties;
    }
    
    // Get servers sorted by penalty (highest first)
    std::vector<std::string> getServersByPenalty(double responseTimeThreshold,
                                                 double responseTimeMultiplier,
                                                 double errorPenalty) {
        auto penalties = getServerPenalties(responseTimeThreshold, responseTimeMultiplier, errorPenalty);
        
        struct ServerPenalty {
            std::string serverID;
            double penalty;
        };
        
        std::vector<ServerPenalty> serverPenalties;
        for (const auto& pair : penalties) {
            serverPenalties.push_back({pair.first, pair.second});
        }
        
        // Sort by penalty (highest first), then by server ID
        std::sort(serverPenalties.begin(), serverPenalties.end(),
                  [](const ServerPenalty& a, const ServerPenalty& b) {
                      if (a.penalty != b.penalty) {
                          return a.penalty > b.penalty;
                      }
                      return a.serverID < b.serverID;
                  });
        
        std::vector<std::string> result;
        for (const auto& sp : serverPenalties) {
            result.push_back(sp.serverID);
        }
        
        return result;
    }
    
    // Get all events for a specific server
    std::vector<ServerEvent> getEventsForServer(const std::string& serverID) {
        std::vector<ServerEvent> result;
        for (const auto& event : events) {
            if (event.serverID == serverID) {
                result.push_back(event);
            }
        }
        return result;
    }
};

int main() {
    std::cout << std::string(60, '=') << std::endl;
    std::cout << "Server Penalty Calculator" << std::endl;
    std::cout << std::string(60, '=') << std::endl;
    
    ServerPenaltyCalculator calc;
    
    // Example 1: Response time penalties
    std::cout << "\nExample 1: Response Time Penalties" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    auto baseTime = std::chrono::system_clock::now();
    calc.addEvent("server1", "response_time", 150.0, baseTime);
    calc.addEvent("server1", "response_time", 120.0, baseTime + std::chrono::seconds(1));
    calc.addEvent("server2", "response_time", 200.0, baseTime);
    calc.addEvent("server2", "response_time", 110.0, baseTime + std::chrono::seconds(1));
    calc.addEvent("server3", "response_time", 90.0, baseTime);
    
    double responseTimeThreshold = 100.0;  // 100ms threshold
    double responseTimeMultiplier = 0.5;   // $0.50 per ms over threshold
    double errorPenalty = 10.0;            // $10 per error
    
    std::cout << "Response Time Threshold: " << responseTimeThreshold << " ms" << std::endl;
    std::cout << "Response Time Multiplier: $" << responseTimeMultiplier 
              << " per ms over threshold" << std::endl;
    std::cout << "Error Penalty: $" << errorPenalty << " per error" << std::endl << std::endl;
    
    auto penalties = calc.getServerPenalties(responseTimeThreshold, responseTimeMultiplier, errorPenalty);
    
    std::cout << "Server Penalties:" << std::endl;
    for (const auto& pair : penalties) {
        std::cout << "  " << pair.first << ": $" << pair.second << std::endl;
    }
    
    // Example 2: Combined response time and errors
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 2: Combined Response Time and Errors" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    ServerPenaltyCalculator calc2;
    calc2.addEvent("serverA", "response_time", 150.0, baseTime);
    calc2.addEvent("serverA", "error", 1.0, baseTime + std::chrono::seconds(1));
    calc2.addEvent("serverA", "error", 1.0, baseTime + std::chrono::seconds(2));
    calc2.addEvent("serverB", "response_time", 120.0, baseTime);
    calc2.addEvent("serverB", "response_time", 130.0, baseTime + std::chrono::seconds(1));
    calc2.addEvent("serverC", "error", 1.0, baseTime);
    calc2.addEvent("serverC", "error", 1.0, baseTime + std::chrono::seconds(1));
    calc2.addEvent("serverC", "error", 1.0, baseTime + std::chrono::seconds(2));
    
    auto penalties2 = calc2.getServerPenalties(responseTimeThreshold, responseTimeMultiplier, errorPenalty);
    
    std::cout << "Server Penalties:" << std::endl;
    for (const auto& pair : penalties2) {
        auto events = calc2.getEventsForServer(pair.first);
        std::cout << "  " << pair.first << ": $" << pair.second
                  << " (events: " << events.size() << ")" << std::endl;
    }
    
    // Example 3: Servers sorted by penalty
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 3: Servers Sorted by Penalty (Highest First)" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    auto sortedServers = calc2.getServersByPenalty(responseTimeThreshold, responseTimeMultiplier, errorPenalty);
    
    std::cout << "Servers sorted by penalty:" << std::endl;
    for (size_t i = 0; i < sortedServers.size(); i++) {
        std::cout << "  " << (i + 1) << ". " << sortedServers[i] 
                  << ": $" << penalties2[sortedServers[i]] << std::endl;
    }
    
    // Example 4: Time-windowed penalties
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 4: Recent Events Only (Last 5 minutes)" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    ServerPenaltyCalculator calc3;
    auto now = std::chrono::system_clock::now();
    calc3.addEvent("serverX", "response_time", 150.0, now - std::chrono::minutes(2));
    calc3.addEvent("serverX", "response_time", 200.0, now - std::chrono::minutes(10));
    calc3.addEvent("serverY", "error", 1.0, now - std::chrono::minutes(1));
    calc3.addEvent("serverY", "error", 1.0, now - std::chrono::minutes(6));
    
    auto timeWindow = std::chrono::minutes(5);
    auto cutoffTime = now - timeWindow;
    
    std::map<std::string, double> penalties3;
    std::map<std::string, bool> serverSet;
    
    for (const auto& event : calc3.getEventsForServer("serverX")) {
        if (event.timestamp > cutoffTime) {
            serverSet["serverX"] = true;
        }
    }
    for (const auto& event : calc3.getEventsForServer("serverY")) {
        if (event.timestamp > cutoffTime) {
            serverSet["serverY"] = true;
        }
    }
    
    for (const auto& pair : serverSet) {
        double penalty = 0.0;
        auto events = calc3.getEventsForServer(pair.first);
        for (const auto& event : events) {
            if (event.timestamp > cutoffTime) {
                if (event.eventType == "response_time") {
                    if (event.value > responseTimeThreshold) {
                        penalty += (event.value - responseTimeThreshold) * responseTimeMultiplier;
                    }
                } else if (event.eventType == "error") {
                    penalty += errorPenalty;
                }
            }
        }
        penalties3[pair.first] = penalty;
    }
    
    std::cout << "Server Penalties (last 5 minutes only):" << std::endl;
    for (const auto& pair : penalties3) {
        std::cout << "  " << pair.first << ": $" << pair.second << std::endl;
    }
    
    return 0;
}

