#include <iostream>
#include <string>
#include <vector>
#include <sstream>

// Compress a single word using run-length encoding
// e.g., "stripe" -> "s4e" (s + 4 middle chars + e)
std::string compressWord(const std::string& word) {
    if (word.length() <= 2) {
        return word;
    }
    // Keep first and last char, replace middle with count
    int middle = word.length() - 2;
    return word.substr(0, 1) + std::to_string(middle) + word.substr(word.length() - 1);
}

// Split string by delimiter
std::vector<std::string> split(const std::string& str, char delimiter) {
    std::vector<std::string> tokens;
    std::stringstream ss(str);
    std::string token;
    while (std::getline(ss, token, delimiter)) {
        tokens.push_back(token);
    }
    return tokens;
}

// Join vector of strings with delimiter
std::string join(const std::vector<std::string>& parts, const std::string& delimiter) {
    if (parts.empty()) return "";
    std::string result = parts[0];
    for (size_t i = 1; i < parts.size(); i++) {
        result += delimiter + parts[i];
    }
    return result;
}

// Compress minor parts (separated by '.')
// minor_parts indicates how many parts from the end to compress
std::string compressMinorParts(const std::string& major, int minorParts) {
    std::vector<std::string> parts = split(major, '.');
    
    // Determine starting index for compression
    // We compress from the end, minorParts at a time
    int startIdx = parts.size() - minorParts;
    if (startIdx < 0) {
        startIdx = 0;
    }
    
    // Compress parts from startIdx onwards
    for (size_t i = startIdx; i < parts.size(); i++) {
        parts[i] = compressWord(parts[i]);
    }
    
    return join(parts, ".");
}

// Compress major parts (separated by '/')
// minor_parts indicates how many minor parts to compress in each major part
std::string compressString(const std::string& str, int minorParts) {
    std::vector<std::string> majorParts = split(str, '/');
    
    for (size_t i = 0; i < majorParts.size(); i++) {
        majorParts[i] = compressMinorParts(majorParts[i], minorParts);
    }
    
    return join(majorParts, "/");
}

// Perform compression in two phases
std::string performCompression(const std::string& str, int minorParts) {
    // Phase 1: Compress the first major part
    std::cout << "Original string: " << str << std::endl;
    std::cout << "Minor parts to compress: " << minorParts << std::endl << std::endl;
    
    std::string phase1 = compressString(str, minorParts);
    std::cout << "After Phase 1 (compress first major part): " << phase1 << std::endl;
    
    // Phase 2: Compress all major parts
    std::vector<std::string> majorParts = split(phase1, '/');
    for (size_t i = 0; i < majorParts.size(); i++) {
        majorParts[i] = compressMinorParts(majorParts[i], minorParts);
    }
    std::string phase2 = join(majorParts, "/");
    std::cout << "After Phase 2 (compress all major parts): " << phase2 << std::endl << std::endl;
    
    return phase2;
}

int main() {
    std::cout << std::string(60, '=') << std::endl;
    std::cout << "Example 1:" << std::endl;
    std::cout << std::string(60, '=') << std::endl;
    std::string str1 = "stripe.com/payments/checkout/customer.john.doe";
    std::string result1 = performCompression(str1, 2);
    std::cout << "Final result: " << result1 << std::endl << std::endl;
    
    std::cout << std::string(60, '=') << std::endl;
    std::cout << "Example 2:" << std::endl;
    std::cout << std::string(60, '=') << std::endl;
    std::string str2 = "www.api.stripe.com/checkout";
    std::string result2 = performCompression(str2, 3);
    std::cout << "Final result: " << result2 << std::endl << std::endl;
    
    // Additional test case
    std::cout << std::string(60, '=') << std::endl;
    std::cout << "Test Case 3 (single char words):" << std::endl;
    std::cout << std::string(60, '=') << std::endl;
    std::string str3 = "a.b.c/d.e.f";
    std::string result3 = performCompression(str3, 2);
    std::cout << "Final result: " << result3 << std::endl;
    
    return 0;
}

