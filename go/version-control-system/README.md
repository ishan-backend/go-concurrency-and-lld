# Version Control System

An in-memory Version Control System (VCS) implementation in Go that supports tracking file changes, committing updates, branching, and retrieving historical versions.

## Overview

This is a simplified version control system similar to Git, designed to demonstrate core VCS concepts including:

- **File Management**: Create, update, and delete files in a hierarchical directory structure
- **Commits**: Create commits with full snapshots of the repository state
- **Branching**: Create and switch between branches
- **Commit History**: View commit history for each branch
- **Rollback**: Revert the repository to any previous commit

## Features

### Functional Requirements

✅ **File Versioning**: Support creation and versioning of multiple files within a hierarchical directory structure  
✅ **Commits**: Allow users to commit changes across the entire repository (no staging area)  
✅ **Branching**: Support basic branching operations (create and switch branches)  
✅ **Commit History**: Maintain a commit history for each branch  
✅ **Rollback**: Allow users to rollback to any previous commit in the history  
✅ **Full Snapshots**: Store full snapshots of the file system at the time of each commit  
✅ **History Viewing**: Enable viewing the commit history of the repository  

### Design Principles

- **Modularity**: Clear separation between models and repository logic
- **Maintainability**: Object-oriented design with clear interfaces
- **Reliability**: Consistency and correctness across commits, branches, and rollbacks
- **Usability**: Simple interface for core operations

## Architecture

### Core Components

#### Models

- **FileNode**: Represents a file with name, content, path, and modification time
- **DirectoryNode**: Represents a directory with files and subdirectories (hierarchical structure)
- **Commit**: Represents a commit with ID, message, timestamp, full snapshot, parent commit, and author
- **Branch**: Represents a branch with name, head commit, and commit history

#### Repository

The `Repository` class is the main entry point that manages:
- File system state (root directory)
- Branches and current branch
- Commits and commit history
- All VCS operations

## Usage

### Basic Operations

#### Create a Repository

```go
repo := NewRepository("my-project", "John Doe")
```

#### Add Files

```go
repo.AddFile("README.md", "# My Project")
repo.AddFile("src/main.go", "package main\n\nfunc main() {}")
repo.AddFile("config/app.json", `{"version": "1.0.0"}`)
```

#### Update Files

```go
repo.UpdateFile("README.md", "# Updated Project")
```

#### Delete Files

```go
repo.DeleteFile("old-file.txt")
```

#### Commit Changes

```go
commitID, err := repo.Commit("Add initial project files")
```

#### Create a Branch

```go
err := repo.CreateBranch("feature/new-feature")
```

#### Switch Branch

```go
err := repo.SwitchBranch("feature/new-feature")
```

#### Rollback to a Commit

```go
err := repo.Rollback(commitID)
```

#### View Commit History

```go
history := repo.GetCommitHistory()
for _, commit := range history {
    fmt.Println(commit.String())
}
```

#### List Files

```go
files := repo.ListFiles()
for _, file := range files {
    fmt.Println(file)
}
```

#### Get File Content

```go
content, err := repo.GetFileContent("src/main.go")
```

### Example Workflow

```go
package main

import (
    "fmt"
)

func main() {
    // Create repository
    repo := NewRepository("my-project", "Alice")
    
    // Add files
    repo.AddFile("main.go", "package main\n\nfunc main() {}")
    repo.AddFile("README.md", "# Project")
    
    // Commit
    commitID1, _ := repo.Commit("Initial commit")
    
    // Update file
    repo.UpdateFile("main.go", "package main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello\")\n}")
    commitID2, _ := repo.Commit("Add fmt import")
    
    // Create and switch to feature branch
    repo.CreateBranch("feature")
    repo.SwitchBranch("feature")
    
    // Make changes on feature branch
    repo.AddFile("feature.go", "package main\n\nfunc feature() {}")
    commitID3, _ := repo.Commit("Add feature")
    
    // Switch back to main
    repo.SwitchBranch("main")
    
    // Rollback
    repo.Rollback(commitID1)
    
    // View history
    history := repo.GetCommitHistory()
    for _, commit := range history {
        fmt.Println(commit.String())
    }
}
```

## Project Structure

```
version-control-system/
├── models/
│   ├── file.go          # FileNode model
│   ├── directory.go     # DirectoryNode model
│   ├── commit.go        # Commit model
│   ├── branch.go        # Branch model
│   └── utils.go         # Utility functions
├── repository.go        # Main repository implementation
├── main.go              # Example usage and demo
├── go.mod               # Go module file
└── README.md            # This file
```

## Design Decisions

### Full Snapshots vs Deltas

The system stores **full snapshots** of the file system for each commit rather than deltas (differences). This simplifies the implementation but uses more memory. In a production system, deltas would be preferred for efficiency.

### In-Memory Storage

All data is stored in memory. For persistence, you would need to add serialization/deserialization to disk.

### No Staging Area

Unlike Git, this system commits all changes at once. There's no staging area where you can selectively add files to a commit.

### Branch Implementation

Branches maintain their own commit history. When switching branches, the entire file system state is restored from the branch's head commit.

### Rollback Behavior

Rollback removes commits after the target commit from the branch history. This is a simplified approach - in real systems, you might want to create a new commit that reverts changes instead.

## Running the Demo

```bash
cd version-control-system
go run .
```

This will execute the demo program that demonstrates:
1. Creating a repository
2. Adding and updating files
3. Creating commits
4. Creating and switching branches
5. Rolling back to previous commits
6. Viewing commit history

## Testing

To test the system, you can create your own test file:

```go
package main

import (
    "testing"
)

func TestRepository(t *testing.T) {
    repo := NewRepository("test", "Test User")
    
    // Test adding files
    err := repo.AddFile("test.txt", "content")
    if err != nil {
        t.Errorf("Failed to add file: %v", err)
    }
    
    // Test committing
    commitID, err := repo.Commit("Test commit")
    if err != nil {
        t.Errorf("Failed to commit: %v", err)
    }
    
    // Test rollback
    err = repo.Rollback(commitID)
    if err != nil {
        t.Errorf("Failed to rollback: %v", err)
    }
}
```

## Limitations

1. **In-Memory Only**: No persistence to disk
2. **Full Snapshots**: Uses more memory than delta-based systems
3. **No Merging**: Branch merging is not implemented
4. **No Conflict Resolution**: Not applicable since merging isn't supported
5. **No Remote Operations**: No push/pull functionality
6. **Single User**: No multi-user or collaboration features
7. **No File Permissions**: All files are treated equally

## Future Enhancements

- [ ] Add delta-based storage for efficiency
- [ ] Implement branch merging
- [ ] Add persistence to disk
- [ ] Support for file permissions and metadata
- [ ] Add diff viewing between commits
- [ ] Implement tags for marking important commits
- [ ] Add remote repository support
- [ ] Implement conflict resolution for merging

## References

Design inspired by: [Design Version Control System - AlgoMaster](https://algomaster.io/learn/lld/design-version-control-system)

## License

This project is part of the LowLevelDesign collection.

