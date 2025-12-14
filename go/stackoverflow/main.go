package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"stackoverflow/pkg/handler"
	"stackoverflow/pkg/repository"
	"stackoverflow/pkg/service"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize components
	repo := repository.NewRepository()
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	// Setup routes
	router := mux.NewRouter()

	// API Routes
	router.HandleFunc("/api/posts", h.CreatePost).Methods("POST")
	router.HandleFunc("/api/comments", h.AddComment).Methods("POST")
	
	// For upvote, extract comment ID from URL and pass via context
	router.HandleFunc("/api/comments/{id}/upvote", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		commentID := vars["id"]
		ctx := context.WithValue(r.Context(), "comment_id", commentID)
		h.UpvoteComment(w, r.WithContext(ctx))
	}).Methods("POST")

	// Get all comments recursively for a post
	router.HandleFunc("/api/posts/{id}/comments", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		postID := vars["id"]
		ctx := context.WithValue(r.Context(), "post_id", postID)
		h.GetAllCommentsRecursively(w, r.WithContext(ctx))
	}).Methods("GET")

	// Add CORS middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// Start server
	fmt.Println("Stack Overflow API Server starting on :8080")
	fmt.Println("Available endpoints:")
	fmt.Println("  POST   /api/posts                    - Create a new post")
	fmt.Println("  POST   /api/comments                 - Add a comment or reply")
	fmt.Println("  POST   /api/comments/{id}/upvote     - Upvote a comment")
	fmt.Println("  GET    /api/posts/{id}/comments      - Get all comments and replies recursively")
	fmt.Println()

	log.Fatal(http.ListenAndServe(":8080", router))
}

