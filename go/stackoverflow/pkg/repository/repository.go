package repository

import (
	"errors"
	"stackoverflow/pkg/models"
	"sync"
)

// Repository handles data storage and retrieval
type Repository struct {
	posts    map[string]*models.Post
	comments map[string]*models.Comment
	mu       sync.RWMutex
}

// NewRepository creates a new repository instance
func NewRepository() *Repository {
	return &Repository{
		posts:    make(map[string]*models.Post),
		comments: make(map[string]*models.Comment),
	}
}

// Post operations
func (r *Repository) CreatePost(post *models.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.posts[post.ID]; exists {
		return errors.New("post already exists")
	}

	r.posts[post.ID] = post
	return nil
}

func (r *Repository) GetPost(id string) (*models.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, exists := r.posts[id]
	if !exists {
		return nil, errors.New("post not found")
	}

	return post, nil
}

// Comment operations
func (r *Repository) CreateComment(comment *models.Comment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.comments[comment.ID]; exists {
		return errors.New("comment already exists")
	}

	r.comments[comment.ID] = comment
	return nil
}

func (r *Repository) GetComment(id string) (*models.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	comment, exists := r.comments[id]
	if !exists {
		return nil, errors.New("comment not found")
	}

	return comment, nil
}

func (r *Repository) GetCommentsByPost(postID string) []*models.Comment {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var comments []*models.Comment
	for _, comment := range r.comments {
		if comment.PostID == postID {
			comments = append(comments, comment)
		}
	}

	return comments
}

// UpvoteComment increments the upvote count for a comment
func (r *Repository) UpvoteComment(commentID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	comment, exists := r.comments[commentID]
	if !exists {
		return errors.New("comment not found")
	}

	comment.Upvote()
	return nil
}

// GetAllCommentsRecursively returns all comments for a post, including all nested replies, in a flattened list
func (r *Repository) GetAllCommentsRecursively(postID string) []*models.Comment {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Get all comments for this post
	allComments := make([]*models.Comment, 0)
	for _, comment := range r.comments {
		if comment.PostID == postID {
			allComments = append(allComments, comment)
		}
	}

	// Build a map for quick lookup
	commentMap := make(map[string]*models.Comment)
	for _, comment := range allComments {
		commentMap[comment.ID] = comment
	}

	// Find top-level comments (no parent)
	topLevelComments := make([]*models.Comment, 0)
	for _, comment := range allComments {
		if comment.ParentID == nil {
			topLevelComments = append(topLevelComments, comment)
		}
	}

	// Recursively collect all comments and their replies in a flattened list
	result := make([]*models.Comment, 0)
	visited := make(map[string]bool)

	var collectReplies func(commentID string)
	collectReplies = func(commentID string) {
		if visited[commentID] {
			return
		}
		visited[commentID] = true

		comment, exists := commentMap[commentID]
		if !exists {
			return
		}

		// Add the comment to result
		result = append(result, comment)

		// Find all direct replies to this comment
		for _, c := range allComments {
			if c.ParentID != nil && *c.ParentID == commentID && !visited[c.ID] {
				collectReplies(c.ID)
			}
		}
	}

	// Start from top-level comments and recursively collect all replies
	for _, topLevel := range topLevelComments {
		collectReplies(topLevel.ID)
	}

	// Include any orphaned comments (comments with invalid parent references)
	for _, comment := range allComments {
		if !visited[comment.ID] {
			result = append(result, comment)
		}
	}

	return result
}
