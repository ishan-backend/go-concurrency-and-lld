#include <iostream>
#include <string>
#include <vector>
#include <map>
#include <climits>
#include <algorithm>

// Record represents a database record (map with string keys and integer values)
using Record = std::map<std::string, int>;

// RecordComparator compares records based on a key and direction
class RecordComparator {
private:
    std::string key;
    std::string direction; // "asc" or "desc"
    
    // Get value for the key, defaulting to 0 if not present
    int getValue(const Record& rec) const {
        auto it = rec.find(key);
        if (it != rec.end()) {
            return it->second;
        }
        return 0;
    }
    
public:
    RecordComparator(const std::string& k, const std::string& dir) 
        : key(k), direction(dir) {}
    
    // Compare two records
    // Returns -1 if rec1 comes before rec2, 0 if equal, 1 if rec1 comes after rec2
    int compare(const Record& rec1, const Record& rec2) const {
        int val1 = getValue(rec1);
        int val2 = getValue(rec2);
        
        if (val1 < val2) {
            return (direction == "asc") ? -1 : 1;
        } else if (val1 > val2) {
            return (direction == "asc") ? 1 : -1;
        }
        return 0;
    }
};

// MinByKey returns the record with the minimum value for the specified key
Record minByKey(const std::string& key, const std::vector<Record>& records) {
    return firstByKey(key, "asc", records);
}

// FirstByKey returns the first record according to key and direction
Record firstByKey(const std::string& key, const std::string& direction, 
                   const std::vector<Record>& records) {
    if (records.empty()) {
        return Record();
    }
    
    RecordComparator comparator(key, direction);
    Record result = records[0];
    
    for (size_t i = 1; i < records.size(); i++) {
        if (comparator.compare(records[i], result) < 0) {
            result = records[i];
        }
    }
    
    return result;
}

// Helper function to print a record
void printRecord(const Record& rec) {
    std::cout << "{";
    bool first = true;
    for (const auto& pair : rec) {
        if (!first) std::cout << ", ";
        std::cout << "\"" << pair.first << "\": " << pair.second;
        first = false;
    }
    std::cout << "}";
}

// Helper function to print records vector
void printRecords(const std::vector<Record>& records) {
    std::cout << "[";
    for (size_t i = 0; i < records.size(); i++) {
        if (i > 0) std::cout << ", ";
        printRecord(records[i]);
    }
    std::cout << "]";
}

int main() {
    std::cout << std::string(60, '=') << std::endl;
    std::cout << "Analytical Database - Part 1: min_by_key" << std::endl;
    std::cout << std::string(60, '=') << std::endl;
    
    // Example 1: Basic min_by_key
    std::cout << "\nExample 1: Basic min_by_key" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    std::vector<Record> records1 = {
        {{"a", 1}, {"b", 2}},
        {{"a", 2}}
    };
    auto result1 = minByKey("a", records1);
    std::cout << "Records: ";
    printRecords(records1);
    std::cout << "\nmin_by_key(\"a\"): ";
    printRecord(result1);
    std::cout << std::endl;
    
    // Example 2: Different order
    std::cout << "\nExample 2: Different order" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    std::vector<Record> records2 = {
        {{"a", 2}},
        {{"a", 1}, {"b", 2}}
    };
    auto result2 = minByKey("a", records2);
    std::cout << "Records: ";
    printRecords(records2);
    std::cout << "\nmin_by_key(\"a\"): ";
    printRecord(result2);
    std::cout << std::endl;
    
    // Example 3: Missing key (treated as 0)
    std::cout << "\nExample 3: Missing key (treated as 0)" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    std::vector<Record> records3 = {
        {{"a", 1}, {"b", 2}},
        {{"a", 2}}
    };
    auto result3 = minByKey("b", records3);
    std::cout << "Records: ";
    printRecords(records3);
    std::cout << "\nmin_by_key(\"b\"): ";
    printRecord(result3);
    std::cout << " (missing key treated as 0)" << std::endl;
    
    // Example 4: Empty record
    std::cout << "\nExample 4: Empty record" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    std::vector<Record> records4 = {{}};
    auto result4 = minByKey("a", records4);
    std::cout << "Records: ";
    printRecords(records4);
    std::cout << "\nmin_by_key(\"a\"): ";
    printRecord(result4);
    std::cout << std::endl;
    
    // Example 5: Negative values
    std::cout << "\nExample 5: Negative values" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    std::vector<Record> records5 = {
        {{"a", -1}},
        {{"b", -1}}
    };
    auto result5 = minByKey("b", records5);
    std::cout << "Records: ";
    printRecords(records5);
    std::cout << "\nmin_by_key(\"b\"): ";
    printRecord(result5);
    std::cout << std::endl;
    
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Analytical Database - Part 2: first_by_key" << std::endl;
    std::cout << std::string(60, '=') << std::endl;
    
    // Example 6: first_by_key ascending
    std::cout << "\nExample 6: first_by_key ascending" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    std::vector<Record> records6 = {{{"a", 1}}};
    auto result6 = firstByKey("a", "asc", records6);
    std::cout << "Records: ";
    printRecords(records6);
    std::cout << "\nfirst_by_key(\"a\", \"asc\"): ";
    printRecord(result6);
    std::cout << std::endl;
    
    // Example 7: first_by_key with missing key
    std::cout << "\nExample 7: first_by_key with missing key" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    std::vector<Record> records7 = {
        {{"b", 1}},
        {{"b", -2}},
        {{"a", 10}}
    };
    auto result7 = firstByKey("a", "asc", records7);
    std::cout << "Records: ";
    printRecords(records7);
    std::cout << "\nfirst_by_key(\"a\", \"asc\"): ";
    printRecord(result7);
    std::cout << " (missing key = 0)" << std::endl;
    
    // Example 8: first_by_key descending
    std::cout << "\nExample 8: first_by_key descending" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    std::vector<Record> records8 = {
        {{"b", 1}},
        {{"b", -2}},
        {{"a", 10}}
    };
    auto result8 = firstByKey("a", "desc", records8);
    std::cout << "Records: ";
    printRecords(records8);
    std::cout << "\nfirst_by_key(\"a\", \"desc\"): ";
    printRecord(result8);
    std::cout << std::endl;
    
    // Example 9: first_by_key on different key
    std::cout << "\nExample 9: first_by_key on different key" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    std::vector<Record> records9 = {
        {{"b", 1}},
        {{"b", -2}},
        {{"a", 10}}
    };
    auto result9a = firstByKey("b", "asc", records9);
    auto result9b = firstByKey("b", "desc", records9);
    std::cout << "Records: ";
    printRecords(records9);
    std::cout << "\nfirst_by_key(\"b\", \"asc\"): ";
    printRecord(result9a);
    std::cout << "\nfirst_by_key(\"b\", \"desc\"): ";
    printRecord(result9b);
    std::cout << std::endl;
    
    // Example 10: first_by_key with empty records
    std::cout << "\nExample 10: first_by_key with empty records" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    std::vector<Record> records10 = {
        {},
        {{"a", 10}, {"b", -10}},
        {},
        {{"a", 3}, {"c", 3}}
    };
    auto result10 = firstByKey("a", "desc", records10);
    std::cout << "Records: ";
    printRecords(records10);
    std::cout << "\nfirst_by_key(\"a\", \"desc\"): ";
    printRecord(result10);
    std::cout << std::endl;
    
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Analytical Database - Part 3: RecordComparator" << std::endl;
    std::cout << std::string(60, '=') << std::endl;
    
    // Example 11: Comparator usage
    std::cout << "\nExample 11: Comparator usage" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    RecordComparator cmp("a", "asc");
    Record rec1 = {{"a", 1}};
    Record rec2 = {{"a", 2}};
    Record rec3 = {{"a", 1}};
    
    std::cout << "Compare ";
    printRecord(rec1);
    std::cout << " vs ";
    printRecord(rec2);
    std::cout << ": " << cmp.compare(rec1, rec2) << std::endl;
    
    std::cout << "Compare ";
    printRecord(rec2);
    std::cout << " vs ";
    printRecord(rec1);
    std::cout << ": " << cmp.compare(rec2, rec1) << std::endl;
    
    std::cout << "Compare ";
    printRecord(rec1);
    std::cout << " vs ";
    printRecord(rec3);
    std::cout << ": " << cmp.compare(rec1, rec3) << std::endl;
    
    // Example 12: Comparator with descending order
    std::cout << "\nExample 12: Comparator with descending order" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    RecordComparator cmpDesc("a", "desc");
    std::cout << "Compare ";
    printRecord(rec1);
    std::cout << " vs ";
    printRecord(rec2);
    std::cout << " (desc): " << cmpDesc.compare(rec1, rec2) << std::endl;
    
    std::cout << "Compare ";
    printRecord(rec2);
    std::cout << " vs ";
    printRecord(rec1);
    std::cout << " (desc): " << cmpDesc.compare(rec2, rec1) << std::endl;
    
    return 0;
}

