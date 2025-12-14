package service

import (
	"errors"
	"stackoverflow/pkg/models"
	"stackoverflow/pkg/repository"
)

// Service handles business logic
type Service struct {
	repo *repository.Repository
}

// NewService creates a new service instance
func NewService(repo *repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreatePost creates a new post/question
func (s *Service) CreatePost(title, content, authorID string) (*models.Post, error) {
	if title == "" || content == "" || authorID == "" {
		return nil, errors.New("title, content, and author ID are required")
	}

	post := models.NewPost(title, content, authorID)
	err := s.repo.CreatePost(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

// AddComment creates a new comment or reply
// If parentID is nil, it's a top-level comment
// If parentID is provided, it's a reply to that comment
func (s *Service) AddComment(postID, authorID, content string, parentID *string) (*models.Comment, error) {
	if postID == "" || authorID == "" || content == "" {
		return nil, errors.New("post ID, author ID, and content are required")
	}

	// Verify post exists
	_, err := s.repo.GetPost(postID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	// If parent ID is provided, verify parent comment exists
	if parentID != nil {
		parentComment, err := s.repo.GetComment(*parentID)
		if err != nil {
			return nil, errors.New("parent comment not found")
		}

		// Verify parent comment belongs to the same post
		if parentComment.PostID != postID {
			return nil, errors.New("parent comment does not belong to the specified post")
		}
	}

	comment := models.NewComment(postID, authorID, content, parentID)
	err = s.repo.CreateComment(comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// UpvoteComment increments the upvote count for a comment
func (s *Service) UpvoteComment(commentID string) (*models.Comment, error) {
	if commentID == "" {
		return nil, errors.New("comment ID is required")
	}

	// Verify comment exists
	_, err := s.repo.GetComment(commentID)
	if err != nil {
		return nil, errors.New("comment not found")
	}

	// Upvote the comment
	err = s.repo.UpvoteComment(commentID)
	if err != nil {
		return nil, err
	}

	// Return updated comment
	return s.repo.GetComment(commentID)
}

// GetAllCommentsRecursively returns all comments for a post, including all nested replies, in a flattened list
func (s *Service) GetAllCommentsRecursively(postID string) ([]*models.Comment, error) {
	if postID == "" {
		return nil, errors.New("post ID is required")
	}

	// Verify post exists
	_, err := s.repo.GetPost(postID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	// Get all comments recursively
	comments := s.repo.GetAllCommentsRecursively(postID)
	return comments, nil
}
