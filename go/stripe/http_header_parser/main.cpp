#include <iostream>
#include <string>
#include <map>
#include <sstream>
#include <algorithm>
#include <cctype>

class HTTPHeaderParser {
private:
    std::map<std::string, std::string> headers;
    
    // Convert string to lowercase
    std::string toLower(const std::string& str) {
        std::string result = str;
        std::transform(result.begin(), result.end(), result.begin(), ::tolower);
        return result;
    }
    
    // Trim whitespace from string
    std::string trim(const std::string& str) {
        size_t first = str.find_first_not_of(" \t\n\r");
        if (first == std::string::npos) return "";
        size_t last = str.find_last_not_of(" \t\n\r");
        return str.substr(first, (last - first + 1));
    }
    
public:
    // Parse headers from raw string
    void parseHeaders(const std::string& rawHeaders) {
        std::istringstream stream(rawHeaders);
        std::string line;
        
        while (std::getline(stream, line)) {
            line = trim(line);
            if (line.empty()) continue;
            
            size_t colonPos = line.find(':');
            if (colonPos != std::string::npos) {
                std::string key = trim(line.substr(0, colonPos));
                std::string value = trim(line.substr(colonPos + 1));
                headers[toLower(key)] = value;
            }
        }
    }
    
    // Get header value by key (case-insensitive)
    bool getHeader(const std::string& key, std::string& value) {
        auto it = headers.find(toLower(key));
        if (it != headers.end()) {
            value = it->second;
            return true;
        }
        return false;
    }
    
    // Get Content-Type header
    bool getContentType(std::string& contentType) {
        std::string value;
        if (!getHeader("Content-Type", value)) {
            return false;
        }
        
        // Extract MIME type (before semicolon if present)
        size_t semicolonPos = value.find(';');
        if (semicolonPos != std::string::npos) {
            contentType = trim(value.substr(0, semicolonPos));
        } else {
            contentType = trim(value);
        }
        return true;
    }
    
    // Get Content-Length as integer
    bool getContentLength(int& length) {
        std::string value;
        if (!getHeader("Content-Length", value)) {
            return false;
        }
        
        try {
            length = std::stoi(value);
            return true;
        } catch (...) {
            return false;
        }
    }
    
    // Check if Transfer-Encoding is chunked
    bool isChunked() {
        std::string encoding;
        if (!getHeader("Transfer-Encoding", encoding)) {
            return false;
        }
        return toLower(encoding) == "chunked";
    }
    
    // Get User-Agent header
    bool getUserAgent(std::string& userAgent) {
        return getHeader("User-Agent", userAgent);
    }
    
    // Get all headers
    std::map<std::string, std::string> getAllHeaders() {
        return headers;
    }
};

int main() {
    std::cout << std::string(60, '=') << std::endl;
    std::cout << "HTTP Header Parser" << std::endl;
    std::cout << std::string(60, '=') << std::endl;
    
    // Example 1: Basic headers
    std::string rawHeaders1 = "Content-Type: application/json\n"
                              "Content-Length: 1024\n"
                              "User-Agent: Mozilla/5.0\n"
                              "Accept: application/json";
    
    std::cout << "\nExample 1: Basic Headers" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    std::cout << "Raw Headers:" << std::endl;
    std::cout << rawHeaders1 << std::endl << std::endl;
    
    HTTPHeaderParser parser1;
    parser1.parseHeaders(rawHeaders1);
    
    std::string contentType;
    if (parser1.getContentType(contentType)) {
        std::cout << "Content-Type: " << contentType << std::endl;
    }
    
    int contentLength;
    if (parser1.getContentLength(contentLength)) {
        std::cout << "Content-Length: " << contentLength << std::endl;
    }
    
    std::string userAgent;
    if (parser1.getUserAgent(userAgent)) {
        std::cout << "User-Agent: " << userAgent << std::endl;
    }
    
    std::cout << "Is Chunked: " << (parser1.isChunked() ? "true" : "false") << std::endl;
    
    // Example 2: Headers with Transfer-Encoding
    std::string rawHeaders2 = "Content-Type: text/html; charset=utf-8\n"
                              "Transfer-Encoding: chunked\n"
                              "Server: nginx/1.18.0";
    
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 2: Headers with Transfer-Encoding" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    std::cout << "Raw Headers:" << std::endl;
    std::cout << rawHeaders2 << std::endl << std::endl;
    
    HTTPHeaderParser parser2;
    parser2.parseHeaders(rawHeaders2);
    
    if (parser2.getContentType(contentType)) {
        std::cout << "Content-Type: " << contentType << std::endl;
    }
    
    std::cout << "Is Chunked: " << (parser2.isChunked() ? "true" : "false") << std::endl;
    
    std::string server;
    if (parser2.getHeader("Server", server)) {
        std::cout << "Server: " << server << std::endl;
    }
    
    // Example 3: Case-insensitive lookup
    std::string rawHeaders3 = "content-type: application/xml\n"
                              "CONTENT-LENGTH: 2048\n"
                              "user-agent: curl/7.68.0";
    
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 3: Case-Insensitive Lookup" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    std::cout << "Raw Headers:" << std::endl;
    std::cout << rawHeaders3 << std::endl << std::endl;
    
    HTTPHeaderParser parser3;
    parser3.parseHeaders(rawHeaders3);
    
    if (parser3.getContentType(contentType)) {
        std::cout << "Content-Type: " << contentType << std::endl;
    }
    
    if (parser3.getContentLength(contentLength)) {
        std::cout << "Content-Length: " << contentLength << std::endl;
    }
    
    if (parser3.getUserAgent(userAgent)) {
        std::cout << "User-Agent: " << userAgent << std::endl;
    }
    
    // Example 4: Missing headers
    std::cout << "\n" << std::string(60, '=') << std::endl;
    std::cout << "Example 4: Missing Headers" << std::endl;
    std::cout << std::string(61, '-') << std::endl;
    
    std::string auth;
    if (parser3.getHeader("Authorization", auth)) {
        std::cout << "Authorization header found" << std::endl;
    } else {
        std::cout << "Authorization header not found" << std::endl;
    }
    
    return 0;
}

