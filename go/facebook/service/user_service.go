package service

import (
	"fmt"
	"time"

	"github.com/ishanpandey/LowLevelDesign/go/facebook/models"
	"github.com/ishanpandey/LowLevelDesign/go/facebook/repository"
)

// userServiceImpl implements UserService interface
// Following Single Responsibility Principle (SRP) - only handles user business logic
type userServiceImpl struct {
	userRepo repository.UserRepository // Depends on abstraction, not concrete type (DIP)
}

// NewUserService creates a new instance of UserService
// Following Dependency Inversion Principle (DIP) - accepts repository interface
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userServiceImpl{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (s *userServiceImpl) CreateUser(firstName, lastName string, dob string) (int64, error) {
	// Parse date of birth
	dateOfBirth, err := time.Parse("2006-01-02", dob)
	if err != nil {
		return 0, fmt.Errorf("invalid date format: %w", err)
	}

	// Create new user
	user := models.NewUser(firstName, lastName, dateOfBirth)

	// Save user to repository
	if err := s.userRepo.AddUser(user); err != nil {
		return 0, fmt.Errorf("failed to save user: %w", err)
	}

	return user.ID, nil
}

// GetUser retrieves a user by ID
func (s *userServiceImpl) GetUser(userID int64) (*models.User, error) {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// AddFollower adds a follower relationship
// When user A follows user B, user B gains user A as a follower
func (s *userServiceImpl) AddFollower(followedBy, follows int64) error {
	fmt.Printf("Adding follower = %d for user - %d\n", follows, followedBy)

	// Get the user who is being followed
	user, err := s.userRepo.GetUser(followedBy)
	if err != nil {
		return fmt.Errorf("user %d not found: %w", followedBy, err)
	}

	// Verify that the follower exists
	_, err = s.userRepo.GetUser(follows)
	if err != nil {
		return fmt.Errorf("follower user %d not found: %w", follows, err)
	}

	// Add follower
	user.AddFollower(follows)
	return nil
}

// RemoveFollower removes a follower relationship
func (s *userServiceImpl) RemoveFollower(userID, followerID int64) error {
	fmt.Printf("Removing follower = %d from user - %d\n", followerID, userID)

	// Get the user
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return fmt.Errorf("user %d not found: %w", userID, err)
	}

	// Remove follower
	user.RemoveFollower(followerID)
	return nil
}
