package main

import (
	"fmt"
	"log"
	"net/http"
	"reddit-comments/pkg/handler"
	"reddit-comments/pkg/repository"
	"reddit-comments/pkg/service"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize components
	repo := repository.NewRepository()
	service := service.NewService(repo)
	handler := handler.NewHandler(service)

	// Setup routes
	router := mux.NewRouter()

	// User routes
	router.HandleFunc("/api/users", handler.CreateUser).Methods("POST")
	router.HandleFunc("/api/users", handler.GetAllUsers).Methods("GET")
	router.HandleFunc("/api/users/{id}", handler.GetUser).Methods("GET")

	// Post routes
	router.HandleFunc("/api/posts", handler.CreatePost).Methods("POST")
	router.HandleFunc("/api/posts", handler.GetAllPosts).Methods("GET")
	router.HandleFunc("/api/posts/{id}", handler.GetPost).Methods("GET")
	router.HandleFunc("/api/posts/{id}/comments", handler.GetCommentsByPost).Methods("GET")
	router.HandleFunc("/api/posts/{id}/comments/tree", handler.GetCommentTree).Methods("GET")
	router.HandleFunc("/api/posts/{id}/comments/stats", handler.GetCommentStats).Methods("GET")

	// Comment routes
	router.HandleFunc("/api/comments", handler.CreateComment).Methods("POST")
	router.HandleFunc("/api/comments", handler.GetAllComments).Methods("GET")
	router.HandleFunc("/api/comments/{id}", handler.GetComment).Methods("GET")
	router.HandleFunc("/api/comments/{id}/path", handler.GetCommentPath).Methods("GET")
	router.HandleFunc("/api/comments/{id}/ancestors", handler.GetCommentAncestors).Methods("GET")

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
	fmt.Println("Reddit Comments API Server starting on :8080")
	fmt.Println("Available endpoints:")
	fmt.Println("  POST   /api/users")
	fmt.Println("  GET    /api/users")
	fmt.Println("  GET    /api/users/{id}")
	fmt.Println("  POST   /api/posts")
	fmt.Println("  GET    /api/posts")
	fmt.Println("  GET    /api/posts/{id}")
	fmt.Println("  GET    /api/posts/{id}/comments")
	fmt.Println("  GET    /api/posts/{id}/comments/tree")
	fmt.Println("  GET    /api/posts/{id}/comments/stats")
	fmt.Println("  POST   /api/comments")
	fmt.Println("  GET    /api/comments")
	fmt.Println("  GET    /api/comments/{id}")
	fmt.Println("  GET    /api/comments/{id}/path")
	fmt.Println("  GET    /api/comments/{id}/ancestors")
	fmt.Println()

	// Seed some sample data
	seedSampleData(service)

	log.Fatal(http.ListenAndServe(":8080", router))
}

// seedSampleData creates some sample data for testing
func seedSampleData(service *service.Service) {
	fmt.Println("Seeding sample data...")

	// Create users
	user1, _ := service.CreateUser("alice", "alice@example.com")
	user2, _ := service.CreateUser("bob", "bob@example.com")
	user3, _ := service.CreateUser("charlie", "charlie@example.com")

	// Create posts
	post1, _ := service.CreatePost("Go vs Rust", "Which language is better for systems programming?", user1.ID)
	post2, _ := service.CreatePost("Microservices Architecture", "Best practices for designing microservices", user2.ID)

	// Create comments for post1
	comment1, _ := service.CreateComment(post1.ID, user2.ID, "Go is definitely better for concurrency!", nil)
	comment2, _ := service.CreateComment(post1.ID, user3.ID, "Rust has better memory safety", nil)

	// Create replies to comment1
	reply1, _ := service.CreateComment(post1.ID, user1.ID, "I agree! Goroutines are amazing", &comment1.ID)
	reply2, _ := service.CreateComment(post1.ID, user3.ID, "But Rust's ownership model is revolutionary", &comment1.ID)

	// Create nested reply
	nestedReply, _ := service.CreateComment(post1.ID, user2.ID, "Both have their merits, but Go's simplicity wins", &reply1.ID)

	// Create comments for post2
	comment3, _ := service.CreateComment(post2.ID, user1.ID, "Great post! I've been using Docker for containerization", nil)
	comment4, _ := service.CreateComment(post2.ID, user3.ID, "Kubernetes is essential for orchestration", nil)

	// Create reply to comment3
	reply3, _ := service.CreateComment(post2.ID, user2.ID, "Docker Compose is perfect for local development", &comment3.ID)

	fmt.Printf("Created %d users, %d posts, and %d comments\n", 3, 2, 7)
	fmt.Println("Sample data seeded successfully!")
	fmt.Println()
}
