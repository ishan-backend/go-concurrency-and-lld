package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

/*
🔑 Hashing Logic & Collision Handling
Strategy Used: Base62Encoder (MD5 + Base62)

Takes input: "https://example.com"

Creates MD5 hash → gives 32 hex characters

Takes first 8 characters of that hex → converts into a number

Then converts that number into Base62 (digits + lowercase + uppercase = 62 chars)

Ensures result is always 7 characters long

✅ This gives consistent output for same input — deterministic hashing.

But Hashes Can Collide — How is Collision Handled?

Inside ShortenURL():

for attempt := 0; attempt < s.maxRetries; attempt++ {
    input := originalURL
    if attempt > 0 {
        input = fmt.Sprintf("%s-%d-%d", originalURL, time.Now().UnixNano(), rand.Intn(10000))
    }
    shortCode = s.encoder.Encode(input)
    err = s.repo.Save(urlObj)
    if err == nil { return urlObj, nil }
}


If collision happens (duplicate short code), retry with randomness appended → new hash!

*/

// ============= Domain Models =============

type URL struct {
	ID          string    `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortCode   string    `json:"short_code"`
	CreatedAt   time.Time `json:"created_at"`
	Clicks      int64     `json:"clicks"`
}

// ============= Repository Pattern =============

// URLRepository defines the interface for URL storage operations
type URLRepository interface {
	Save(url *URL) error
	FindByShortCode(shortCode string) (*URL, error)
	FindByOriginalURL(originalURL string) (*URL, error)
	IncrementClicks(shortCode string) error
}

// InMemoryRepository implements URLRepository using in-memory storage
type InMemoryRepository struct {
	mu              sync.RWMutex
	urlsByShortCode map[string]*URL
	urlsByOriginal  map[string]*URL
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		urlsByShortCode: make(map[string]*URL),
		urlsByOriginal:  make(map[string]*URL),
	}
}

func (r *InMemoryRepository) Save(url *URL) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check for collision
	if _, exists := r.urlsByShortCode[url.ShortCode]; exists {
		return errors.New("short code already exists")
	}

	r.urlsByShortCode[url.ShortCode] = url
	r.urlsByOriginal[url.OriginalURL] = url
	return nil
}

func (r *InMemoryRepository) FindByShortCode(shortCode string) (*URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	url, exists := r.urlsByShortCode[shortCode]
	if !exists {
		return nil, errors.New("URL not found")
	}
	return url, nil
}

func (r *InMemoryRepository) FindByOriginalURL(originalURL string) (*URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	url, exists := r.urlsByOriginal[originalURL]
	if !exists {
		return nil, errors.New("URL not found")
	}
	return url, nil
}

func (r *InMemoryRepository) IncrementClicks(shortCode string) error {
	r.mu.RLock()
	url, exists := r.urlsByShortCode[shortCode]
	r.mu.RUnlock()
	if !exists {
		return errors.New("URL not found")
	}
	atomic.AddInt64(&url.Clicks, 1)
	return nil
}

// ============= Strategy Pattern for Encoding =============

// EncodingStrategy defines interface for URL encoding algorithms
type EncodingStrategy interface {
	Encode(input string) string
}

// Base62Encoder uses base62 encoding (MD5 hash truncated)
type Base62Encoder struct{}

func (e *Base62Encoder) Encode(input string) string {
	const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Create MD5 hash
	hash := md5.Sum([]byte(input))
	hashStr := hex.EncodeToString(hash[:])

	// Convert first 8 chars to number and then to base62
	var num uint64
	for i := 0; i < 8 && i < len(hashStr); i++ {
		num = num*16 + uint64(hashStr[i])
	}

	/*
			func PolyHash(input string) int {
		    p := 31
		    m := int(1e9 + 9)
		    hash := 0
		    power := 1
		    for _, c := range input {
		        hash = (hash + int(c-'a'+1)*power) % m
		        power = (power * p) % m
		    }
		    return hash
		}

	*/

	// Convert to base62
	result := make([]byte, 0, 7)
	for num > 0 && len(result) < 7 {
		result = append(result, base62Chars[num%62])
		num /= 62
	}

	// Pad if needed
	for len(result) < 7 {
		result = append(result, base62Chars[0])
	}

	// Reverse
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

// CounterEncoder uses sequential counter with base62
type CounterEncoder struct {
	mu      sync.Mutex
	counter uint64
}

func NewCounterEncoder() *CounterEncoder {
	return &CounterEncoder{counter: 1000000} // Start from 1M to get 6-7 char codes
}

func (e *CounterEncoder) Encode(input string) string {
	e.mu.Lock()
	defer e.mu.Unlock()

	const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	num := e.counter
	e.counter++

	result := make([]byte, 0, 7)
	for num > 0 {
		result = append(result, base62Chars[num%62])
		num /= 62
	}

	// Reverse
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

// ============= Service Layer =============

type URLShortenerService struct {
	repo       URLRepository
	encoder    EncodingStrategy
	baseURL    string
	maxRetries int
}

func NewURLShortenerService(repo URLRepository, encoder EncodingStrategy, baseURL string) *URLShortenerService {
	return &URLShortenerService{
		repo:       repo,
		encoder:    encoder,
		baseURL:    baseURL,
		maxRetries: 5,
	}
}

// ShortenURL creates a short URL with collision handling
func (s *URLShortenerService) ShortenURL(originalURL string) (*URL, error) {
	// Validate URL
	if !s.isValidURL(originalURL) {
		return nil, errors.New("invalid URL format")
	}

	// Check if URL already shortened
	if existing, err := s.repo.FindByOriginalURL(originalURL); err == nil {
		return existing, nil
	}

	// Try to create short code with collision handling
	var shortCode string
	var err error

	for attempt := 0; attempt < s.maxRetries; attempt++ {
		input := originalURL
		if attempt > 0 {
			// Add randomness on retry to avoid collision
			input = fmt.Sprintf("%s-%d-%d", originalURL, time.Now().UnixNano(), rand.Intn(10000))
		}

		shortCode = s.encoder.Encode(input)

		urlObj := &URL{
			ID:          shortCode,
			OriginalURL: originalURL,
			ShortCode:   shortCode,
			CreatedAt:   time.Now(),
			Clicks:      0,
		}

		err = s.repo.Save(urlObj)
		if err == nil {
			return urlObj, nil
		}

		// If collision, retry
		log.Printf("Collision detected for code %s, retry %d/%d", shortCode, attempt+1, s.maxRetries)
	}

	return nil, fmt.Errorf("failed to generate unique short code after %d attempts: %w", s.maxRetries, err)
}

func (s *URLShortenerService) GetOriginalURL(shortCode string) (*URL, error) {
	url, err := s.repo.FindByShortCode(shortCode)
	if err != nil {
		return nil, err
	}

	// Increment click counter synchronously to ensure consistency
	if err := s.repo.IncrementClicks(shortCode); err != nil {
		return nil, err
	}

	return url, nil
}

func (s *URLShortenerService) isValidURL(rawURL string) bool {
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}

// ============= HTTP Handlers =============

type URLShortenerHandler struct {
	service *URLShortenerService
}

func NewURLShortenerHandler(service *URLShortenerService) *URLShortenerHandler {
	return &URLShortenerHandler{service: service}
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	ShortCode   string    `json:"short_code"`
	CreatedAt   time.Time `json:"created_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// HandleShorten creates a short URL
func (h *URLShortenerHandler) HandleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.URL) == "" {
		h.sendError(w, "URL is required", http.StatusBadRequest)
		return
	}

	urlObj, err := h.service.ShortenURL(req.URL)
	if err != nil {
		if strings.Contains(err.Error(), "invalid URL") {
			h.sendError(w, err.Error(), http.StatusBadRequest)
		} else {
			h.sendError(w, "Failed to shorten URL", http.StatusInternalServerError)
		}
		return
	}

	response := ShortenResponse{
		ShortURL:    fmt.Sprintf("%s/%s", h.service.baseURL, urlObj.ShortCode),
		OriginalURL: urlObj.OriginalURL,
		ShortCode:   urlObj.ShortCode,
		CreatedAt:   urlObj.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// HandleRedirect redirects short code to original URL
func (h *URLShortenerHandler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	shortCode := strings.TrimPrefix(r.URL.Path, "/")

	if shortCode == "" {
		h.sendError(w, "Short code is required", http.StatusBadRequest)
		return
	}

	urlObj, err := h.service.GetOriginalURL(shortCode)
	if err != nil {
		h.sendError(w, "URL not found", http.StatusNotFound)
		return
	}

	// 301 for permanent redirect (cached by browsers)
	// 302 for temporary redirect (not cached)
	http.Redirect(w, r, urlObj.OriginalURL, http.StatusMovedPermanently)
}

// HandleStats returns statistics for a short code
func (h *URLShortenerHandler) HandleStats(w http.ResponseWriter, r *http.Request) {
	shortCode := r.URL.Query().Get("code")

	if shortCode == "" {
		h.sendError(w, "Short code is required", http.StatusBadRequest)
		return
	}

	// Fetch URL without incrementing clicks
	urlObj, err := h.service.repo.FindByShortCode(shortCode)
	if err != nil {
		h.sendError(w, "URL not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(urlObj)
}

func (h *URLShortenerHandler) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// ============= Middleware =============

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next(w, r)
		log.Printf("Completed in %v", time.Since(start))
	}
}

// ============= Main Application =============

func main() {
	// Initialize dependencies
	repo := NewInMemoryRepository()

	// Choose encoding strategy: Base62 (hash-based) or Counter (sequential)
	encoder := &Base62Encoder{}
	// encoder := NewCounterEncoder() // Alternative strategy

	baseURL := "http://localhost:8080"
	service := NewURLShortenerService(repo, encoder, baseURL)
	handler := NewURLShortenerHandler(service)

	// Setup routes
	http.HandleFunc("/api/shorten", loggingMiddleware(handler.HandleShorten))
	http.HandleFunc("/api/stats", loggingMiddleware(handler.HandleStats))
	http.HandleFunc("/", loggingMiddleware(handler.HandleRedirect))

	log.Println("URL Shortener starting on :8080")
	log.Println("POST /api/shorten - Create short URL")
	log.Println("GET /{shortCode} - Redirect to original URL")
	log.Println("GET /api/stats?code={shortCode} - Get URL statistics")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
