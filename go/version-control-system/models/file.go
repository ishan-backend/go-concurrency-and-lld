package models

import (
	"time"
)

// FileNode represents a file in the version control system
type FileNode struct {
	Name     string
	Content  string
	Path     string
	Modified time.Time
}

// NewFile creates a new file node
func NewFile(name, path, content string) *FileNode {
	return &FileNode{
		Name:     name,
		Content:  content,
		Path:     path,
		Modified: time.Now(),
	}
}

// Copy creates a deep copy of the file
func (f *FileNode) Copy() *FileNode {
	return &FileNode{
		Name:     f.Name,
		Content:  f.Content,
		Path:     f.Path,
		Modified: f.Modified,
	}
}

