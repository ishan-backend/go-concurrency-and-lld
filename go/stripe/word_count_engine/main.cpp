#include <iostream>
#include <string>
#include <vector>
#include <map>
#include <algorithm>
#include <sstream>
#include <cctype>

struct WordCount {
    std::string word;
    int count;
    
    WordCount(const std::string& w, int c) : word(w), count(c) {}
};

class WordCountEngine {
private:
    std::map<std::string, int> wordCounts;
    
    // Convert string to lowercase
    std::string toLower(const std::string& str) {
        std::string result = str;
        std::transform(result.begin(), result.end(), result.begin(), ::tolower);
        return result;
    }
    
    // Remove punctuation from word
    std::string cleanWord(const std::string& word) {
        std::string cleaned;
        for (char c : word) {
            if (std::isalnum(c)) {
                cleaned += c;
            }
        }
        return cleaned;
    }
    
public:
    // Process a document and count words
    void processDocument(const std::string& document) {
        std::istringstream iss(document);
        std::string word;
        
        while (iss >> word) {
            std::string cleaned = cleanWord(toLower(word));
            if (!cleaned.empty()) {
                wordCounts[cleaned]++;
            }
        }
    }
    
    // Get word counts sorted by frequency (descending), then alphabetically
    std::vector<WordCount> getWordCounts() {
        std::vector<WordCount> result;
        
        for (const auto& pair : wordCounts) {
            result.push_back(WordCount(pair.first, pair.second));
        }
        
        // Sort by count (descending), then by word (ascending)
        std::sort(result.begin(), result.end(),
                  [](const WordCount& a, const WordCount& b) {
                      if (a.count != b.count) {
                          return a.count > b.count;
                      }
                      return a.word < b.word;
                  });
        
        return result;
    }
    
    // Get top N words by frequency
    std::vector<WordCount> getTopWords(int n) {
        auto allCounts = getWordCounts();
        
        if (n > static_cast<int>(allCounts.size())) {
            n = allCounts.size();
        }
        
        return std::vector<WordCount>(allCounts.begin(), allCounts.begin() + n);
    }
    
    // Get frequency of a specific word
    int getWordFrequency(const std::string& word) {
        std::string lowerWord = toLower(word);
        auto it = wordCounts.find(lowerWord);
        if (it != wordCounts.end()) {
            return it->second;
        }
        return 0;
    }
    
    // Reset all word counts
    void reset() {
        wordCounts.clear();
    }
};

int main() {
    std::cout << std::string(60, '=') << std::endl;
    std::cout << "Word Count Engine" << std::endl;
    std::cout << std::string(60, '=') << std::endl;
    
    // Example 1: Basic word counting
    std::cout << "\nExample 1: Basic Word Counting" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    WordCountEngine engine;
    std::string document = "Practice makes perfect. you'll only get Perfect by practice. just practice!";
    engine.processDocument(document);
    
    std::cout << "Document: " << document << "\n\n";
    
    auto wordCounts = engine.getWordCounts();
    std::cout << "Word Counts:" << std::endl;
    for (const auto& wc : wordCounts) {
        std::cout << "  " << wc.word << ": " << wc.count << std::endl;
    }
    
    // Example 2: Top words
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 2: Top Words" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    auto topWords = engine.getTopWords(3);
    std::cout << "Top 3 words:" << std::endl;
    for (size_t i = 0; i < topWords.size(); i++) {
        std::cout << "  " << (i + 1) << ". " << topWords[i].word 
                  << ": " << topWords[i].count << std::endl;
    }
    
    // Example 3: Word frequency lookup
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 3: Word Frequency Lookup" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    std::vector<std::string> testWords = {"practice", "perfect", "makes", "nonexistent"};
    for (const auto& word : testWords) {
        int freq = engine.getWordFrequency(word);
        std::cout << "  '" << word << "': " << freq << std::endl;
    }
    
    // Example 4: Multiple documents
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 4: Multiple Documents" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    WordCountEngine engine2;
    std::vector<std::string> documents = {
        "The quick brown fox jumps over the lazy dog",
        "The lazy dog was jumped over by the quick brown fox",
        "A quick brown fox is always faster than a lazy dog"
    };
    
    for (size_t i = 0; i < documents.size(); i++) {
        std::cout << "Processing document " << (i + 1) << ": " << documents[i] << std::endl;
        engine2.processDocument(documents[i]);
    }
    
    std::cout << "\nCombined Word Counts:" << std::endl;
    auto wordCounts2 = engine2.getWordCounts();
    for (const auto& wc : wordCounts2) {
        std::cout << "  " << wc.word << ": " << wc.count << std::endl;
    }
    
    // Example 5: Case-insensitive and punctuation handling
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 5: Case-Insensitive and Punctuation Handling" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    WordCountEngine engine3;
    std::string doc3 = "Hello, World! HELLO world. hello-WORLD!";
    engine3.processDocument(doc3);
    
    std::cout << "Document: " << doc3 << std::endl;
    std::cout << "\nWord Counts (case-insensitive, punctuation removed):" << std::endl;
    auto wordCounts3 = engine3.getWordCounts();
    for (const auto& wc : wordCounts3) {
        std::cout << "  " << wc.word << ": " << wc.count << std::endl;
    }
    
    return 0;
}

