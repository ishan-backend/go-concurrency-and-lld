# HTTP Header Parser

## Problem Statement

Implement an HTTP header parser that can extract and parse various HTTP headers from a raw header string. The parser should:

1. Parse headers from a raw string format (key-value pairs separated by newlines)
2. Extract specific headers like Content-Type, Content-Length, User-Agent, etc.
3. Handle case-insensitive header key lookups
4. Parse Content-Type and extract the MIME type (ignoring parameters like charset)
5. Parse Content-Length as an integer
6. Check if Transfer-Encoding is chunked

### Input Format
Headers are provided as a multi-line string where each line contains a header in the format:
```
key: value
```

Example:
```
Content-Type: application/json
Content-Length: 1024
User-Agent: Mozilla/5.0
```

## Test Cases

### Test Case 1: Basic Headers
- **Input**:
  ```
  Content-Type: application/json
  Content-Length: 1024
  User-Agent: Mozilla/5.0
  Accept: application/json
  ```
- **Expected Output**:
  - Content-Type: `application/json`
  - Content-Length: `1024`
  - User-Agent: `Mozilla/5.0`
  - Is Chunked: `false`

### Test Case 2: Headers with Parameters
- **Input**:
  ```
  Content-Type: text/html; charset=utf-8
  Transfer-Encoding: chunked
  Server: nginx/1.18.0
  ```
- **Expected Output**:
  - Content-Type: `text/html` (MIME type only, without charset)
  - Is Chunked: `true`
  - Server: `nginx/1.18.0`

### Test Case 3: Case-Insensitive Lookup
- **Input**:
  ```
  content-type: application/xml
  CONTENT-LENGTH: 2048
  user-agent: curl/7.68.0
  ```
- **Expected Output**:
  - Content-Type: `application/xml`
  - Content-Length: `2048`
  - User-Agent: `curl/7.68.0`

### Test Case 4: Missing Headers
- **Input**: Headers without Authorization header
- **Expected Output**: Should return false/not found when querying for Authorization

## Running the Solutions

### Go

To run the Go solution:

```bash
cd go/stripe/http_header_parser
go run main.go
```

### C++

To compile and run the C++ solution:

```bash
cd go/stripe/http_header_parser
g++ -o main main.cpp
./main
```

