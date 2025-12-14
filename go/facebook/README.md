# Facebook - Low Level Design

## Problem Statement

Design and implement a simplified Facebook-like social media system that supports core functionalities including user management, follower relationships, post creation/deletion, and personalized news feed generation. The system should be scalable, maintainable, and follow SOLID principles.

### Core Requirements

#### 1. User Management
- **Create User**: Users should be able to register with basic information:
  - First Name
  - Last Name
  - Date of Birth
  - Optional: Address (zip code, street name, city, country)
  - Optional: Gender (FEMALE, MALE, TRANSGENDER)
  - Optional: Phone Number
  - Optional: Email ID
- **Get User**: Retrieve user information by user ID
- Each user should have a unique auto-generated ID

#### 2. Follower System
- **Add Follower**: User A can follow User B (one-way relationship)
  - When User A follows User B, User B gains User A as a follower
  - A user can have multiple followers
- **Remove Follower**: User B can remove User A from their followers list
- **View Followers**: Users can view their list of followers

#### 3. Post Management
- **Create Post**: Users can create posts with unique post IDs
  - Posts should be stored in chronological order (most recent first)
  - Each post should have a unique auto-generated key for sorting
- **Delete Post**: Users can delete their own posts
- **View Posts**: Users can view all their posts
- Posts should be efficiently retrievable both by ID and in chronological order

#### 4. News Feed Generation
- **Get News Feed**: Generate a personalized news feed for a user containing:
  - Posts from all users that the current user follows
  - Posts sorted by recency (most recent first)
  - Support for limiting the number of posts returned (top N)
- **Paginated Feed**: Support pagination for the news feed
  - Default page size: 3 posts per page
  - Return posts for a specific page number
  - Handle edge cases (empty feed, invalid page numbers)

### Technical Constraints

1. **Data Structures**:
   - Use efficient data structures for fast lookups and chronological ordering
   - Posts should be stored in a way that allows O(1) insertion and deletion
   - Feed generation should efficiently merge posts from multiple followers

2. **Thread Safety** (if applicable):
   - Consider concurrent access scenarios
   - Ensure data consistency during concurrent operations

3. **Memory Efficiency**:
   - Optimize storage for large numbers of users and posts
   - Consider memory-efficient data structures

### Expected Functionalities

Based on the reference implementation, the system should support:

```
1. Create multiple users
2. Add followers (one-way following)
3. Remove followers
4. Create posts for users
5. Delete posts
6. Get news feed (top N posts from all followed users)
7. Get paginated news feed (page-based pagination)
```

### Example Usage Flow

```
1. Create users: User1, User2, User3, User4
2. User1 follows User2, User3
3. User2 follows User1, User4
4. User1 creates posts: Post1, Post2, Post3
5. User2 creates posts: Post4, Post5
6. User3 creates posts: Post6
7. Get news feed for User1 → Should show posts from User2 and User3 (Post6, Post5, Post4, ...)
8. Get paginated feed for User1 (page 1) → Should show first 3 posts
9. User1 deletes Post2
10. Get updated feed for User1
```

---

## SOLID Principles Application

The implementation should strictly adhere to SOLID principles:

### 1. Single Responsibility Principle (SRP)
- **UserService**: Responsible only for user-related operations (create, get, follower management)
- **PostService**: Responsible only for post-related operations (create, delete, retrieve posts)
- **FeedService**: Responsible only for feed generation and pagination logic
- **UserRepository**: Responsible only for user data persistence and retrieval
- Each model/struct should represent a single entity with clear responsibilities

### 2. Open/Closed Principle (OCP)
- Define interfaces for services and repositories
- Allow extension of functionality (e.g., different feed algorithms, different storage backends) without modifying existing code
- Use strategy pattern for different feed generation algorithms if needed
- Use dependency injection to allow swapping implementations

### 3. Liskov Substitution Principle (LSP)
- Any implementation of an interface should be substitutable without breaking the system
- Repository implementations should be interchangeable
- Service implementations should follow the same contract

### 4. Interface Segregation Principle (ISP)
- Create small, focused interfaces instead of large monolithic ones
- Separate interfaces for:
  - `UserRepository` (user data operations)
  - `PostRepository` (post data operations) - if needed
  - `FeedGenerator` (feed generation logic)
  - `FollowerManager` (follower relationship management)
- Clients should not be forced to depend on methods they don't use

### 5. Dependency Inversion Principle (DIP)
- High-level modules (services) should depend on abstractions (interfaces), not concrete implementations
- Services should depend on repository interfaces, not concrete repository implementations
- Use dependency injection to provide dependencies
- Example:
  ```go
  type UserService struct {
      userRepo UserRepository  // interface, not concrete type
  }
  ```

---

## Design Considerations

### Data Structures

1. **User Storage**:
   - Use a map/dictionary for O(1) user lookup by ID
   - Each user maintains a set of follower IDs

2. **Post Storage**:
   - Use a doubly-linked list per user for chronological ordering (O(1) insertion at head)
   - Use a map per user for O(1) post lookup by ID
   - Each post should have a unique key for sorting/prioritization

3. **Feed Generation**:
   - Use a priority queue (max-heap) to merge posts from multiple followers
   - Compare posts by their key (most recent first)
   - Efficiently handle merging of sorted lists

### Architecture Layers

```
┌─────────────────────────────────────┐
│         Service Layer               │
│  (UserService, PostService,         │
│   FeedService)                      │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│      Repository Layer                │
│  (UserRepository, PostRepository)   │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│         Model Layer                  │
│  (User, Post, Address, Gender)      │
└─────────────────────────────────────┘
```

### Key Design Patterns

1. **Repository Pattern**: Abstract data access layer
2. **Service Pattern**: Business logic encapsulation
3. **Singleton Pattern**: For service instances (if needed, but prefer dependency injection)
4. **Strategy Pattern**: For different feed generation algorithms (optional)
5. **Factory Pattern**: For creating service instances with proper dependencies

---

## Implementation Guidelines

### File Structure

```
facebook/
├── README.md
├── go.mod
├── models/
│   ├── user.go
│   ├── post.go
│   ├── address.go
│   └── gender.go
├── repository/
│   ├── user_repository.go
│   └── interfaces.go
├── service/
│   ├── user_service.go
│   ├── post_service.go
│   ├── feed_service.go
│   └── interfaces.go
└── main.go
```

### Interface Definitions

```go
// Repository Interfaces
type UserRepository interface {
    AddUser(user *User) error
    GetUser(id int) (*User, error)
}

// Service Interfaces
type UserService interface {
    CreateUser(firstName, lastName string, dob time.Time) (int, error)
    GetUser(userId int) (*User, error)
    AddFollower(followedBy, follows int) error
    RemoveFollower(userId, followerId int) error
}

type PostService interface {
    AddPost(userId, postId int) error
    DeletePost(userId, postId int) error
    GetPost(userId, postId int) (*Post, error)
}

type FeedService interface {
    GetNewsFeed(userId int, topN *int) ([]*Post, error)
    GetFeedPaginated(userId, pageNumber int) ([]*Post, error)
}
```

### Error Handling

- Use Go's error handling conventions
- Return meaningful error messages
- Handle edge cases (user not found, post not found, invalid page numbers, etc.)

### Testing

- Write unit tests for each service
- Test edge cases (empty followers, no posts, pagination boundaries)
- Test concurrent access scenarios if applicable

---

## Future Enhancements (Optional)

1. **Comments System**: Add ability to comment on posts
2. **Likes System**: Add ability to like posts and comments
3. **Notifications**: Notify users when someone follows them or likes their post
4. **Privacy Settings**: Control who can see posts
5. **Search Functionality**: Search users and posts
6. **Timeline**: Personal timeline vs. news feed
7. **Post Types**: Support different post types (text, image, video)
8. **Real-time Updates**: WebSocket support for real-time feed updates

---

## References

- Original implementation: [Dinkey11/MachineCoding/Facebook](https://github.com/Dinkey11/MachineCoding/tree/main/Facebook)
- SOLID Principles: [Wikipedia - SOLID](https://en.wikipedia.org/wiki/SOLID)

