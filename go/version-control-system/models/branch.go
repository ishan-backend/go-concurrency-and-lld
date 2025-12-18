package models

// Branch represents a branch in the version control system
type Branch struct {
	Name      string
	Head      string // ID of the latest commit
	Commits   []string // List of commit IDs in order
}

// NewBranch creates a new branch
func NewBranch(name, initialCommitID string) *Branch {
	return &Branch{
		Name:    name,
		Head:    initialCommitID,
		Commits: []string{initialCommitID},
	}
}

// AddCommit adds a commit to the branch
func (b *Branch) AddCommit(commitID string) {
	b.Commits = append(b.Commits, commitID)
	b.Head = commitID
}

// GetCommitHistory returns the list of commit IDs in this branch
func (b *Branch) GetCommitHistory() []string {
	return b.Commits
}

