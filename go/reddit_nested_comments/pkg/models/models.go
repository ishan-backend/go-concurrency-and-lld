package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Post represents a post that can have comments
type Post struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	AuthorID     string    `json:"author_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CommentCount int       `json:"comment_count"`
}

// Comment represents a comment that can have replies
type Comment struct {
	ID        string     `json:"id"`
	PostID    string     `json:"post_id"`
	AuthorID  string     `json:"author_id"`
	Content   string     `json:"content"`
	ParentID  *string    `json:"parent_id,omitempty"` // nil for top-level comments
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Replies   []*Comment `json:"replies,omitempty"`
	Depth     int        `json:"depth"` // Depth level for nested display
}

// CommentTree represents the hierarchical structure of comments
type CommentTree struct {
	Comment *Comment      `json:"comment"`
	Replies []CommentTree `json:"replies"`
}

// CommentPath represents the path from root to a specific comment
type CommentPath struct {
	Comments []*Comment `json:"comments"`
}

// NewUser creates a new user
func NewUser(username, email string) *User {
	return &User{
		ID:       uuid.New().String(),
		Username: username,
		Email:    email,
	}
}

// NewPost creates a new post
func NewPost(title, content, authorID string) *Post {
	now := time.Now()
	return &Post{
		ID:           uuid.New().String(),
		Title:        title,
		Content:      content,
		AuthorID:     authorID,
		CreatedAt:    now,
		UpdatedAt:    now,
		CommentCount: 0,
	}
}

// NewComment creates a new comment
func NewComment(postID, authorID, content string, parentID *string) *Comment {
	now := time.Now()
	depth := 0
	if parentID != nil {
		depth = 1 // Will be calculated properly by the service
	}

	return &Comment{
		ID:        uuid.New().String(),
		PostID:    postID,
		AuthorID:  authorID,
		Content:   content,
		ParentID:  parentID,
		CreatedAt: now,
		UpdatedAt: now,
		Replies:   make([]*Comment, 0),
		Depth:     depth,
	}
}

// AddReply adds a reply to a comment
func (c *Comment) AddReply(reply *Comment) {
	c.Replies = append(c.Replies, reply)
}

// IsTopLevel checks if the comment is a top-level comment
func (c *Comment) IsTopLevel() bool {
	return c.ParentID == nil
}

// SetDepth sets the depth of the comment in the tree
func (c *Comment) SetDepth(depth int) {
	c.Depth = depth
	for _, reply := range c.Replies {
		reply.SetDepth(depth + 1)
	}
}
