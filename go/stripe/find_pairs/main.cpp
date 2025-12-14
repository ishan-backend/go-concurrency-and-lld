#include <iostream>
#include <vector>
#include <map>
#include <algorithm>
#include <string>
#include <sstream>

struct Pair {
    int first;
    int second;
    
    Pair(int f, int s) : first(f), second(s) {}
};

class FindPairs {
private:
    std::vector<int> numbers;
    
public:
    FindPairs(const std::vector<int>& nums) : numbers(nums) {}
    
    // Find all pairs of indices whose values sum to target
    std::vector<Pair> findPairsWithSum(int target) {
        std::vector<Pair> pairs;
        std::map<int, std::vector<int>> seen; // value -> list of indices
        
        // First pass: build index map
        for (size_t i = 0; i < numbers.size(); i++) {
            seen[numbers[i]].push_back(i);
        }
        
        // Second pass: find pairs
        std::map<std::string, bool> usedPairs;
        
        for (size_t i = 0; i < numbers.size(); i++) {
            int complement = target - numbers[i];
            
            if (seen.find(complement) != seen.end()) {
                for (int j : seen[complement]) {
                    if (static_cast<int>(i) < j) { // Ensure i < j to avoid duplicates
                        std::ostringstream oss;
                        oss << i << "," << j;
                        std::string key = oss.str();
                        
                        if (usedPairs.find(key) == usedPairs.end()) {
                            pairs.push_back(Pair(static_cast<int>(i), j));
                            usedPairs[key] = true;
                        }
                    }
                }
            }
        }
        
        return pairs;
    }
    
    // Find pairs sorted by first index, then second index
    std::vector<Pair> findPairsWithSumSorted(int target) {
        auto pairs = findPairsWithSum(target);
        
        std::sort(pairs.begin(), pairs.end(),
                  [](const Pair& a, const Pair& b) {
                      if (a.first != b.first) {
                          return a.first < b.first;
                      }
                      return a.second < b.second;
                  });
        
        return pairs;
    }
    
    // Count the number of pairs that sum to target
    int countPairs(int target) {
        return findPairsWithSum(target).size();
    }
    
    // Find pairs using two-pointer technique (requires sorted array)
    std::vector<Pair> findPairsWithSumTwoPointer(int target) {
        // Create a copy and sort with indices
        struct IndexedNum {
            int value;
            int index;
        };
        
        std::vector<IndexedNum> indexed;
        for (size_t i = 0; i < numbers.size(); i++) {
            indexed.push_back({numbers[i], static_cast<int>(i)});
        }
        
        std::sort(indexed.begin(), indexed.end(),
                  [](const IndexedNum& a, const IndexedNum& b) {
                      return a.value < b.value;
                  });
        
        std::vector<Pair> pairs;
        int left = 0;
        int right = indexed.size() - 1;
        
        while (left < right) {
            int sum = indexed[left].value + indexed[right].value;
            
            if (sum == target) {
                // Add pair (ensure smaller index first)
                int firstIdx = indexed[left].index;
                int secondIdx = indexed[right].index;
                if (firstIdx < secondIdx) {
                    pairs.push_back(Pair(firstIdx, secondIdx));
                } else {
                    pairs.push_back(Pair(secondIdx, firstIdx));
                }
                
                // Move pointers to find more pairs
                left++;
                right--;
                
                // Skip duplicates
                while (left < right && indexed[left].value == indexed[left - 1].value) {
                    left++;
                }
                while (left < right && indexed[right].value == indexed[right + 1].value) {
                    right--;
                }
            } else if (sum < target) {
                left++;
            } else {
                right--;
            }
        }
        
        return pairs;
    }
};

int main() {
    std::cout << std::string(60, '=') << std::endl;
    std::cout << "Find Pairs - Two Sum Problem" << std::endl;
    std::cout << std::string(60, '=') << std::endl;
    
    // Example 1: Basic two sum
    std::cout << "\nExample 1: Basic Two Sum" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    std::vector<int> numbers1 = {2, 7, 11, 15};
    int target1 = 9;
    FindPairs fp1(numbers1);
    
    std::cout << "Numbers: ";
    for (int n : numbers1) std::cout << n << " ";
    std::cout << "\nTarget: " << target1 << "\n\n";
    
    auto pairs1 = fp1.findPairsWithSumSorted(target1);
    std::cout << "Pairs (indices):" << std::endl;
    for (const auto& pair : pairs1) {
        std::cout << "  (" << pair.first << ", " << pair.second 
                  << ") -> values: [" << numbers1[pair.first] << ", " 
                  << numbers1[pair.second] << "] = " 
                  << (numbers1[pair.first] + numbers1[pair.second]) << std::endl;
    }
    
    // Example 2: Multiple pairs
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 2: Multiple Pairs" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    std::vector<int> numbers2 = {1, 2, 3, 4, 5, 6};
    int target2 = 7;
    FindPairs fp2(numbers2);
    
    std::cout << "Numbers: ";
    for (int n : numbers2) std::cout << n << " ";
    std::cout << "\nTarget: " << target2 << "\n\n";
    
    auto pairs2 = fp2.findPairsWithSumSorted(target2);
    std::cout << "Pairs (indices):" << std::endl;
    for (const auto& pair : pairs2) {
        std::cout << "  (" << pair.first << ", " << pair.second 
                  << ") -> values: [" << numbers2[pair.first] << ", " 
                  << numbers2[pair.second] << "] = " 
                  << (numbers2[pair.first] + numbers2[pair.second]) << std::endl;
    }
    
    std::cout << "\nTotal pairs: " << fp2.countPairs(target2) << std::endl;
    
    // Example 3: Duplicate values
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 3: Duplicate Values" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    std::vector<int> numbers3 = {3, 3, 4, 5, 6, 3};
    int target3 = 6;
    FindPairs fp3(numbers3);
    
    std::cout << "Numbers: ";
    for (int n : numbers3) std::cout << n << " ";
    std::cout << "\nTarget: " << target3 << "\n\n";
    
    auto pairs3 = fp3.findPairsWithSumSorted(target3);
    std::cout << "Pairs (indices):" << std::endl;
    for (const auto& pair : pairs3) {
        std::cout << "  (" << pair.first << ", " << pair.second 
                  << ") -> values: [" << numbers3[pair.first] << ", " 
                  << numbers3[pair.second] << "] = " 
                  << (numbers3[pair.first] + numbers3[pair.second]) << std::endl;
    }
    
    // Example 4: Two-pointer approach
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 4: Two-Pointer Approach" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    std::vector<int> numbers4 = {1, 3, 4, 5, 6, 8, 9};
    int target4 = 10;
    FindPairs fp4(numbers4);
    
    std::cout << "Numbers: ";
    for (int n : numbers4) std::cout << n << " ";
    std::cout << "\nTarget: " << target4 << "\n\n";
    
    auto pairs4 = fp4.findPairsWithSumTwoPointer(target4);
    std::cout << "Pairs using two-pointer (indices):" << std::endl;
    for (const auto& pair : pairs4) {
        std::cout << "  (" << pair.first << ", " << pair.second 
                  << ") -> values: [" << numbers4[pair.first] << ", " 
                  << numbers4[pair.second] << "] = " 
                  << (numbers4[pair.first] + numbers4[pair.second]) << std::endl;
    }
    
    // Example 5: No pairs found
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 5: No Pairs Found" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    std::vector<int> numbers5 = {1, 2, 3, 4};
    int target5 = 10;
    FindPairs fp5(numbers5);
    
    std::cout << "Numbers: ";
    for (int n : numbers5) std::cout << n << " ";
    std::cout << "\nTarget: " << target5 << "\n\n";
    
    auto pairs5 = fp5.findPairsWithSum(target5);
    if (pairs5.empty()) {
        std::cout << "No pairs found" << std::endl;
    } else {
        std::cout << "Pairs found:" << std::endl;
        for (const auto& pair : pairs5) {
            std::cout << "  (" << pair.first << ", " << pair.second << ")" << std::endl;
        }
    }
    
    return 0;
}

