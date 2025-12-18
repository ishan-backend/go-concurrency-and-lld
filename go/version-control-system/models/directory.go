package models

import (
	"time"
)

// DirectoryNode represents a directory in the version control system
type DirectoryNode struct {
	Name      string
	Path      string
	Files     map[string]*FileNode      // filename -> FileNode
	SubDirs   map[string]*DirectoryNode // dirname -> DirectoryNode
	Modified  time.Time
}

// NewDirectory creates a new directory node
func NewDirectory(name, path string) *DirectoryNode {
	return &DirectoryNode{
		Name:     name,
		Path:     path,
		Files:    make(map[string]*FileNode),
		SubDirs:  make(map[string]*DirectoryNode),
		Modified: time.Now(),
	}
}

// AddFile adds a file to the directory
func (d *DirectoryNode) AddFile(file *FileNode) {
	d.Files[file.Name] = file
	d.Modified = time.Now()
}

// AddSubDirectory adds a subdirectory to the directory
func (d *DirectoryNode) AddSubDirectory(dir *DirectoryNode) {
	d.SubDirs[dir.Name] = dir
	d.Modified = time.Now()
}

// GetFile retrieves a file by name
func (d *DirectoryNode) GetFile(name string) (*FileNode, bool) {
	file, exists := d.Files[name]
	return file, exists
}

// GetSubDirectory retrieves a subdirectory by name
func (d *DirectoryNode) GetSubDirectory(name string) (*DirectoryNode, bool) {
	dir, exists := d.SubDirs[name]
	return dir, exists
}

// Copy creates a deep copy of the directory and all its contents
func (d *DirectoryNode) Copy() *DirectoryNode {
	copy := &DirectoryNode{
		Name:     d.Name,
		Path:     d.Path,
		Files:    make(map[string]*FileNode),
		SubDirs:  make(map[string]*DirectoryNode),
		Modified: d.Modified,
	}

	// Copy all files
	for name, file := range d.Files {
		copy.Files[name] = file.Copy()
	}

	// Recursively copy all subdirectories
	for name, subDir := range d.SubDirs {
		copy.SubDirs[name] = subDir.Copy()
	}

	return copy
}

// FindFile searches for a file by path (e.g., "src/main.go" or "main.go")
func (d *DirectoryNode) FindFile(path string) (*FileNode, bool) {
	parts := splitPath(path)
	if len(parts) == 0 {
		return nil, false
	}

	current := d
	for i := 0; i < len(parts)-1; i++ {
		subDir, exists := current.SubDirs[parts[i]]
		if !exists {
			return nil, false
		}
		current = subDir
	}

	file, exists := current.Files[parts[len(parts)-1]]
	return file, exists
}

// FindOrCreateDirectory finds or creates a directory path
func (d *DirectoryNode) FindOrCreateDirectory(path string) *DirectoryNode {
	parts := splitPath(path)
	if len(parts) == 0 {
		return d
	}

	current := d
	for _, part := range parts {
		if subDir, exists := current.SubDirs[part]; exists {
			current = subDir
		} else {
			newDir := NewDirectory(part, current.Path+"/"+part)
			current.AddSubDirectory(newDir)
			current = newDir
		}
	}
	return current
}

// FindDirectory finds a directory by path
func (d *DirectoryNode) FindDirectory(path string) (*DirectoryNode, bool) {
	parts := splitPath(path)
	if len(parts) == 0 {
		return d, true
	}

	current := d
	for _, part := range parts {
		subDir, exists := current.SubDirs[part]
		if !exists {
			return nil, false
		}
		current = subDir
	}
	return current, true
}

// String returns a string representation of the directory structure
func (d *DirectoryNode) String() string {
	return d.stringRecursive("")
}

func (d *DirectoryNode) stringRecursive(indent string) string {
	result := indent + d.Name + "/\n"
	
	for _, file := range d.Files {
		result += indent + "  " + file.Name + "\n"
	}
	
	for _, subDir := range d.SubDirs {
		result += subDir.stringRecursive(indent + "  ")
	}
	
	return result
}

// splitPath splits a path string into parts
func splitPath(path string) []string {
	if path == "" {
		return []string{}
	}
	
	parts := []string{}
	current := ""
	for _, char := range path {
		if char == '/' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

