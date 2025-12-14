package service

import (
	"container/heap"
	"fmt"

	"github.com/ishanpandey/LowLevelDesign/go/facebook/models"
	"github.com/ishanpandey/LowLevelDesign/go/facebook/repository"
)

const (
	DefaultFeedSize = 1000 // Large number to represent "all posts"
	DefaultPageSize = 3
)

// postHeap is a max-heap for posts (sorted by key, most recent first)
type postHeap []*models.Post

func (h postHeap) Len() int           { return len(h) }
func (h postHeap) Less(i, j int) bool { return h[i].Key > h[j].Key } // Max heap
func (h postHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *postHeap) Push(x interface{}) {
	*h = append(*h, x.(*models.Post))
}

func (h *postHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// feedServiceImpl implements FeedService interface
// Following Single Responsibility Principle (SRP) - only handles feed generation logic
type feedServiceImpl struct {
	userRepo repository.UserRepository // Depends on abstraction (DIP)
}

// NewFeedService creates a new instance of FeedService
// Following Dependency Inversion Principle (DIP) - accepts repository interface
func NewFeedService(userRepo repository.UserRepository) FeedService {
	return &feedServiceImpl{
		userRepo: userRepo,
	}
}

// GetNewsFeed generates a news feed for a user containing posts from all followed users
// Uses a priority queue (max-heap) to efficiently merge posts from multiple followers
func (s *feedServiceImpl) GetNewsFeed(userID int64, topN *int) ([]*models.Post, error) {
	// Get the user
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return nil, fmt.Errorf("user %d not found: %w", userID, err)
	}

	// Determine the number of posts to return
	limit := DefaultFeedSize
	if topN != nil {
		limit = *topN
	}

	// Get all followers
	followerIDs := user.GetFollowerIDs()
	if len(followerIDs) == 0 {
		return []*models.Post{}, nil // Empty feed
	}

	// Initialize priority queue (max-heap) for merging posts
	pq := &postHeap{}
	heap.Init(pq)

	// Add the first post from each follower to the priority queue
	for _, followerID := range followerIDs {
		follower, err := s.userRepo.GetUser(followerID)
		if err != nil {
			continue // Skip if follower not found
		}

		firstPost := follower.Head.Next
		// Check if there are actual posts (not the dummy tail)
		if firstPost.ID != 10000 {
			heap.Push(pq, firstPost)
		}
	}

	// Merge posts using priority queue
	feed := make([]*models.Post, 0, limit)
	count := 0

	for pq.Len() > 0 && count < limit {
		// Get the post with highest key (most recent)
		curr := heap.Pop(pq).(*models.Post)
		feed = append(feed, curr)
		count++

		// Add the next post from the same user to the queue
		nextPost := curr.Next
		if nextPost.ID != 10000 { // Not the dummy tail
			heap.Push(pq, nextPost)
		}
	}

	return feed, nil
}

// GetFeedPaginated returns a paginated feed for a user
func (s *feedServiceImpl) GetFeedPaginated(userID int64, pageNumber int) ([]*models.Post, error) {
	if pageNumber < 1 {
		return nil, fmt.Errorf("page number must be >= 1")
	}

	// Get all posts
	allPosts, err := s.GetNewsFeed(userID, nil)
	if err != nil {
		return nil, err
	}

	// Calculate pagination
	start := (pageNumber - 1) * DefaultPageSize
	end := start + DefaultPageSize

	if start >= len(allPosts) {
		return []*models.Post{}, nil // Empty page
	}

	if end > len(allPosts) {
		end = len(allPosts)
	}

	return allPosts[start:end], nil
}

// PrintFeed is a helper function to print the feed (for demonstration)
func PrintFeed(posts []*models.Post, userID int64) {
	postIDs := make([]int, len(posts))
	for i, post := range posts {
		postIDs[i] = post.ID
	}
	fmt.Printf("User - %d has feed -> %v\n", userID, postIDs)
}
