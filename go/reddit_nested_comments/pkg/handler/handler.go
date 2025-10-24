package handler

import (
	"encoding/json"
	"net/http"
	"reddit-comments/pkg/models"
	"reddit-comments/pkg/service"
	"strconv"
	"strings"
)

// Handler handles HTTP requests
type Handler struct {
	service *service.Service
}

// NewHandler creates a new handler instance
func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ErrorResponse sends an error response
func (h *Handler) ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Error:   message,
	})
}

// SuccessResponse sends a success response
func (h *Handler) SuccessResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    data,
	})
}

// CreateUser handles user creation
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	
	user, err := h.service.CreateUser(req.Username, req.Email)
	if err != nil {
		h.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	
	h.SuccessResponse(w, user)
}

// GetUser handles user retrieval
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	// Extract user ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	if path == "" {
		h.ErrorResponse(w, http.StatusBadRequest, "User ID is required")
		return
	}
	
	user, err := h.service.GetUser(path)
	if err != nil {
		h.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	
	h.SuccessResponse(w, user)
}

// GetAllUsers handles getting all users
func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	users := h.service.GetAllUsers()
	h.SuccessResponse(w, users)
}

// CreatePost handles post creation
func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	var req struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		AuthorID string `json:"author_id"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	
	post, err := h.service.CreatePost(req.Title, req.Content, req.AuthorID)
	if err != nil {
		h.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	
	h.SuccessResponse(w, post)
}

// GetPost handles post retrieval
func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	// Extract post ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	if path == "" {
		h.ErrorResponse(w, http.StatusBadRequest, "Post ID is required")
		return
	}
	
	post, err := h.service.GetPost(path)
	if err != nil {
		h.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	
	h.SuccessResponse(w, post)
}

// GetAllPosts handles getting all posts
func (h *Handler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	posts := h.service.GetAllPosts()
	h.SuccessResponse(w, posts)
}

// CreateComment handles comment creation
func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	var req struct {
		PostID   string  `json:"post_id"`
		AuthorID string  `json:"author_id"`
		Content  string  `json:"content"`
		ParentID *string `json:"parent_id,omitempty"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	
	comment, err := h.service.CreateComment(req.PostID, req.AuthorID, req.Content, req.ParentID)
	if err != nil {
		h.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	
	h.SuccessResponse(w, comment)
}

// GetComment handles comment retrieval
func (h *Handler) GetComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	// Extract comment ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/comments/")
	if path == "" {
		h.ErrorResponse(w, http.StatusBadRequest, "Comment ID is required")
		return
	}
	
	comment, err := h.service.GetComment(path)
	if err != nil {
		h.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	
	h.SuccessResponse(w, comment)
}

// GetCommentsByPost handles getting comments for a post
func (h *Handler) GetCommentsByPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	// Extract post ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	path = strings.TrimSuffix(path, "/comments")
	if path == "" {
		h.ErrorResponse(w, http.StatusBadRequest, "Post ID is required")
		return
	}
	
	comments := h.service.GetCommentsByPost(path)
	h.SuccessResponse(w, comments)
}

// GetCommentTree handles getting the comment tree for a post
func (h *Handler) GetCommentTree(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	// Extract post ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	path = strings.TrimSuffix(path, "/comments/tree")
	if path == "" {
		h.ErrorResponse(w, http.StatusBadRequest, "Post ID is required")
		return
	}
	
	tree, err := h.service.GetCommentTree(path)
	if err != nil {
		h.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	
	h.SuccessResponse(w, tree)
}

// GetCommentPath handles getting the path to a comment
func (h *Handler) GetCommentPath(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	// Extract comment ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/comments/")
	path = strings.TrimSuffix(path, "/path")
	if path == "" {
		h.ErrorResponse(w, http.StatusBadRequest, "Comment ID is required")
		return
	}
	
	commentPath, err := h.service.GetCommentPath(path)
	if err != nil {
		h.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	
	h.SuccessResponse(w, commentPath)
}

// GetCommentAncestors handles getting parent comments
func (h *Handler) GetCommentAncestors(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	// Extract comment ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/comments/")
	path = strings.TrimSuffix(path, "/ancestors")
	if path == "" {
		h.ErrorResponse(w, http.StatusBadRequest, "Comment ID is required")
		return
	}
	
	ancestors, err := h.service.GetCommentAncestors(path)
	if err != nil {
		h.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	
	h.SuccessResponse(w, ancestors)
}

// GetCommentStats handles getting comment statistics
func (h *Handler) GetCommentStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	// Extract post ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	path = strings.TrimSuffix(path, "/comments/stats")
	if path == "" {
		h.ErrorResponse(w, http.StatusBadRequest, "Post ID is required")
		return
	}
	
	stats, err := h.service.GetCommentStats(path)
	if err != nil {
		h.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	
	h.SuccessResponse(w, stats)
}

// GetAllComments handles getting all comments
func (h *Handler) GetAllComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	
	comments := h.service.GetAllComments()
	h.SuccessResponse(w, comments)
}
