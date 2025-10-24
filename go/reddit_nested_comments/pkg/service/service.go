package service

import (
	"errors"
	"reddit-comments/pkg/models"
	"reddit-comments/pkg/repository"
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

// User operations
func (s *Service) CreateUser(username, email string) (*models.User, error) {
	if username == "" || email == "" {
		return nil, errors.New("username and email are required")
	}

	user := models.NewUser(username, email)
	err := s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetUser(id string) (*models.User, error) {
	return s.repo.GetUser(id)
}

func (s *Service) GetAllUsers() []*models.User {
	return s.repo.GetAllUsers()
}

// Post operations
func (s *Service) CreatePost(title, content, authorID string) (*models.Post, error) {
	if title == "" || content == "" || authorID == "" {
		return nil, errors.New("title, content, and author ID are required")
	}

	// Verify author exists
	_, err := s.repo.GetUser(authorID)
	if err != nil {
		return nil, errors.New("author not found")
	}

	post := models.NewPost(title, content, authorID)
	err = s.repo.CreatePost(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *Service) GetPost(id string) (*models.Post, error) {
	return s.repo.GetPost(id)
}

func (s *Service) GetAllPosts() []*models.Post {
	return s.repo.GetAllPosts()
}

// Comment operations
func (s *Service) CreateComment(postID, authorID, content string, parentID *string) (*models.Comment, error) {
	if postID == "" || authorID == "" || content == "" {
		return nil, errors.New("post ID, author ID, and content are required")
	}

	// Verify post exists
	_, err := s.repo.GetPost(postID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	// Verify author exists
	_, err = s.repo.GetUser(authorID)
	if err != nil {
		return nil, errors.New("author not found")
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

	// Update post comment count
	err = s.repo.UpdatePostCommentCount(postID, 1)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *Service) GetComment(id string) (*models.Comment, error) {
	return s.repo.GetComment(id)
}

func (s *Service) GetCommentsByPost(postID string) []*models.Comment {
	return s.repo.GetCommentsByPost(postID)
}

func (s *Service) GetTopLevelComments(postID string) []*models.Comment {
	return s.repo.GetTopLevelComments(postID)
}

func (s *Service) GetAllComments() []*models.Comment {
	return s.repo.GetAllComments()
}

// GetCommentPath returns the path from root to the specified comment
func (s *Service) GetCommentPath(commentID string) (*models.CommentPath, error) {
	return s.repo.GetCommentPath(commentID)
}

// GetCommentTree returns the complete comment tree for a post
func (s *Service) GetCommentTree(postID string) ([]models.CommentTree, error) {
	// Verify post exists
	_, err := s.repo.GetPost(postID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	return s.repo.BuildCommentTree(postID), nil
}

// GetCommentAncestors returns all parent comments for a given comment
func (s *Service) GetCommentAncestors(commentID string) ([]*models.Comment, error) {
	path, err := s.repo.GetCommentPath(commentID)
	if err != nil {
		return nil, err
	}

	// Return all comments except the last one (which is the comment itself)
	if len(path.Comments) <= 1 {
		return []*models.Comment{}, nil
	}

	return path.Comments[:len(path.Comments)-1], nil
}

// GetCommentWithReplies returns a comment with all its replies
func (s *Service) GetCommentWithReplies(commentID string) (*models.Comment, error) {
	comment, err := s.repo.GetComment(commentID)
	if err != nil {
		return nil, err
	}

	// The repository already maintains the replies structure
	return comment, nil
}

// GetCommentDepth returns the depth of a comment in the tree
func (s *Service) GetCommentDepth(commentID string) (int, error) {
	path, err := s.repo.GetCommentPath(commentID)
	if err != nil {
		return 0, err
	}

	return len(path.Comments) - 1, nil
}

// GetCommentStats returns statistics about comments for a post
func (s *Service) GetCommentStats(postID string) (map[string]interface{}, error) {
	// Verify post exists
	_, err := s.repo.GetPost(postID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	comments := s.repo.GetCommentsByPost(postID)
	topLevelComments := s.repo.GetTopLevelComments(postID)

	stats := map[string]interface{}{
		"total_comments":     len(comments),
		"top_level_comments": len(topLevelComments),
		"nested_comments":    len(comments) - len(topLevelComments),
	}

	// Calculate max depth
	maxDepth := 0
	for _, comment := range comments {
		depth, _ := s.GetCommentDepth(comment.ID)
		if depth > maxDepth {
			maxDepth = depth
		}
	}
	stats["max_depth"] = maxDepth

	return stats, nil
}
