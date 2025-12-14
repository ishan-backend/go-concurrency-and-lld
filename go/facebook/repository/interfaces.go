package repository

import "github.com/ishanpandey/LowLevelDesign/go/facebook/models"

// UserRepository defines the interface for user data operations
// Following Interface Segregation Principle (ISP)
type UserRepository interface {
	AddUser(user *models.User) error
	GetUser(id int64) (*models.User, error)
}
