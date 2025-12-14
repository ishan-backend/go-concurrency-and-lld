package models

import (
	"time"

	"github.com/google/uuid"
)

// Post represents a question/post in Stack Overflow
type Post struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  string    `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Comment represents a comment or reply in Stack Overflow
type Comment struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	AuthorID  string    `json:"author_id"`
	Content   string    `json:"content"`
	ParentID  *string   `json:"parent_id,omitempty"` // nil for top-level comments, set for replies
	Upvotes   int       `json:"upvotes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewPost creates a new post
func NewPost(title, content, authorID string) *Post {
	now := time.Now()
	return &Post{
		ID:        uuid.New().String(),
		Title:     title,
		Content:   content,
		AuthorID:  authorID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewComment creates a new comment
func NewComment(postID, authorID, content string, parentID *string) *Comment {
	now := time.Now()
	return &Comment{
		ID:        uuid.New().String(),
		PostID:    postID,
		AuthorID:  authorID,
		Content:   content,
		ParentID:  parentID,
		Upvotes:   0,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// IsTopLevel checks if the comment is a top-level comment (not a reply)
func (c *Comment) IsTopLevel() bool {
	return c.ParentID == nil
}

// Upvote increments the upvote count
func (c *Comment) Upvote() {
	c.Upvotes++
}
