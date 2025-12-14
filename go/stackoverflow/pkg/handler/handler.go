package handler

import (
	"encoding/json"
	"net/http"
	"stackoverflow/pkg/service"
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

// CreatePost handles POST /api/posts - Create a new post
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

// AddComment handles POST /api/comments - Add a comment or reply
func (h *Handler) AddComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		PostID   string  `json:"post_id"`
		AuthorID string  `json:"author_id"`
		Content  string  `json:"content"`
		ParentID *string `json:"parent_id,omitempty"` // Optional: if provided, it's a reply
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	comment, err := h.service.AddComment(req.PostID, req.AuthorID, req.Content, req.ParentID)
	if err != nil {
		h.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.SuccessResponse(w, comment)
}

// UpvoteComment handles POST /api/comments/{id}/upvote - Upvote a comment
func (h *Handler) UpvoteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract comment ID from URL path
	// The route will be /api/comments/{id}/upvote, so we need to extract {id}
	// This will be handled by mux.Vars in main.go
	commentID := r.Context().Value("comment_id").(string)
	if commentID == "" {
		h.ErrorResponse(w, http.StatusBadRequest, "Comment ID is required")
		return
	}

	comment, err := h.service.UpvoteComment(commentID)
	if err != nil {
		h.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	h.SuccessResponse(w, comment)
}

// GetAllCommentsRecursively handles GET /api/posts/{id}/comments - Get all comments and replies recursively
func (h *Handler) GetAllCommentsRecursively(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract post ID from URL path
	// This will be handled by mux.Vars in main.go
	postID := r.Context().Value("post_id").(string)
	if postID == "" {
		h.ErrorResponse(w, http.StatusBadRequest, "Post ID is required")
		return
	}

	comments, err := h.service.GetAllCommentsRecursively(postID)
	if err != nil {
		h.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	h.SuccessResponse(w, comments)
}
