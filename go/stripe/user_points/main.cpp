#include <iostream>
#include <string>
#include <vector>
#include <map>
#include <algorithm>
#include <ctime>
#include <chrono>

enum class TransactionType {
    EARN,
    SPEND,
    EXPIRE
};

struct Transaction {
    std::string userID;
    TransactionType type;
    int amount;
    std::chrono::system_clock::time_point timestamp;
    std::string description;
    
    Transaction(const std::string& uid, TransactionType t, int amt, 
                const std::string& desc)
        : userID(uid), type(t), amount(amt), 
          timestamp(std::chrono::system_clock::now()), description(desc) {}
};

struct PointBalance {
    int amount;
    std::chrono::system_clock::time_point expirationTime;
    
    PointBalance(int amt, const std::chrono::system_clock::time_point& exp)
        : amount(amt), expirationTime(exp) {}
};

class UserPointsManager {
private:
    std::vector<Transaction> transactions;
    // userID -> list of point balances with expiration times
    std::map<std::string, std::vector<PointBalance>> userBalances;
    
    // Spend points using FIFO (oldest points first)
    void spendPoints(const std::string& userID, int amount) {
        if (userBalances.find(userID) == userBalances.end()) {
            return;
        }
        
        auto& balance = userBalances[userID];
        int remaining = amount;
        
        for (auto it = balance.begin(); it != balance.end() && remaining > 0;) {
            if (it->amount <= remaining) {
                remaining -= it->amount;
                it = balance.erase(it);
            } else {
                it->amount -= remaining;
                remaining = 0;
            }
        }
    }
    
    // Expire points
    void expirePoints(const std::string& userID, int amount) {
        if (userBalances.find(userID) == userBalances.end()) {
            return;
        }
        
        auto& balance = userBalances[userID];
        auto now = std::chrono::system_clock::now();
        int remaining = amount;
        
        for (auto it = balance.begin(); it != balance.end() && remaining > 0;) {
            if (it->expirationTime <= now) {
                if (it->amount <= remaining) {
                    remaining -= it->amount;
                    it = balance.erase(it);
                } else {
                    it->amount -= remaining;
                    remaining = 0;
                }
            } else {
                ++it;
            }
        }
    }
    
    // Clean expired points for a user
    void cleanExpiredPoints(const std::string& userID) {
        if (userBalances.find(userID) == userBalances.end()) {
            return;
        }
        
        auto& balance = userBalances[userID];
        auto now = std::chrono::system_clock::now();
        
        balance.erase(
            std::remove_if(balance.begin(), balance.end(),
                [now](const PointBalance& pb) {
                    return pb.expirationTime <= now;
                }),
            balance.end()
        );
    }
    
public:
    // Add a transaction and update balances
    void addTransaction(const std::string& userID, TransactionType txnType, 
                       int amount, const std::string& description,
                       const std::chrono::system_clock::time_point* expirationTime = nullptr) {
        transactions.push_back(Transaction(userID, txnType, amount, description));
        
        if (userBalances.find(userID) == userBalances.end()) {
            userBalances[userID] = std::vector<PointBalance>();
        }
        
        switch (txnType) {
            case TransactionType::EARN: {
                auto expTime = expirationTime != nullptr 
                    ? *expirationTime 
                    : std::chrono::system_clock::now() + std::chrono::hours(365 * 24);
                userBalances[userID].push_back(PointBalance(amount, expTime));
                break;
            }
            case TransactionType::SPEND:
                spendPoints(userID, amount);
                break;
            case TransactionType::EXPIRE:
                expirePoints(userID, amount);
                break;
        }
    }
    
    // Get current balance for a user (excluding expired points)
    int getBalance(const std::string& userID) {
        cleanExpiredPoints(userID);
        
        if (userBalances.find(userID) == userBalances.end()) {
            return 0;
        }
        
        int total = 0;
        for (const auto& entry : userBalances[userID]) {
            total += entry.amount;
        }
        return total;
    }
    
    // Get all transactions for a user
    std::vector<Transaction> getTransactions(const std::string& userID) {
        std::vector<Transaction> result;
        for (const auto& txn : transactions) {
            if (txn.userID == userID) {
                result.push_back(txn);
            }
        }
        return result;
    }
    
    // Get balances for all users
    std::map<std::string, int> getAllBalances() {
        std::map<std::string, int> result;
        for (const auto& pair : userBalances) {
            result[pair.first] = getBalance(pair.first);
        }
        return result;
    }
    
    // Get top users sorted by balance (highest first)
    std::vector<std::string> getTopUsers(int limit) {
        auto balances = getAllBalances();
        
        struct UserBalance {
            std::string userID;
            int balance;
        };
        
        std::vector<UserBalance> userBalancesList;
        for (const auto& pair : balances) {
            userBalancesList.push_back({pair.first, pair.second});
        }
        
        // Sort by balance (highest first), then by userID
        std::sort(userBalancesList.begin(), userBalancesList.end(),
                  [](const UserBalance& a, const UserBalance& b) {
                      if (a.balance != b.balance) {
                          return a.balance > b.balance;
                      }
                      return a.userID < b.userID;
                  });
        
        std::vector<std::string> result;
        for (size_t i = 0; i < userBalancesList.size() && i < static_cast<size_t>(limit); i++) {
            result.push_back(userBalancesList[i].userID);
        }
        
        return result;
    }
};

std::string transactionTypeToString(TransactionType type) {
    switch (type) {
        case TransactionType::EARN: return "earn";
        case TransactionType::SPEND: return "spend";
        case TransactionType::EXPIRE: return "expire";
        default: return "unknown";
    }
}

int main() {
    std::cout << std::string(60, '=') << std::endl;
    std::cout << "User Points Manager" << std::endl;
    std::cout << std::string(60, '=') << std::endl;
    
    UserPointsManager manager;
    
    // Example 1: Basic earn and spend
    std::cout << "\nExample 1: Basic Earn and Spend" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    manager.addTransaction("user1", TransactionType::EARN, 100, "Sign up bonus");
    manager.addTransaction("user1", TransactionType::EARN, 50, "Referral bonus");
    manager.addTransaction("user1", TransactionType::SPEND, 75, "Purchase item");
    
    std::cout << "User1 balance: " << manager.getBalance("user1") << " points" << std::endl;
    
    // Example 2: Points expiration
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 2: Points Expiration" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    UserPointsManager manager2;
    auto now = std::chrono::system_clock::now();
    
    auto exp1 = now + std::chrono::hours(30 * 24);
    auto exp2 = now + std::chrono::hours(60 * 24);
    
    manager2.addTransaction("user2", TransactionType::EARN, 100, "Bonus 1", &exp1);
    manager2.addTransaction("user2", TransactionType::EARN, 200, "Bonus 2", &exp2);
    
    std::cout << "User2 balance: " << manager2.getBalance("user2") << " points" << std::endl;
    
    // Example 3: FIFO spending
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 3: FIFO Spending (Oldest Points First)" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    UserPointsManager manager3;
    
    manager3.addTransaction("user3", TransactionType::EARN, 100, "Earn 1");
    manager3.addTransaction("user3", TransactionType::EARN, 50, "Earn 2");
    manager3.addTransaction("user3", TransactionType::EARN, 25, "Earn 3");
    
    std::cout << "User3 balance before spend: " << manager3.getBalance("user3") << " points" << std::endl;
    
    manager3.addTransaction("user3", TransactionType::SPEND, 75, "Spend 75");
    
    std::cout << "User3 balance after spending 75: " << manager3.getBalance("user3") << " points" << std::endl;
    
    // Example 4: Top users
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 4: Top Users by Balance" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    UserPointsManager manager4;
    manager4.addTransaction("alice", TransactionType::EARN, 500, "Bonus");
    manager4.addTransaction("bob", TransactionType::EARN, 300, "Bonus");
    manager4.addTransaction("bob", TransactionType::EARN, 200, "Referral");
    manager4.addTransaction("charlie", TransactionType::EARN, 100, "Bonus");
    manager4.addTransaction("charlie", TransactionType::SPEND, 50, "Purchase");
    
    auto topUsers = manager4.getTopUsers(3);
    std::cout << "Top 3 users by balance:" << std::endl;
    for (size_t i = 0; i < topUsers.size(); i++) {
        int balance = manager4.getBalance(topUsers[i]);
        std::cout << "  " << (i + 1) << ". " << topUsers[i] << ": " << balance << " points" << std::endl;
    }
    
    // Example 5: Transaction history
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 5: Transaction History" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    auto transactions = manager4.getTransactions("bob");
    std::cout << "Transactions for bob:" << std::endl;
    for (size_t i = 0; i < transactions.size(); i++) {
        std::cout << "  " << (i + 1) << ". " << transactionTypeToString(transactions[i].type)
                  << ": " << transactions[i].amount << " points - " 
                  << transactions[i].description << std::endl;
    }
    
    return 0;
}

