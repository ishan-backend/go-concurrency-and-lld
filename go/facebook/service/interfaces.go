package service

import "github.com/ishanpandey/LowLevelDesign/go/facebook/models"

// UserService defines the interface for user-related operations
// Following Interface Segregation Principle (ISP)
type UserService interface {
	CreateUser(firstName, lastName string, dob string) (int64, error)
	GetUser(userID int64) (*models.User, error)
	AddFollower(followedBy, follows int64) error
	RemoveFollower(userID, followerID int64) error
}

// PostService defines the interface for post-related operations
type PostService interface {
	AddPost(userID int64, postID int) error
	DeletePost(userID int64, postID int) error
	GetPost(userID int64, postID int) (*models.Post, error)
}

// FeedService defines the interface for feed generation operations
type FeedService interface {
	GetNewsFeed(userID int64, topN *int) ([]*models.Post, error)
	GetFeedPaginated(userID int64, pageNumber int) ([]*models.Post, error)
}
