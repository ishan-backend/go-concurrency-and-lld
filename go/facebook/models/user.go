package models

import (
	"sync/atomic"
	"time"
)

var userIDGenerator int64 = 0

// User represents a user in the system
type User struct {
	ID          int64
	FirstName   string
	LastName    string
	DateOfBirth time.Time
	Address     *Address
	Gender      Gender
	PhoneNumber string
	EmailID     string
	Followers   map[int64]bool // Set of follower IDs
	Posts       map[int]*Post  // Map of postID -> Post for O(1) lookup
	Head        *Post          // Dummy head node for doubly-linked list
	Tail        *Post          // Dummy tail node for doubly-linked list
}

// NewUser creates a new User instance
func NewUser(firstName, lastName string, dob time.Time) *User {
	userID := atomic.AddInt64(&userIDGenerator, 1)

	// Create dummy head and tail nodes for doubly-linked list
	head := NewPost(-10000)
	tail := NewPost(10000)
	head.Next = tail
	tail.Prev = head

	return &User{
		ID:          userID,
		FirstName:   firstName,
		LastName:    lastName,
		DateOfBirth: dob,
		Followers:   make(map[int64]bool),
		Posts:       make(map[int]*Post),
		Head:        head,
		Tail:        tail,
	}
}

// AddFollower adds a follower to the user
func (u *User) AddFollower(followerID int64) {
	u.Followers[followerID] = true
}

// RemoveFollower removes a follower from the user
func (u *User) RemoveFollower(followerID int64) {
	delete(u.Followers, followerID)
}

// GetFollowerIDs returns a slice of follower IDs
func (u *User) GetFollowerIDs() []int64 {
	followerIDs := make([]int64, 0, len(u.Followers))
	for id := range u.Followers {
		followerIDs = append(followerIDs, id)
	}
	return followerIDs
}

// GetPostIDs returns a slice of post IDs
func (u *User) GetPostIDs() []int {
	postIDs := make([]int, 0, len(u.Posts))
	for id := range u.Posts {
		postIDs = append(postIDs, id)
	}
	return postIDs
}
