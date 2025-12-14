#include <iostream>
#include <string>
#include <vector>
#include <sstream>
#include <algorithm>
#include <regex>

// Charge represents a payment charge
struct Charge {
    std::string cardCountry;
    std::string currency;
    int amount;
    std::string ipCountry;
    
    Charge() : amount(0) {}
};

// ParseCharge parses a charge string
Charge parseCharge(const std::string& chargeStr) {
    Charge charge;
    
    std::string str = chargeStr;
    // Remove "CHARGE: " prefix
    if (str.find("CHARGE: ") == 0) {
        str = str.substr(8);
    }
    
    // Split by &
    std::istringstream iss(str);
    std::string token;
    while (std::getline(iss, token, '&')) {
        size_t eqPos = token.find('=');
        if (eqPos != std::string::npos) {
            std::string key = token.substr(0, eqPos);
            std::string value = token.substr(eqPos + 1);
            
            // Trim whitespace
            key.erase(0, key.find_first_not_of(" \t"));
            key.erase(key.find_last_not_of(" \t") + 1);
            value.erase(0, value.find_first_not_of(" \t"));
            value.erase(value.find_last_not_of(" \t") + 1);
            
            if (key == "card_country") {
                charge.cardCountry = value;
            } else if (key == "currency") {
                charge.currency = value;
            } else if (key == "amount") {
                charge.amount = std::stoi(value);
            } else if (key == "ip_country") {
                charge.ipCountry = value;
            }
        }
    }
    
    return charge;
}

// Rule represents a rule (ALLOW or BLOCK)
struct Rule {
    std::string type;      // "ALLOW" or "BLOCK"
    std::string condition; // e.g., "amount<100"
};

// ParseRule parses a rule string
Rule parseRule(const std::string& ruleStr) {
    Rule rule;
    
    if (ruleStr.find("ALLOW:") == 0) {
        rule.type = "ALLOW";
        rule.condition = ruleStr.substr(6);
    } else if (ruleStr.find("BLOCK:") == 0) {
        rule.type = "BLOCK";
        rule.condition = ruleStr.substr(6);
    }
    
    // Trim whitespace
    rule.condition.erase(0, rule.condition.find_first_not_of(" \t"));
    rule.condition.erase(rule.condition.find_last_not_of(" \t") + 1);
    
    return rule;
}

// Compare integers
bool compareInt(int a, const std::string& op, int b) {
    if (op == ">") return a > b;
    if (op == "<") return a < b;
    if (op == ">=") return a >= b;
    if (op == "<=") return a <= b;
    if (op == "==") return a == b;
    if (op == "!=") return a != b;
    return false;
}

// Compare strings
bool compareString(const std::string& a, const std::string& op, const std::string& b) {
    if (op == "==") return a == b;
    if (op == "!=") return a != b;
    return false;
}

// Evaluate a comparison condition
bool evaluateComparison(const std::string& condition, const std::string& op, const Charge& charge) {
    size_t opPos = condition.find(op);
    if (opPos == std::string::npos) {
        return false;
    }
    
    std::string left = condition.substr(0, opPos);
    std::string right = condition.substr(opPos + op.length());
    
    // Trim whitespace
    left.erase(0, left.find_first_not_of(" \t"));
    left.erase(left.find_last_not_of(" \t") + 1);
    right.erase(0, right.find_first_not_of(" \t"));
    right.erase(right.find_last_not_of(" \t") + 1);
    
    // Handle amount comparisons
    if (left == "amount") {
        int rightVal = std::stoi(right);
        return compareInt(charge.amount, op, rightVal);
    }
    
    // Handle country comparisons
    if (left == "card_country" && right == "ip_country") {
        return compareString(charge.cardCountry, op, charge.ipCountry);
    }
    
    return false;
}

// EvaluateCondition evaluates a condition against a charge
bool evaluateCondition(const std::string& condition, const Charge& charge) {
    std::string cond = condition;
    
    // Trim whitespace
    cond.erase(0, cond.find_first_not_of(" \t"));
    cond.erase(cond.find_last_not_of(" \t") + 1);
    
    // Handle AND conditions
    if (cond.find(" AND ") != std::string::npos) {
        std::istringstream iss(cond);
        std::string part;
        std::vector<std::string> parts;
        
        while (std::getline(iss, part, '&')) {
            if (part.find(" AND ") != std::string::npos) {
                size_t andPos = part.find(" AND ");
                parts.push_back(part.substr(0, andPos));
                parts.push_back(part.substr(andPos + 5));
            } else {
                parts.push_back(part);
            }
        }
        
        // Simple split by " AND "
        parts.clear();
        size_t start = 0;
        while (true) {
            size_t andPos = cond.find(" AND ", start);
            if (andPos == std::string::npos) {
                parts.push_back(cond.substr(start));
                break;
            }
            parts.push_back(cond.substr(start, andPos - start));
            start = andPos + 5;
        }
        
        for (const auto& p : parts) {
            std::string trimmed = p;
            trimmed.erase(0, trimmed.find_first_not_of(" \t"));
            trimmed.erase(trimmed.find_last_not_of(" \t") + 1);
            if (!evaluateCondition(trimmed, charge)) {
                return false;
            }
        }
        return true;
    }
    
    // Handle comparison operators
    std::vector<std::string> operators = {"!=", "==", ">=", "<=", ">", "<"};
    for (const auto& op : operators) {
        if (cond.find(op) != std::string::npos) {
            return evaluateComparison(cond, op, charge);
        }
    }
    
    return false;
}

// ProcessRadarRules processes charge and rules
int processRadarRules(const std::string& input) {
    std::vector<std::string> parts;
    std::istringstream iss(input);
    std::string token;
    
    while (std::getline(iss, token, ',')) {
        parts.push_back(token);
    }
    
    if (parts.empty()) {
        return 1; // Default allow
    }
    
    // Parse charge (first part)
    Charge charge = parseCharge(parts[0]);
    
    // Process rules in order
    for (size_t i = 1; i < parts.size(); i++) {
        std::string ruleStr = parts[i];
        
        // Trim whitespace
        ruleStr.erase(0, ruleStr.find_first_not_of(" \t"));
        ruleStr.erase(ruleStr.find_last_not_of(" \t") + 1);
        
        if (ruleStr.empty()) {
            continue;
        }
        
        Rule rule = parseRule(ruleStr);
        if (rule.type.empty()) {
            continue;
        }
        
        bool conditionMet = evaluateCondition(rule.condition, charge);
        
        if (rule.type == "ALLOW") {
            if (conditionMet) {
                return 1; // Allowed
            }
        } else if (rule.type == "BLOCK") {
            if (conditionMet) {
                return 0; // Blocked
            }
        }
    }
    
    // Default: if no rule matched, allow
    return 1;
}

int main() {
    std::cout << std::string(60, '=') << std::endl;
    std::cout << "Stripe Radar Rules Engine" << std::endl;
    std::cout << std::string(60, '=') << std::endl;
    
    // Example 1: Basic block
    std::cout << "\nExample 1: Basic Block" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    std::string input1 = "CHARGE: card_country=US&currency=USD&amount=150&ip_country=CA,BLOCK:amount > 100";
    int result1 = processRadarRules(input1);
    std::cout << "Input: " << input1 << std::endl;
    std::cout << "Result: " << result1 << " (0 = blocked, 1 = allowed)" << std::endl;
    std::cout << "Expected: 0 (blocked because amount 150 > 100)" << std::endl;
    
    // Example 2: Compound block
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 2: Compound Block" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    std::string input2 = "CHARGE: card_country=US&currency=USD&amount=150&ip_country=CA,ALLOW:amount<100,BLOCK:card_country != ip_country AND amount > 100";
    int result2 = processRadarRules(input2);
    std::cout << "Input: " << input2 << std::endl;
    std::cout << "Result: " << result2 << " (0 = blocked, 1 = allowed)" << std::endl;
    std::cout << "Expected: 0 (blocked: card_country(US) != ip_country(CA) AND amount(150) > 100)" << std::endl;
    
    // Example 3: ALLOW rule matches
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 3: ALLOW Rule Matches" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    std::string input3 = "CHARGE: card_country=US&currency=USD&amount=50&ip_country=CA,ALLOW:amount<100,BLOCK:card_country != ip_country AND amount > 100";
    int result3 = processRadarRules(input3);
    std::cout << "Input: " << input3 << std::endl;
    std::cout << "Result: " << result3 << " (0 = blocked, 1 = allowed)" << std::endl;
    std::cout << "Expected: 1 (allowed because amount 50 < 100)" << std::endl;
    
    // Example 4: Country mismatch only
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 4: Country Mismatch Only" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    std::string input4 = "CHARGE: card_country=US&currency=USD&amount=50&ip_country=CA,BLOCK:card_country != ip_country";
    int result4 = processRadarRules(input4);
    std::cout << "Input: " << input4 << std::endl;
    std::cout << "Result: " << result4 << " (0 = blocked, 1 = allowed)" << std::endl;
    std::cout << "Expected: 0 (blocked: US != CA)" << std::endl;
    
    // Example 5: No rules match
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 5: No Rules Match" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    std::string input5 = "CHARGE: card_country=US&currency=USD&amount=50&ip_country=US,ALLOW:amount>100,BLOCK:card_country != ip_country";
    int result5 = processRadarRules(input5);
    std::cout << "Input: " << input5 << std::endl;
    std::cout << "Result: " << result5 << " (0 = blocked, 1 = allowed)" << std::endl;
    std::cout << "Expected: 1 (default allow, no rules matched)" << std::endl;
    
    return 0;
}

