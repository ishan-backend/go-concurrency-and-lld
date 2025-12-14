package service

import (
	"fmt"

	"github.com/ishanpandey/LowLevelDesign/go/facebook/models"
	"github.com/ishanpandey/LowLevelDesign/go/facebook/repository"
)

// postServiceImpl implements PostService interface
// Following Single Responsibility Principle (SRP) - only handles post business logic
type postServiceImpl struct {
	userRepo repository.UserRepository // Depends on abstraction (DIP)
}

// NewPostService creates a new instance of PostService
// Following Dependency Inversion Principle (DIP) - accepts repository interface
func NewPostService(userRepo repository.UserRepository) PostService {
	return &postServiceImpl{
		userRepo: userRepo,
	}
}

// AddPost adds a post for a user
// Posts are inserted at the head of the doubly-linked list for O(1) insertion
func (s *postServiceImpl) AddPost(userID int64, postID int) error {
	fmt.Printf("Adding Post = %d from user - %d\n", postID, userID)

	// Get the user
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return fmt.Errorf("user %d not found: %w", userID, err)
	}

	// Check if post already exists
	if _, exists := user.Posts[postID]; exists {
		return fmt.Errorf("post %d already exists for user %d", postID, userID)
	}

	// Create new post
	post := models.NewPost(postID)

	// Insert at head of doubly-linked list (most recent first)
	head := user.Head
	mostRecent := head.Next

	head.Next = post
	post.Next = mostRecent
	mostRecent.Prev = post
	post.Prev = head

	// Add to posts map for O(1) lookup
	user.Posts[postID] = post

	return nil
}

// DeletePost deletes a post for a user
func (s *postServiceImpl) DeletePost(userID int64, postID int) error {
	fmt.Printf("Deleting Post = %d from user - %d\n", postID, userID)

	// Get the user
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return fmt.Errorf("user %d not found: %w", userID, err)
	}

	// Check if post exists
	post, exists := user.Posts[postID]
	if !exists {
		return fmt.Errorf("post %d not found for user %d", postID, userID)
	}

	// Remove from doubly-linked list
	prev := post.Prev
	next := post.Next
	prev.Next = next
	next.Prev = prev

	// Remove from posts map
	delete(user.Posts, postID)

	return nil
}

// GetPost retrieves a post for a user
func (s *postServiceImpl) GetPost(userID int64, postID int) (*models.Post, error) {
	// Get the user
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return nil, fmt.Errorf("user %d not found: %w", userID, err)
	}

	// Get post from map
	post, exists := user.Posts[postID]
	if !exists {
		return nil, fmt.Errorf("post %d not found for user %d", postID, userID)
	}

	return post, nil
}
