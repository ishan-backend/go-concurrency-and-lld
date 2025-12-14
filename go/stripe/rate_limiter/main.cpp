#include <iostream>
#include <string>
#include <vector>
#include <map>
#include <mutex>
#include <chrono>
#include <algorithm>

class RateLimiter {
private:
    std::mutex mu;
    std::map<std::string, std::vector<std::chrono::system_clock::time_point>> requests;
    int maxRequests;
    std::chrono::duration<double> windowDuration;
    
public:
    RateLimiter(int maxReq, std::chrono::duration<double> windowDur)
        : maxRequests(maxReq), windowDuration(windowDur) {}
    
    // Check if a request is allowed
    bool isAllowed(const std::string& userID) {
        std::lock_guard<std::mutex> lock(mu);
        
        auto now = std::chrono::system_clock::now();
        auto windowStart = now - windowDuration;
        
        // Get or create request history for user
        auto& userRequests = requests[userID];
        
        // Remove requests outside the window
        userRequests.erase(
            std::remove_if(userRequests.begin(), userRequests.end(),
                [windowStart](const auto& reqTime) {
                    return reqTime <= windowStart;
                }),
            userRequests.end()
        );
        
        // Check if we're at the limit
        if (userRequests.size() >= static_cast<size_t>(maxRequests)) {
            return false;
        }
        
        // Add current request
        userRequests.push_back(now);
        return true;
    }
    
    // Get remaining requests for a user
    int getRemainingRequests(const std::string& userID) {
        std::lock_guard<std::mutex> lock(mu);
        
        auto now = std::chrono::system_clock::now();
        auto windowStart = now - windowDuration;
        
        auto it = requests.find(userID);
        if (it == requests.end()) {
            return maxRequests;
        }
        
        // Count valid requests
        int count = 0;
        for (const auto& reqTime : it->second) {
            if (reqTime > windowStart) {
                count++;
            }
        }
        
        return maxRequests - count;
    }
    
    // Reset all rate limit data
    void reset() {
        std::lock_guard<std::mutex> lock(mu);
        requests.clear();
    }
};

class TokenBucket {
public:
    double tokens;
    std::chrono::system_clock::time_point lastRefill;
    
    TokenBucket(double maxTokens) 
        : tokens(maxTokens), lastRefill(std::chrono::system_clock::now()) {}
};

class TokenBucketRateLimiter {
private:
    std::mutex mu;
    std::map<std::string, TokenBucket> buckets;
    int maxTokens;
    double refillRate; // tokens per second
    std::chrono::duration<double> refillInterval;
    
public:
    TokenBucketRateLimiter(int maxTok, double refill, std::chrono::duration<double> interval)
        : maxTokens(maxTok), refillRate(refill), refillInterval(interval) {}
    
    // Check if a request is allowed using token bucket
    bool isAllowed(const std::string& userID) {
        std::lock_guard<std::mutex> lock(mu);
        
        auto now = std::chrono::system_clock::now();
        
        // Get or create bucket for user
        auto it = buckets.find(userID);
        if (it == buckets.end()) {
            buckets[userID] = TokenBucket(maxTokens);
            return true;
        }
        
        auto& bucket = it->second;
        
        // Refill tokens based on time elapsed
        auto elapsed = std::chrono::duration<double>(now - bucket.lastRefill).count();
        double tokensToAdd = elapsed * refillRate;
        bucket.tokens = std::min(bucket.tokens + tokensToAdd, static_cast<double>(maxTokens));
        bucket.lastRefill = now;
        
        // Check if we have tokens
        if (bucket.tokens >= 1.0) {
            bucket.tokens -= 1.0;
            return true;
        }
        
        return false;
    }
    
    // Get remaining tokens for a user
    int getRemainingTokens(const std::string& userID) {
        std::lock_guard<std::mutex> lock(mu);
        
        auto it = buckets.find(userID);
        if (it == buckets.end()) {
            return maxTokens;
        }
        
        auto& bucket = it->second;
        auto now = std::chrono::system_clock::now();
        auto elapsed = std::chrono::duration<double>(now - bucket.lastRefill).count();
        double tokensToAdd = elapsed * refillRate;
        double currentTokens = std::min(bucket.tokens + tokensToAdd, static_cast<double>(maxTokens));
        
        return static_cast<int>(currentTokens);
    }
};

int main() {
    std::cout << std::string(60, '=') << std::endl;
    std::cout << "Rate Limiter - Sliding Window" << std::endl;
    std::cout << std::string(60, '=') << std::endl;
    
    // Example 1: Sliding Window Rate Limiter
    std::cout << "\nExample 1: Sliding Window Rate Limiter" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    std::cout << "Limit: 5 requests per 10 seconds" << std::endl;
    
    RateLimiter limiter(5, std::chrono::seconds(10));
    
    std::string userID = "user1";
    for (int i = 0; i < 7; i++) {
        bool allowed = limiter.isAllowed(userID);
        int remaining = limiter.getRemainingRequests(userID);
        std::cout << "Request " << (i + 1) << ": Allowed=" << (allowed ? "true" : "false")
                  << ", Remaining=" << remaining << std::endl;
        std::this_thread::sleep_for(std::chrono::milliseconds(500));
    }
    
    // Example 2: Rate limit exceeded
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 2: Rate Limit Exceeded" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    std::cout << "Limit: 3 requests per 5 seconds" << std::endl;
    
    RateLimiter limiter2(3, std::chrono::seconds(5));
    std::string userID2 = "user2";
    
    for (int i = 0; i < 5; i++) {
        bool allowed = limiter2.isAllowed(userID2);
        int remaining = limiter2.getRemainingRequests(userID2);
        std::cout << "Request " << (i + 1) << ": Allowed=" << (allowed ? "true" : "false")
                  << ", Remaining=" << remaining << std::endl;
    }
    
    // Example 3: Multiple users
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 3: Multiple Users" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    std::cout << "Limit: 2 requests per 3 seconds per user" << std::endl;
    
    RateLimiter limiter3(2, std::chrono::seconds(3));
    
    std::vector<std::string> users = {"alice", "bob", "alice", "bob", "alice"};
    for (size_t i = 0; i < users.size(); i++) {
        bool allowed = limiter3.isAllowed(users[i]);
        int remaining = limiter3.getRemainingRequests(users[i]);
        std::cout << "User " << users[i] << ", Request " << (i + 1)
                  << ": Allowed=" << (allowed ? "true" : "false")
                  << ", Remaining=" << remaining << std::endl;
    }
    
    // Example 4: Token Bucket Rate Limiter
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Rate Limiter - Token Bucket" << std::endl;
    std::cout << std::string(60, '=') << std::endl;
    std::cout << "\nExample 4: Token Bucket Rate Limiter" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    std::cout << "Max tokens: 5, Refill rate: 1 token per second" << std::endl;
    
    TokenBucketRateLimiter tokenLimiter(5, 1.0, std::chrono::seconds(1));
    std::string userID3 = "user3";
    
    for (int i = 0; i < 7; i++) {
        bool allowed = tokenLimiter.isAllowed(userID3);
        int remaining = tokenLimiter.getRemainingTokens(userID3);
        std::cout << "Request " << (i + 1) << ": Allowed=" << (allowed ? "true" : "false")
                  << ", Remaining tokens=" << remaining << std::endl;
        if (i < 4) {
            std::this_thread::sleep_for(std::chrono::milliseconds(200));
        } else {
            std::this_thread::sleep_for(std::chrono::seconds(1));
        }
    }
    
    // Example 5: Token bucket with slow requests
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 5: Token Bucket - Slow Requests" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    std::cout << "Max tokens: 3, Refill rate: 0.5 tokens per second" << std::endl;
    
    TokenBucketRateLimiter tokenLimiter2(3, 0.5, std::chrono::seconds(1));
    std::string userID4 = "user4";
    
    std::cout << "Making 3 requests quickly, then waiting for refill..." << std::endl;
    for (int i = 0; i < 3; i++) {
        bool allowed = tokenLimiter2.isAllowed(userID4);
        int remaining = tokenLimiter2.getRemainingTokens(userID4);
        std::cout << "Request " << (i + 1) << ": Allowed=" << (allowed ? "true" : "false")
                  << ", Remaining tokens=" << remaining << std::endl;
    }
    
    std::cout << "\nWaiting 3 seconds for token refill..." << std::endl;
    std::this_thread::sleep_for(std::chrono::seconds(3));
    
    bool allowed = tokenLimiter2.isAllowed(userID4);
    int remaining = tokenLimiter2.getRemainingTokens(userID4);
    std::cout << "Request after wait: Allowed=" << (allowed ? "true" : "false")
              << ", Remaining tokens=" << remaining << std::endl;
    
    return 0;
}

