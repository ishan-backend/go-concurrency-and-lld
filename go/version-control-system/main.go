package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== Version Control System Demo ===\n")

	// Create a new repository
	repo := NewRepository("my-project", "John Doe")
	fmt.Printf("Created repository: %s\n", repo.Name)
	fmt.Printf("Current branch: %s\n\n", repo.GetCurrentBranch())

	// Add some files
	fmt.Println("--- Adding files ---")
	repo.AddFile("README.md", "# My Project\n\nThis is a sample project.")
	repo.AddFile("src/main.go", "package main\n\nfunc main() {\n\tprintln(\"Hello, World!\")\n}")
	repo.AddFile("src/utils.go", "package main\n\nfunc helper() {\n\t// Helper function\n}")
	repo.AddFile("config.json", "{\n\t\"version\": \"1.0.0\"\n}")

	fmt.Println("Files added:")
	for _, file := range repo.ListFiles() {
		fmt.Printf("  - %s\n", file)
	}
	fmt.Println()

	// Commit the changes
	fmt.Println("--- Committing changes ---")
	commitID1, err := repo.Commit("Add initial project files")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Commit created: %s\n\n", commitID1)

	// Update a file
	fmt.Println("--- Updating file ---")
	repo.UpdateFile("src/main.go", "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}")
	fmt.Println("Updated src/main.go")
	
	commitID2, err := repo.Commit("Update main.go to use fmt package")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Commit created: %s\n\n", commitID2)

	// Add another file
	fmt.Println("--- Adding more files ---")
	repo.AddFile("tests/test.go", "package tests\n\nfunc TestMain() {\n\t// Test code\n}")
	commitID3, err := repo.Commit("Add test file")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Commit created: %s\n\n", commitID3)

	// View commit history
	fmt.Println("--- Commit History ---")
	history := repo.GetCommitHistory()
	for i, commit := range history {
		fmt.Printf("%d. %s\n", i+1, commit.String())
	}
	fmt.Println()

	// Create a new branch
	fmt.Println("--- Creating feature branch ---")
	err = repo.CreateBranch("feature/new-feature")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Created branch: feature/new-feature")

	// Switch to the new branch
	err = repo.SwitchBranch("feature/new-feature")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Switched to branch: %s\n\n", repo.GetCurrentBranch())

	// Make changes on the feature branch
	fmt.Println("--- Making changes on feature branch ---")
	repo.AddFile("feature.go", "package main\n\nfunc newFeature() {\n\t// New feature implementation\n}")
	commitID4, err := repo.Commit("Add new feature")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Commit created: %s\n\n", commitID4)

	// View branches
	fmt.Println("--- Available Branches ---")
	branches := repo.GetBranches()
	for _, branch := range branches {
		current := ""
		if branch == repo.GetCurrentBranch() {
			current = " (current)"
		}
		fmt.Printf("  - %s%s\n", branch, current)
	}
	fmt.Println()

	// Switch back to main
	fmt.Println("--- Switching back to main branch ---")
	err = repo.SwitchBranch("main")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Switched to branch: %s\n", repo.GetCurrentBranch())
	
	// Verify feature.go doesn't exist on main
	_, err = repo.GetFileContent("feature.go")
	if err != nil {
		fmt.Println("feature.go does not exist on main branch (as expected)\n")
	}

	// Rollback to an earlier commit
	fmt.Println("--- Rolling back to commit 2 ---")
	err = repo.Rollback(commitID2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Rolled back to commit: %s\n\n", commitID2)

	// Verify the rollback
	fmt.Println("--- Files after rollback ---")
	files := repo.ListFiles()
	for _, file := range files {
		fmt.Printf("  - %s\n", file)
	}
	fmt.Println()

	// Check if test.go was removed (it was added in commit 3)
	_, err = repo.GetFileContent("tests/test.go")
	if err != nil {
		fmt.Println("tests/test.go was removed after rollback (as expected)")
	}

	// View updated commit history
	fmt.Println("\n--- Updated Commit History (after rollback) ---")
	history = repo.GetCommitHistory()
	for i, commit := range history {
		fmt.Printf("%d. %s\n", i+1, commit.String())
	}

	fmt.Println("\n=== Demo Complete ===")
}

// Helper function to demonstrate file operations
func demonstrateFileOperations(repo *Repository) {
	fmt.Println("\n=== File Operations Demo ===")

	// Add a file
	repo.AddFile("demo.txt", "Initial content")
	fmt.Println("Added demo.txt")

	// Read file content
	content, _ := repo.GetFileContent("demo.txt")
	fmt.Printf("Content: %s\n", content)

	// Update file
	repo.UpdateFile("demo.txt", "Updated content")
	content, _ = repo.GetFileContent("demo.txt")
	fmt.Printf("Updated content: %s\n", content)

	// Delete file
	repo.DeleteFile("demo.txt")
	fmt.Println("Deleted demo.txt")
}

