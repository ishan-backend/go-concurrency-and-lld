#include <iostream>
#include <string>
#include <vector>
#include <map>
#include <sstream>
#include <algorithm>

class MutualRankCalculator {
private:
    // user -> list of ranked users (in order of preference)
    std::map<std::string, std::vector<std::string>> rankings;
    
    // Trim whitespace
    std::string trim(const std::string& str) {
        size_t first = str.find_first_not_of(" \t");
        if (first == std::string::npos) return "";
        size_t last = str.find_last_not_of(" \t");
        return str.substr(first, (last - first + 1));
    }
    
    // Split string by delimiter
    std::vector<std::string> split(const std::string& str, char delimiter) {
        std::vector<std::string> tokens;
        std::stringstream ss(str);
        std::string token;
        while (std::getline(ss, token, delimiter)) {
            token = trim(token);
            if (!token.empty()) {
                tokens.push_back(token);
            }
        }
        return tokens;
    }
    
public:
    // Add a user's ranking of other users
    void addRanking(const std::string& userID, const std::string& rankingsStr) {
        rankings[userID] = split(rankingsStr, ',');
    }
    
    // Get the rank of targetUser in userID's ranking (1-indexed)
    // Returns 0 if not found
    int getRank(const std::string& userID, const std::string& targetUser) {
        if (rankings.find(userID) == rankings.end()) {
            return 0;
        }
        
        const auto& rankedUsers = rankings[userID];
        for (size_t i = 0; i < rankedUsers.size(); i++) {
            if (rankedUsers[i] == targetUser) {
                return i + 1; // 1-indexed
            }
        }
        return 0;
    }
    
    // Get mutual rank between two users
    // Mutual rank is the sum of ranks: rank of userB in userA's list + rank of userA in userB's list
    bool getMutualRank(const std::string& userA, const std::string& userB, int& mutualRank) {
        int rankA = getRank(userA, userB);
        int rankB = getRank(userB, userA);
        
        if (rankA == 0 || rankB == 0) {
            return false; // One or both users don't have the other in their ranking
        }
        
        mutualRank = rankA + rankB;
        return true;
    }
    
    // Find best mutual matches for a given user
    std::vector<std::string> findBestMatches(const std::string& userID) {
        struct Match {
            std::string userID;
            int mutualRank;
        };
        
        std::vector<Match> matches;
        
        // Get all users that this user has ranked
        if (rankings.find(userID) == rankings.end()) {
            return {};
        }
        
        const auto& rankedUsers = rankings[userID];
        
        // Check mutual rank with each ranked user
        for (const auto& otherUser : rankedUsers) {
            int mutualRank;
            if (getMutualRank(userID, otherUser, mutualRank)) {
                matches.push_back({otherUser, mutualRank});
            }
        }
        
        // Sort by mutual rank (lower is better), then by userID for stability
        std::sort(matches.begin(), matches.end(), [](const Match& a, const Match& b) {
            if (a.mutualRank != b.mutualRank) {
                return a.mutualRank < b.mutualRank;
            }
            return a.userID < b.userID;
        });
        
        std::vector<std::string> result;
        for (const auto& m : matches) {
            result.push_back(m.userID);
        }
        
        return result;
    }
    
    // Get all mutual ranks for a given user
    std::map<std::string, int> getAllMutualRanks(const std::string& userID) {
        std::map<std::string, int> result;
        
        if (rankings.find(userID) == rankings.end()) {
            return result;
        }
        
        const auto& rankedUsers = rankings[userID];
        for (const auto& otherUser : rankedUsers) {
            int mutualRank;
            if (getMutualRank(userID, otherUser, mutualRank)) {
                result[otherUser] = mutualRank;
            }
        }
        
        return result;
    }
};

int main() {
    std::cout << std::string(60, '=') << std::endl;
    std::cout << "Mutual Rank Calculator" << std::endl;
    std::cout << std::string(60, '=') << std::endl;
    
    MutualRankCalculator calc;
    
    // Example 1: Basic mutual rankings
    std::cout << "\nExample 1: Basic Mutual Rankings" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    calc.addRanking("alice", "bob,charlie,dave");
    calc.addRanking("bob", "alice,charlie,dave");
    calc.addRanking("charlie", "alice,bob,dave");
    calc.addRanking("dave", "alice,bob,charlie");
    
    std::cout << "Rankings:" << std::endl;
    std::cout << "  alice: bob, charlie, dave" << std::endl;
    std::cout << "  bob: alice, charlie, dave" << std::endl;
    std::cout << "  charlie: alice, bob, dave" << std::endl;
    std::cout << "  dave: alice, bob, charlie" << std::endl;
    std::cout << std::endl;
    
    // Check mutual ranks
    std::vector<std::pair<std::string, std::string>> testPairs = {
        {"alice", "bob"},
        {"alice", "charlie"},
        {"bob", "charlie"}
    };
    
    for (const auto& pair : testPairs) {
        const std::string& userA = pair.first;
        const std::string& userB = pair.second;
        int rankA = calc.getRank(userA, userB);
        int rankB = calc.getRank(userB, userA);
        int mutualRank;
        
        if (calc.getMutualRank(userA, userB, mutualRank)) {
            std::cout << "  " << userA << " <-> " << userB 
                      << ": rank(" << userA << "->" << userB << ")=" << rankA
                      << ", rank(" << userB << "->" << userA << ")=" << rankB
                      << ", mutual rank=" << mutualRank << std::endl;
        } else {
            std::cout << "  " << userA << " <-> " << userB 
                      << ": No mutual ranking" << std::endl;
        }
    }
    
    // Find best matches for alice
    std::cout << "\nBest matches for alice:" << std::endl;
    std::vector<std::string> bestMatches = calc.findBestMatches("alice");
    for (size_t i = 0; i < bestMatches.size(); i++) {
        int mutualRank;
        calc.getMutualRank("alice", bestMatches[i], mutualRank);
        std::cout << "  " << (i + 1) << ". " << bestMatches[i] 
                  << " (mutual rank: " << mutualRank << ")" << std::endl;
    }
    
    // Example 2: Asymmetric rankings
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 2: Asymmetric Rankings" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    MutualRankCalculator calc2;
    calc2.addRanking("user1", "user2,user3,user4");
    calc2.addRanking("user2", "user4,user1,user3");
    calc2.addRanking("user3", "user1,user2,user4");
    calc2.addRanking("user4", "user2,user1,user3");
    
    std::cout << "Rankings:" << std::endl;
    std::cout << "  user1: user2, user3, user4" << std::endl;
    std::cout << "  user2: user4, user1, user3" << std::endl;
    std::cout << "  user3: user1, user2, user4" << std::endl;
    std::cout << "  user4: user2, user1, user3" << std::endl;
    std::cout << std::endl;
    
    // Check mutual ranks
    std::vector<std::pair<std::string, std::string>> testPairs2 = {
        {"user1", "user2"},
        {"user1", "user3"},
        {"user2", "user4"}
    };
    
    for (const auto& pair : testPairs2) {
        const std::string& userA = pair.first;
        const std::string& userB = pair.second;
        int rankA = calc2.getRank(userA, userB);
        int rankB = calc2.getRank(userB, userA);
        int mutualRank;
        
        if (calc2.getMutualRank(userA, userB, mutualRank)) {
            std::cout << "  " << userA << " <-> " << userB 
                      << ": rank(" << userA << "->" << userB << ")=" << rankA
                      << ", rank(" << userB << "->" << userA << ")=" << rankB
                      << ", mutual rank=" << mutualRank << std::endl;
        }
    }
    
    // Find best matches for user1
    std::cout << "\nBest matches for user1:" << std::endl;
    std::vector<std::string> bestMatches2 = calc2.findBestMatches("user1");
    for (size_t i = 0; i < bestMatches2.size(); i++) {
        int mutualRank;
        calc2.getMutualRank("user1", bestMatches2[i], mutualRank);
        std::cout << "  " << (i + 1) << ". " << bestMatches2[i] 
                  << " (mutual rank: " << mutualRank << ")" << std::endl;
    }
    
    // Example 3: One-way ranking (no mutual rank)
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 3: One-Way Ranking" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    MutualRankCalculator calc3;
    calc3.addRanking("A", "B,C");
    calc3.addRanking("B", "C,D");
    // A ranks B, but B doesn't rank A
    
    int mutualRank;
    if (calc3.getMutualRank("A", "B", mutualRank)) {
        std::cout << "A <-> B: mutual rank = " << mutualRank << std::endl;
    } else {
        std::cout << "A <-> B: No mutual ranking (B doesn't rank A)" << std::endl;
    }
    
    return 0;
}

