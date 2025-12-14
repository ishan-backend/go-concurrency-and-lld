package repository

import (
	"errors"
	"sync"

	"github.com/ishanpandey/LowLevelDesign/go/facebook/models"
)

// userRepositoryImpl implements UserRepository interface
// Following Single Responsibility Principle (SRP) - only handles user data persistence
type userRepositoryImpl struct {
	userMap map[int64]*models.User
	mu      sync.RWMutex // For thread-safe operations
}

// NewUserRepository creates a new instance of UserRepository
// Following Dependency Inversion Principle (DIP) - returns interface, not concrete type
func NewUserRepository() UserRepository {
	return &userRepositoryImpl{
		userMap: make(map[int64]*models.User),
	}
}

// AddUser adds a user to the repository
func (r *userRepositoryImpl) AddUser(user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	if user.ID == 0 {
		return errors.New("user ID cannot be zero")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.userMap[user.ID] = user
	return nil
}

// GetUser retrieves a user by ID
func (r *userRepositoryImpl) GetUser(id int64) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.userMap[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}
