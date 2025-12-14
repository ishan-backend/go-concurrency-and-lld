# Stack Overflow API

A simplified version of Stack Overflow API implementation in Go with 4 core endpoints:
1. Create Post
2. Add Comment / Reply Comment
3. Upvote Comment
4. Get All Comments Recursively

## Architecture

The project follows a clean architecture pattern with clear separation of concerns:

```
stackoverflow/
├── pkg/
│   ├── models/      # Domain models (Post, Comment)
│   ├── repository/  # Data access layer (in-memory storage)
│   ├── service/     # Business logic layer
│   └── handler/     # HTTP handlers (API endpoints)
├── main.go          # Server setup and routing
├── go.mod           # Go module definition
└── README.md        # This file
```

## API Endpoints

### 1. Create Post
**POST** `/api/posts`

Creates a new question/post.

**Request Body:**
```json
{
  "title": "How to implement a REST API in Go?",
  "content": "I'm new to Go and want to learn how to build REST APIs...",
  "author_id": "user-123"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "post-uuid",
    "title": "How to implement a REST API in Go?",
    "content": "I'm new to Go and want to learn how to build REST APIs...",
    "author_id": "user-123",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### 2. Add Comment / Reply Comment
**POST** `/api/comments`

Adds a comment to a post. If `parent_id` is provided, it creates a reply to that comment.

**Request Body (Top-level comment):**
```json
{
  "post_id": "post-uuid",
  "author_id": "user-456",
  "content": "Great question! Here's how you can do it..."
}
```

**Request Body (Reply to comment):**
```json
{
  "post_id": "post-uuid",
  "author_id": "user-789",
  "content": "I agree with your answer!",
  "parent_id": "comment-uuid"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "comment-uuid",
    "post_id": "post-uuid",
    "author_id": "user-456",
    "content": "Great question! Here's how you can do it...",
    "parent_id": null,
    "upvotes": 0,
    "created_at": "2024-01-15T10:35:00Z",
    "updated_at": "2024-01-15T10:35:00Z"
  }
}
```

### 3. Upvote Comment
**POST** `/api/comments/{id}/upvote`

Increments the upvote count for a comment.

**Request:** No body required

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "comment-uuid",
    "post_id": "post-uuid",
    "author_id": "user-456",
    "content": "Great question! Here's how you can do it...",
    "parent_id": null,
    "upvotes": 1,
    "created_at": "2024-01-15T10:35:00Z",
    "updated_at": "2024-01-15T10:35:00Z"
  }
}
```

### 4. Get All Comments Recursively
**GET** `/api/posts/{id}/comments`

Returns all comments and replies for a post in a flattened list. Comments are ordered with parent comments before their replies, recursively traversing the entire comment tree.

**Request:** No body required

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "comment-1-uuid",
      "post_id": "post-uuid",
      "author_id": "user-456",
      "content": "Great question! Here's how you can do it...",
      "parent_id": null,
      "upvotes": 5,
      "created_at": "2024-01-15T10:35:00Z",
      "updated_at": "2024-01-15T10:35:00Z"
    },
    {
      "id": "comment-2-uuid",
      "post_id": "post-uuid",
      "author_id": "user-789",
      "content": "I agree with your answer!",
      "parent_id": "comment-1-uuid",
      "upvotes": 2,
      "created_at": "2024-01-15T10:40:00Z",
      "updated_at": "2024-01-15T10:40:00Z"
    },
    {
      "id": "comment-3-uuid",
      "post_id": "post-uuid",
      "author_id": "user-101",
      "content": "Thanks for the clarification!",
      "parent_id": "comment-2-uuid",
      "upvotes": 0,
      "created_at": "2024-01-15T10:45:00Z",
      "updated_at": "2024-01-15T10:45:00Z"
    }
  ]
}
```

**Note:** The response includes all comments and nested replies in a flattened array, ordered such that parent comments appear before their replies. This allows you to easily iterate through all comments while maintaining the hierarchical relationship through the `parent_id` field.

## Error Responses

All endpoints return errors in the following format:

```json
{
  "success": false,
  "error": "Error message here"
}
```

Common error scenarios:
- Missing required fields
- Post not found
- Comment not found
- Parent comment not found
- Parent comment doesn't belong to the specified post

## Running the Server

1. Navigate to the project directory:
```bash
cd go/stackoverflow
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the server:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## Example Usage

### Create a Post
```bash
curl -X POST http://localhost:8080/api/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "How to use goroutines?",
    "content": "Can someone explain goroutines?",
    "author_id": "user-1"
  }'
```

### Add a Comment
```bash
curl -X POST http://localhost:8080/api/comments \
  -H "Content-Type: application/json" \
  -d '{
    "post_id": "post-uuid-from-previous-response",
    "author_id": "user-2",
    "content": "Goroutines are lightweight threads managed by Go runtime."
  }'
```

### Reply to a Comment
```bash
curl -X POST http://localhost:8080/api/comments \
  -H "Content-Type: application/json" \
  -d '{
    "post_id": "post-uuid",
    "author_id": "user-3",
    "content": "Thanks for the explanation!",
    "parent_id": "comment-uuid-from-previous-response"
  }'
```

### Upvote a Comment
```bash
curl -X POST http://localhost:8080/api/comments/{comment-id}/upvote
```

### Get All Comments Recursively
```bash
curl -X GET http://localhost:8080/api/posts/{post-id}/comments
```

## Design Patterns

1. **Repository Pattern**: Abstracts data access layer
2. **Service Layer Pattern**: Encapsulates business logic
3. **Handler Pattern**: Handles HTTP requests and responses

## Features

- ✅ Thread-safe in-memory storage using `sync.RWMutex`
- ✅ Support for nested comments (replies)
- ✅ Upvote functionality
- ✅ Input validation
- ✅ Error handling
- ✅ CORS support
- ✅ RESTful API design

## Future Enhancements

- Database persistence
- User authentication
- Downvote functionality
- Comment editing/deletion
- Pagination for comments
- Sorting comments by upvotes
- Rate limiting

