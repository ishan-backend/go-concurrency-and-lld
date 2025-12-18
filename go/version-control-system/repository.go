package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
	"vcs/models"
)

// Repository represents the version control system repository
type Repository struct {
	Name           string
	Root           *models.DirectoryNode
	Branches       map[string]*models.Branch
	CurrentBranch  string
	Commits        map[string]*models.Commit // commitID -> Commit
	Author         string
}

// NewRepository creates a new repository
func NewRepository(name, author string) *Repository {
	root := models.NewDirectory("root", "/")
	initialCommitID := generateCommitID("initial", author)
	initialCommit := models.NewCommit(initialCommitID, "Initial commit", author, root, "")
	
	mainBranch := models.NewBranch("main", initialCommitID)
	
	repo := &Repository{
		Name:          name,
		Root:          root,
		Branches:      make(map[string]*models.Branch),
		CurrentBranch: "main",
		Commits:       make(map[string]*models.Commit),
		Author:        author,
	}
	
	repo.Branches["main"] = mainBranch
	repo.Commits[initialCommitID] = initialCommit
	
	return repo
}

// AddFile adds a file to the repository
func (r *Repository) AddFile(path, content string) error {
	parts := splitPath(path)
	if len(parts) == 0 {
		return errors.New("invalid file path")
	}

	fileName := parts[len(parts)-1]
	dirPath := ""
	if len(parts) > 1 {
		dirPath = joinPath(parts[:len(parts)-1])
	}

	var targetDir *models.DirectoryNode
	if dirPath == "" {
		targetDir = r.Root
	} else {
		targetDir = r.Root.FindOrCreateDirectory(dirPath)
	}

	file := models.NewFile(fileName, path, content)
	targetDir.AddFile(file)
	return nil
}

// UpdateFile updates an existing file in the repository
func (r *Repository) UpdateFile(path, content string) error {
	file, exists := r.Root.FindFile(path)
	if !exists {
		return fmt.Errorf("file not found: %s", path)
	}
	file.Content = content
	file.Modified = time.Now()
	return nil
}

// DeleteFile deletes a file from the repository
func (r *Repository) DeleteFile(path string) error {
	parts := splitPath(path)
	if len(parts) == 0 {
		return errors.New("invalid file path")
	}

	fileName := parts[len(parts)-1]
	dirPath := ""
	if len(parts) > 1 {
		dirPath = joinPath(parts[:len(parts)-1])
	}

	var targetDir *models.DirectoryNode
	if dirPath == "" {
		targetDir = r.Root
	} else {
		dir, exists := r.Root.FindDirectory(dirPath)
		if !exists {
			return fmt.Errorf("directory not found: %s", dirPath)
		}
		targetDir = dir
	}

	if _, exists := targetDir.Files[fileName]; !exists {
		return fmt.Errorf("file not found: %s", path)
	}

	delete(targetDir.Files, fileName)
	return nil
}

// Commit creates a new commit with the current state of the repository
func (r *Repository) Commit(message string) (string, error) {
	if message == "" {
		return "", errors.New("commit message cannot be empty")
	}

	currentBranch := r.Branches[r.CurrentBranch]
	if currentBranch == nil {
		return "", fmt.Errorf("branch not found: %s", r.CurrentBranch)
	}

	parentID := currentBranch.Head
	commitID := generateCommitID(message+r.Root.String(), r.Author)
	
	commit := models.NewCommit(commitID, message, r.Author, r.Root, parentID)
	
	r.Commits[commitID] = commit
	currentBranch.AddCommit(commitID)
	
	return commitID, nil
}

// CreateBranch creates a new branch from the current branch
func (r *Repository) CreateBranch(branchName string) error {
	if branchName == "" {
		return errors.New("branch name cannot be empty")
	}

	if _, exists := r.Branches[branchName]; exists {
		return fmt.Errorf("branch already exists: %s", branchName)
	}

	currentBranch := r.Branches[r.CurrentBranch]
	if currentBranch == nil {
		return fmt.Errorf("current branch not found: %s", r.CurrentBranch)
	}

	newBranch := &models.Branch{
		Name:    branchName,
		Head:    currentBranch.Head,
		Commits: make([]string, len(currentBranch.Commits)),
	}
	copy(newBranch.Commits, currentBranch.Commits)
	
	r.Branches[branchName] = newBranch
	return nil
}

// SwitchBranch switches to a different branch and restores its state
func (r *Repository) SwitchBranch(branchName string) error {
	if _, exists := r.Branches[branchName]; !exists {
		return fmt.Errorf("branch not found: %s", branchName)
	}

	if r.CurrentBranch == branchName {
		return nil // Already on this branch
	}

	branch := r.Branches[branchName]
	if branch.Head == "" {
		return fmt.Errorf("branch has no commits: %s", branchName)
	}

	commit, exists := r.Commits[branch.Head]
	if !exists {
		return fmt.Errorf("commit not found: %s", branch.Head)
	}

	// Restore the state from the commit
	r.Root = commit.Root.Copy()
	r.CurrentBranch = branchName
	
	return nil
}

// Rollback reverts the repository to a specific commit
func (r *Repository) Rollback(commitID string) error {
	commit, exists := r.Commits[commitID]
	if !exists {
		return fmt.Errorf("commit not found: %s", commitID)
	}

	// Check if commit belongs to current branch
	currentBranch := r.Branches[r.CurrentBranch]
	found := false
	for _, id := range currentBranch.Commits {
		if id == commitID {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("commit %s does not belong to current branch %s", commitID, r.CurrentBranch)
	}

	// Restore the state from the commit
	r.Root = commit.Root.Copy()
	
	// Update branch head to this commit
	currentBranch.Head = commitID
	
	// Remove commits after this one from the branch
	index := -1
	for i, id := range currentBranch.Commits {
		if id == commitID {
			index = i
			break
		}
	}
	if index >= 0 {
		currentBranch.Commits = currentBranch.Commits[:index+1]
	}
	
	return nil
}

// GetCommitHistory returns the commit history for the current branch
func (r *Repository) GetCommitHistory() []*models.Commit {
	currentBranch := r.Branches[r.CurrentBranch]
	if currentBranch == nil {
		return []*models.Commit{}
	}

	history := make([]*models.Commit, 0, len(currentBranch.Commits))
	for _, commitID := range currentBranch.Commits {
		if commit, exists := r.Commits[commitID]; exists {
			history = append(history, commit)
		}
	}
	return history
}

// GetBranches returns all branch names
func (r *Repository) GetBranches() []string {
	branches := make([]string, 0, len(r.Branches))
	for name := range r.Branches {
		branches = append(branches, name)
	}
	return branches
}

// GetCurrentBranch returns the name of the current branch
func (r *Repository) GetCurrentBranch() string {
	return r.CurrentBranch
}

// GetFileContent retrieves the content of a file
func (r *Repository) GetFileContent(path string) (string, error) {
	file, exists := r.Root.FindFile(path)
	if !exists {
		return "", fmt.Errorf("file not found: %s", path)
	}
	return file.Content, nil
}

// ListFiles lists all files in the repository
func (r *Repository) ListFiles() []string {
	return r.listFilesRecursive(r.Root, "")
}

func (r *Repository) listFilesRecursive(dir *models.DirectoryNode, prefix string) []string {
	files := make([]string, 0)
	
	for _, file := range dir.Files {
		path := prefix + "/" + file.Name
		if prefix == "" {
			path = file.Name
		}
		files = append(files, path)
	}
	
	for name, subDir := range dir.SubDirs {
		subPrefix := prefix + "/" + name
		if prefix == "" {
			subPrefix = name
		}
		files = append(files, r.listFilesRecursive(subDir, subPrefix)...)
	}
	
	return files
}

// Helper functions

func generateCommitID(data, author string) string {
	hash := sha256.Sum256([]byte(data + author + fmt.Sprintf("%d", time.Now().Unix())))
	return hex.EncodeToString(hash[:])[:16]
}

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

func joinPath(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result += "/" + parts[i]
	}
	return result
}

