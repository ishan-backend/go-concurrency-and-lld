#include <iostream>
#include <string>
#include <vector>
#include <map>
#include <algorithm>
#include <ctime>
#include <chrono>
#include <iomanip>
#include <sstream>

struct Transaction {
    double amount;
    std::chrono::system_clock::time_point date;
    std::string description;
    
    Transaction(double amt, const std::chrono::system_clock::time_point& d, const std::string& desc)
        : amount(amt), date(d), description(desc) {}
};

class StripeCapital {
private:
    std::vector<Transaction> transactions;
    
    // Format time to YYYY-MM string
    std::string formatMonth(const std::chrono::system_clock::time_point& tp) {
        auto time_t = std::chrono::system_clock::to_time_t(tp);
        std::tm* tm = std::localtime(&time_t);
        std::ostringstream oss;
        oss << std::setfill('0') << (1900 + tm->tm_year) << "-"
            << std::setw(2) << (tm->tm_mon + 1);
        return oss.str();
    }
    
public:
    // Add a transaction to the history
    void addTransaction(double amount, const std::chrono::system_clock::time_point& date, 
                       const std::string& description) {
        transactions.push_back(Transaction(amount, date, description));
    }
    
    // Calculate total capital from all transactions
    double calculateTotalCapital() {
        double total = 0.0;
        for (const auto& txn : transactions) {
            total += txn.amount;
        }
        return total;
    }
    
    // Calculate average monthly revenue
    double calculateAverageMonthlyRevenue() {
        if (transactions.empty()) {
            return 0.0;
        }
        
        // Group transactions by month
        std::map<std::string, double> monthlyRevenue;
        for (const auto& txn : transactions) {
            if (txn.amount > 0) { // Only count positive amounts as revenue
                std::string monthKey = formatMonth(txn.date);
                monthlyRevenue[monthKey] += txn.amount;
            }
        }
        
        if (monthlyRevenue.empty()) {
            return 0.0;
        }
        
        double total = 0.0;
        for (const auto& pair : monthlyRevenue) {
            total += pair.second;
        }
        
        return total / monthlyRevenue.size();
    }
    
    // Check if business is eligible for a loan
    bool isEligibleForLoan(double threshold) {
        double avgRevenue = calculateAverageMonthlyRevenue();
        return avgRevenue >= threshold;
    }
    
    // Calculate potential loan amount based on revenue
    double getLoanAmount(double multiplier) {
        double avgRevenue = calculateAverageMonthlyRevenue();
        return avgRevenue * multiplier;
    }
    
    // Get monthly revenue breakdown
    std::map<std::string, double> getMonthlyRevenueBreakdown() {
        std::map<std::string, double> monthlyRevenue;
        for (const auto& txn : transactions) {
            if (txn.amount > 0) {
                std::string monthKey = formatMonth(txn.date);
                monthlyRevenue[monthKey] += txn.amount;
            }
        }
        return monthlyRevenue;
    }
    
    // Get top N months by revenue
    std::vector<std::string> getTopMonths(int n) {
        auto breakdown = getMonthlyRevenueBreakdown();
        
        struct MonthRevenue {
            std::string month;
            double revenue;
        };
        
        std::vector<MonthRevenue> months;
        for (const auto& pair : breakdown) {
            months.push_back({pair.first, pair.second});
        }
        
        // Sort by revenue (descending)
        std::sort(months.begin(), months.end(),
                  [](const MonthRevenue& a, const MonthRevenue& b) {
                      if (a.revenue != b.revenue) {
                          return a.revenue > b.revenue;
                      }
                      return a.month > b.month;
                  });
        
        std::vector<std::string> result;
        for (size_t i = 0; i < months.size() && i < static_cast<size_t>(n); i++) {
            result.push_back(months[i].month);
        }
        
        return result;
    }
};

int main() {
    std::cout << std::string(60, '=') << std::endl;
    std::cout << "Stripe Capital - Capital Financing System" << std::endl;
    std::cout << std::string(60, '=') << std::endl;
    
    StripeCapital capital;
    
    // Example 1: Basic eligibility check
    std::cout << "\nExample 1: Basic Eligibility Check" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    auto baseDate = std::chrono::system_clock::now();
    capital.addTransaction(12000, baseDate, "Payment 1");
    capital.addTransaction(15000, baseDate + std::chrono::hours(24 * 30), "Payment 2");
    capital.addTransaction(11000, baseDate + std::chrono::hours(24 * 60), "Payment 3");
    capital.addTransaction(13000, baseDate + std::chrono::hours(24 * 90), "Payment 4");
    
    double totalCapital = capital.calculateTotalCapital();
    double avgRevenue = capital.calculateAverageMonthlyRevenue();
    bool eligible = capital.isEligibleForLoan(10000);
    
    std::cout << std::fixed << std::setprecision(2);
    std::cout << "Total Capital: $" << totalCapital << std::endl;
    std::cout << "Average Monthly Revenue: $" << avgRevenue << std::endl;
    std::cout << "Eligible for Loan (threshold $10,000): " << (eligible ? "true" : "false") << std::endl;
    
    // Example 2: Loan amount calculation
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 2: Loan Amount Calculation" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    double loanAmount = capital.getLoanAmount(3.0); // 3x monthly revenue
    std::cout << "Potential Loan Amount (3x monthly revenue): $" << loanAmount << std::endl;
    
    // Example 3: Monthly breakdown
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 3: Monthly Revenue Breakdown" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    auto breakdown = capital.getMonthlyRevenueBreakdown();
    std::cout << "Monthly Revenue Breakdown:" << std::endl;
    for (const auto& pair : breakdown) {
        std::cout << "  " << pair.first << ": $" << pair.second << std::endl;
    }
    
    // Example 4: Top months
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 4: Top Months by Revenue" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    auto topMonths = capital.getTopMonths(3);
    std::cout << "Top 3 months:" << std::endl;
    for (size_t i = 0; i < topMonths.size(); i++) {
        std::cout << "  " << (i + 1) << ". " << topMonths[i] 
                  << ": $" << breakdown[topMonths[i]] << std::endl;
    }
    
    // Example 5: Not eligible case
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 5: Low Revenue Business" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    StripeCapital capital2;
    capital2.addTransaction(5000, baseDate, "Payment 1");
    capital2.addTransaction(6000, baseDate + std::chrono::hours(24 * 30), "Payment 2");
    capital2.addTransaction(4000, baseDate + std::chrono::hours(24 * 60), "Payment 3");
    
    double avgRevenue2 = capital2.calculateAverageMonthlyRevenue();
    bool eligible2 = capital2.isEligibleForLoan(10000);
    
    std::cout << "Average Monthly Revenue: $" << avgRevenue2 << std::endl;
    std::cout << "Eligible for Loan (threshold $10,000): " << (eligible2 ? "true" : "false") << std::endl;
    
    return 0;
}

