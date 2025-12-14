package models

import (
	"sync/atomic"
)

var postKeyGenerator int64 = 0

// Post represents a post in the system
type Post struct {
	ID      int
	Key     int64 // Used for sorting/prioritization
	Content string
	Next    *Post
	Prev    *Post
}

// NewPost creates a new Post instance
func NewPost(id int) *Post {
	return &Post{
		ID:      id,
		Key:     atomic.AddInt64(&postKeyGenerator, 1),
		Content: "",
		Next:    nil,
		Prev:    nil,
	}
}

// SetContent sets the content of the post
func (p *Post) SetContent(content string) {
	p.Content = content
}
