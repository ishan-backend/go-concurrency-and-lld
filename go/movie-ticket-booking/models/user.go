package models

import "sync/atomic"

var userIDGenerator int64 = 0

// User represents a user in the system
type User struct {
	ID    int64
	Name  string
	Email string
	Phone string
}

// NewUser creates a new User instance
func NewUser(name, email, phone string) *User {
	userID := atomic.AddInt64(&userIDGenerator, 1)
	return &User{
		ID:    userID,
		Name:  name,
		Email: email,
		Phone: phone,
	}
}
