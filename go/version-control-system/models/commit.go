package models

import (
	"fmt"
	"time"
)

// Commit represents a commit in the version control system
type Commit struct {
	ID          string
	Message     string
	Timestamp   time.Time
	Root        *DirectoryNode // Full snapshot of the file system
	ParentID    string          // ID of the parent commit (empty for initial commit)
	Author      string
}

// NewCommit creates a new commit
func NewCommit(id, message, author string, root *DirectoryNode, parentID string) *Commit {
	return &Commit{
		ID:        id,
		Message:   message,
		Timestamp: time.Now(),
		Root:      root.Copy(), // Store a snapshot
		ParentID:  parentID,
		Author:    author,
	}
}

// String returns a string representation of the commit
func (c *Commit) String() string {
	return fmt.Sprintf("Commit[%s] %s - %s (%s)", 
		c.ID[:8], c.Message, c.Author, c.Timestamp.Format("2006-01-02 15:04:05"))
}

