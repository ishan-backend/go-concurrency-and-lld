package main

import (
	"fmt"
	"strings"
)

// HTTPHeaderParser parses HTTP headers and extracts specific values
type HTTPHeaderParser struct {
	headers map[string]string
}

// NewHTTPHeaderParser creates a new parser from raw header string
func NewHTTPHeaderParser(rawHeaders string) *HTTPHeaderParser {
	parser := &HTTPHeaderParser{
		headers: make(map[string]string),
	}

	// Parse headers (format: "key1: value1\nkey2: value2")
	lines := strings.Split(rawHeaders, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(strings.ToLower(parts[0]))
			value := strings.TrimSpace(parts[1])
			parser.headers[key] = value
		}
	}

	return parser
}

// GetHeader retrieves a header value by key (case-insensitive)
func (p *HTTPHeaderParser) GetHeader(key string) (string, bool) {
	key = strings.ToLower(key)
	value, exists := p.headers[key]
	return value, exists
}

// GetContentType extracts the content type from Content-Type header
func (p *HTTPHeaderParser) GetContentType() (string, bool) {
	contentType, exists := p.GetHeader("Content-Type")
	if !exists {
		return "", false
	}

	// Extract MIME type (before semicolon if present)
	parts := strings.Split(contentType, ";")
	return strings.TrimSpace(parts[0]), true
}

// GetContentLength extracts the content length as integer
func (p *HTTPHeaderParser) GetContentLength() (int, bool) {
	contentLength, exists := p.GetHeader("Content-Length")
	if !exists {
		return 0, false
	}

	var length int
	_, err := fmt.Sscanf(contentLength, "%d", &length)
	if err != nil {
		return 0, false
	}

	return length, true
}

// IsChunked checks if Transfer-Encoding is chunked
func (p *HTTPHeaderParser) IsChunked() bool {
	encoding, exists := p.GetHeader("Transfer-Encoding")
	if !exists {
		return false
	}
	return strings.ToLower(encoding) == "chunked"
}

// GetUserAgent extracts the User-Agent header
func (p *HTTPHeaderParser) GetUserAgent() (string, bool) {
	return p.GetHeader("User-Agent")
}

// GetAllHeaders returns all headers as a map
func (p *HTTPHeaderParser) GetAllHeaders() map[string]string {
	result := make(map[string]string)
	for k, v := range p.headers {
		result[k] = v
	}
	return result
}

func main() {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("HTTP Header Parser")
	fmt.Println(strings.Repeat("=", 60))

	// Example 1: Basic headers
	rawHeaders1 := `Content-Type: application/json
Content-Length: 1024
User-Agent: Mozilla/5.0
Accept: application/json`

	fmt.Println("\nExample 1: Basic Headers")
	fmt.Println("-" + strings.Repeat("-", 60))
	fmt.Println("Raw Headers:")
	fmt.Println(rawHeaders1)
	fmt.Println()

	parser1 := NewHTTPHeaderParser(rawHeaders1)

	if contentType, ok := parser1.GetContentType(); ok {
		fmt.Printf("Content-Type: %s\n", contentType)
	}

	if contentLength, ok := parser1.GetContentLength(); ok {
		fmt.Printf("Content-Length: %d\n", contentLength)
	}

	if userAgent, ok := parser1.GetUserAgent(); ok {
		fmt.Printf("User-Agent: %s\n", userAgent)
	}

	fmt.Printf("Is Chunked: %v\n", parser1.IsChunked())

	// Example 2: Headers with Transfer-Encoding
	rawHeaders2 := `Content-Type: text/html; charset=utf-8
Transfer-Encoding: chunked
Server: nginx/1.18.0`

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 2: Headers with Transfer-Encoding")
	fmt.Println("-" + strings.Repeat("-", 60))
	fmt.Println("Raw Headers:")
	fmt.Println(rawHeaders2)
	fmt.Println()

	parser2 := NewHTTPHeaderParser(rawHeaders2)

	if contentType, ok := parser2.GetContentType(); ok {
		fmt.Printf("Content-Type: %s\n", contentType)
	}

	fmt.Printf("Is Chunked: %v\n", parser2.IsChunked())

	if server, ok := parser2.GetHeader("Server"); ok {
		fmt.Printf("Server: %s\n", server)
	}

	// Example 3: Case-insensitive lookup
	rawHeaders3 := `content-type: application/xml
CONTENT-LENGTH: 2048
user-agent: curl/7.68.0`

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 3: Case-Insensitive Lookup")
	fmt.Println("-" + strings.Repeat("-", 60))
	fmt.Println("Raw Headers:")
	fmt.Println(rawHeaders3)
	fmt.Println()

	parser3 := NewHTTPHeaderParser(rawHeaders3)

	if contentType, ok := parser3.GetContentType(); ok {
		fmt.Printf("Content-Type: %s\n", contentType)
	}

	if contentLength, ok := parser3.GetContentLength(); ok {
		fmt.Printf("Content-Length: %d\n", contentLength)
	}

	if userAgent, ok := parser3.GetUserAgent(); ok {
		fmt.Printf("User-Agent: %s\n", userAgent)
	}

	// Example 4: Missing headers
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Example 4: Missing Headers")
	fmt.Println("-" + strings.Repeat("-", 60))

	if _, ok := parser3.GetHeader("Authorization"); ok {
		fmt.Println("Authorization header found")
	} else {
		fmt.Println("Authorization header not found")
	}
}
