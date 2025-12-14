package main

import (
	"fmt"
	"log"

	"github.com/ishanpandey/LowLevelDesign/go/facebook/repository"
	"github.com/ishanpandey/LowLevelDesign/go/facebook/service"
)

func main() {
	// Initialize repository (following DIP - high-level modules depend on abstractions)
	userRepo := repository.NewUserRepository()

	// Initialize services with dependency injection (following DIP)
	userService := service.NewUserService(userRepo)
	postService := service.NewPostService(userRepo)
	feedService := service.NewFeedService(userRepo)

	// Create users
	fmt.Println("=== Creating Users ===")
	user1ID, err := userService.CreateUser("Dinkey", "Mahawar", "1998-02-11")
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	fmt.Printf("Created user with ID: %d\n", user1ID)

	user2ID, err := userService.CreateUser("Rashi", "Mahawar", "1998-02-11")
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	fmt.Printf("Created user with ID: %d\n", user2ID)

	user3ID, err := userService.CreateUser("Shruti", "Mahawar", "1998-02-11")
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	fmt.Printf("Created user with ID: %d\n", user3ID)

	user4ID, err := userService.CreateUser("Shubh", "Mahawar", "1998-02-11")
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	fmt.Printf("Created user with ID: %d\n", user4ID)

	user5ID, err := userService.CreateUser("Bhavya", "Gupta", "1998-02-11")
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	fmt.Printf("Created user with ID: %d\n", user5ID)

	user6ID, err := userService.CreateUser("Knaha", "Goyal", "1998-02-11")
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	fmt.Printf("Created user with ID: %d\n", user6ID)

	user7ID, err := userService.CreateUser("Chaitnaya", "Sharma", "1998-02-11")
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	fmt.Printf("Created user with ID: %d\n", user7ID)

	user8ID, err := userService.CreateUser("Vishnu", "UK", "1998-02-11")
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	fmt.Printf("Created user with ID: %d\n", user8ID)

	user9ID, err := userService.CreateUser("Tanisha", "Yadav", "1998-02-11")
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	fmt.Printf("Created user with ID: %d\n", user9ID)

	// Add followers
	fmt.Println("\n=== Adding Followers ===")
	userService.AddFollower(user1ID, user2ID)
	userService.AddFollower(user1ID, user3ID)
	userService.AddFollower(user1ID, user8ID)
	userService.AddFollower(user2ID, user3ID)
	userService.AddFollower(user2ID, user4ID)
	userService.AddFollower(user2ID, user1ID)
	userService.AddFollower(user2ID, user5ID)
	userService.AddFollower(user4ID, user7ID)
	userService.AddFollower(user4ID, user6ID)
	userService.AddFollower(user1ID, user4ID)
	userService.AddFollower(user3ID, user7ID)
	userService.AddFollower(user5ID, user6ID)

	// Print followers
	fmt.Println("\n=== User Followers ===")
	user1, _ := userService.GetUser(user1ID)
	fmt.Printf("User %d followers: %v\n", user1ID, user1.GetFollowerIDs())

	// Remove follower
	fmt.Println("\n=== Removing Followers ===")
	userService.RemoveFollower(user1ID, user2ID)
	user1, _ = userService.GetUser(user1ID)
	fmt.Printf("User %d followers after removal: %v\n", user1ID, user1.GetFollowerIDs())

	user5, _ := userService.GetUser(user5ID)
	fmt.Printf("User %d followers: %v\n", user5ID, user5.GetFollowerIDs())
	userService.RemoveFollower(user5ID, user6ID)
	user5, _ = userService.GetUser(user5ID)
	fmt.Printf("User %d followers after removal: %v\n", user5ID, user5.GetFollowerIDs())

	user2, _ := userService.GetUser(user2ID)
	fmt.Printf("User %d followers: %v\n", user2ID, user2.GetFollowerIDs())
	userService.RemoveFollower(user2ID, user4ID)
	user2, _ = userService.GetUser(user2ID)
	fmt.Printf("User %d followers after removal: %v\n", user2ID, user2.GetFollowerIDs())

	// Add posts
	fmt.Println("\n=== Adding Posts ===")
	postService.AddPost(user1ID, 1)
	postService.AddPost(user1ID, 2)
	postService.AddPost(user1ID, 12)
	postService.AddPost(user1ID, 13)
	postService.AddPost(user1ID, 16)
	postService.AddPost(user1ID, 17)
	postService.AddPost(user2ID, 5)
	postService.AddPost(user2ID, 1)
	postService.AddPost(user2ID, 4)
	postService.AddPost(user2ID, 3)
	postService.AddPost(user2ID, 8)
	postService.AddPost(user2ID, 7)
	postService.AddPost(user3ID, 9)
	postService.AddPost(user4ID, 10)
	postService.AddPost(user1ID, 14)
	postService.AddPost(user1ID, 15)
	postService.AddPost(user4ID, 11)
	postService.AddPost(user1ID, 3)
	postService.AddPost(user1ID, 4)

	// Print posts
	fmt.Println("\n=== User Posts ===")
	user2, _ = userService.GetUser(user2ID)
	fmt.Printf("User %d posts: %v\n", user2ID, user2.GetPostIDs())

	// Delete post
	fmt.Println("\n=== Deleting Posts ===")
	postService.DeletePost(user2ID, 7)
	user2, _ = userService.GetUser(user2ID)
	fmt.Printf("User %d posts after deletion: %v\n", user2ID, user2.GetPostIDs())

	// Get news feed
	fmt.Println("\n=== News Feed ===")
	topN := 4
	feed, err := feedService.GetNewsFeed(user2ID, &topN)
	if err != nil {
		log.Fatalf("Failed to get news feed: %v", err)
	}
	service.PrintFeed(feed, user2ID)

	// Get paginated feed
	fmt.Println("\n=== Paginated Feed ===")
	feed, err = feedService.GetFeedPaginated(user2ID, 2)
	if err != nil {
		log.Fatalf("Failed to get paginated feed: %v", err)
	}
	service.PrintFeed(feed, user2ID)

	feed, err = feedService.GetFeedPaginated(user2ID, 1)
	if err != nil {
		log.Fatalf("Failed to get paginated feed: %v", err)
	}
	service.PrintFeed(feed, user2ID)

	feed, err = feedService.GetFeedPaginated(user2ID, 4)
	if err != nil {
		log.Fatalf("Failed to get paginated feed: %v", err)
	}
	service.PrintFeed(feed, user2ID)
}
