package repository

import (
	"errors"
	"reddit-comments/pkg/models"
	"sync"
)

// Repository handles data storage and retrieval
type Repository struct {
	users    map[string]*models.User
	posts    map[string]*models.Post
	comments map[string]*models.Comment
	mu       sync.RWMutex
}

// NewRepository creates a new repository instance
func NewRepository() *Repository {
	return &Repository{
		users:    make(map[string]*models.User),
		posts:    make(map[string]*models.Post),
		comments: make(map[string]*models.Comment),
	}
}

// User operations
func (r *Repository) CreateUser(user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return errors.New("user already exists")
	}

	r.users[user.ID] = user
	return nil
}

func (r *Repository) GetUser(id string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (r *Repository) GetAllUsers() []*models.User {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*models.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users
}

// Post operations
func (r *Repository) CreatePost(post *models.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.posts[post.ID]; exists {
		return errors.New("post already exists")
	}

	r.posts[post.ID] = post
	return nil
}

func (r *Repository) GetPost(id string) (*models.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, exists := r.posts[id]
	if !exists {
		return nil, errors.New("post not found")
	}

	return post, nil
}

func (r *Repository) GetAllPosts() []*models.Post {
	r.mu.RLock()
	defer r.mu.RUnlock()

	posts := make([]*models.Post, 0, len(r.posts))
	for _, post := range r.posts {
		posts = append(posts, post)
	}

	return posts
}

func (r *Repository) UpdatePostCommentCount(postID string, delta int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	post, exists := r.posts[postID]
	if !exists {
		return errors.New("post not found")
	}

	post.CommentCount += delta
	return nil
}

// Comment operations
func (r *Repository) CreateComment(comment *models.Comment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.comments[comment.ID]; exists {
		return errors.New("comment already exists")
	}

	r.comments[comment.ID] = comment

	// If this is a reply, add it to parent's replies
	if comment.ParentID != nil {
		if parent, exists := r.comments[*comment.ParentID]; exists {
			parent.AddReply(comment)
		}
	}

	return nil
}

func (r *Repository) GetComment(id string) (*models.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	comment, exists := r.comments[id]
	if !exists {
		return nil, errors.New("comment not found")
	}

	return comment, nil
}

func (r *Repository) GetCommentsByPost(postID string) []*models.Comment {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var comments []*models.Comment
	for _, comment := range r.comments {
		if comment.PostID == postID {
			comments = append(comments, comment)
		}
	}

	return comments
}

func (r *Repository) GetTopLevelComments(postID string) []*models.Comment {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var comments []*models.Comment
	for _, comment := range r.comments {
		if comment.PostID == postID && comment.IsTopLevel() {
			comments = append(comments, comment)
		}
	}

	return comments
}

func (r *Repository) GetAllComments() []*models.Comment {
	r.mu.RLock()
	defer r.mu.RUnlock()

	comments := make([]*models.Comment, 0, len(r.comments))
	for _, comment := range r.comments {
		comments = append(comments, comment)
	}

	return comments
}

// GetCommentPath returns the path from root to the specified comment
func (r *Repository) GetCommentPath(commentID string) (*models.CommentPath, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	comment, exists := r.comments[commentID]
	if !exists {
		return nil, errors.New("comment not found")
	}

	path := &models.CommentPath{
		Comments: make([]*models.Comment, 0),
	}

	// Build path from comment to root
	current := comment
	for current != nil {
		path.Comments = append([]*models.Comment{current}, path.Comments...)

		if current.ParentID != nil {
			var found bool
			current, found = r.comments[*current.ParentID]
			if !found {
				break
			}
		} else {
			break
		}
	}

	return path, nil
}

// BuildCommentTree builds the complete comment tree for a post
func (r *Repository) BuildCommentTree(postID string) []models.CommentTree {
	r.mu.RLock()
	defer r.mu.RUnlock()

	topLevelComments := r.GetTopLevelComments(postID)
	trees := make([]models.CommentTree, 0, len(topLevelComments))

	for _, comment := range topLevelComments {
		tree := r.buildTreeFromComment(comment)
		trees = append(trees, tree)
	}

	return trees
}

// buildTreeFromComment recursively builds a comment tree
func (r *Repository) buildTreeFromComment(comment *models.Comment) models.CommentTree {
	tree := models.CommentTree{
		Comment: comment,
		Replies: make([]models.CommentTree, 0),
	}

	// Set depth for this comment and its replies
	comment.SetDepth(0)

	// Recursively build replies
	for _, reply := range comment.Replies {
		replyTree := r.buildTreeFromComment(reply)
		tree.Replies = append(tree.Replies, replyTree)
	}

	return tree
}
