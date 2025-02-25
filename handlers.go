package main

import (
	"context"
	"fmt"
	"os"
	"time"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
    "strconv"
    "log"
	"github.com/KrishKoria/Gator/internal/database"
)


func handlerLogin(s *state, cmd command) error {
    if len(cmd.Args) == 0 {
        return fmt.Errorf("enter a username")
    }
    username := cmd.Args[0]
    _, err := s.DBQueries.GetUser(context.Background(), username)
    if err != nil {
        fmt.Println("Error: User does not exist")
        os.Exit(1)
    }

    err = s.Config.SetUser(username)
    if err != nil {
        return fmt.Errorf("error setting user: %v", err)
    }
    fmt.Printf("User set to %s\n", username)
    return nil
}

func handlerRegister(s *state, cmd command) error {
    if len(cmd.Args) == 0 {
        return fmt.Errorf("enter a username")
    }
    username := cmd.Args[0]

    _, err := s.DBQueries.GetUser(context.Background(), username)
    if err == nil {
        fmt.Println("Error: User already exists")
        os.Exit(1)
    }

    userID := uuid.New()
    now := time.Now()
    user, err := s.DBQueries.CreateUser(context.Background(), database.CreateUserParams{
        ID:        userID,
        CreatedAt: now,
        UpdatedAt: now,
        Name:      username,
    })
    if err != nil {
        return fmt.Errorf("error creating user: %v", err)
    }

    err = s.Config.SetUser(username)
    if err != nil {
        return fmt.Errorf("error setting user: %v", err)
    }

    fmt.Printf("User created: %+v\n", user)
    return nil
}

func handlerUsers(s *state, cmd command) error {
    users, err := s.DBQueries.GetUsers(context.Background())
    if err != nil {
        return fmt.Errorf("error getting users: %v", err)
    }

    currentUser := s.Config.CurrentUserName
    for _, user := range users {
        if user.Name == currentUser {
            fmt.Printf("* %s (current)\n", user.Name)
        } else {
            fmt.Printf("* %s\n", user.Name)
        }
    }
    return nil
}

func handlerReset(s *state, cmd command) error {
    err := s.DBQueries.DeleteAllUsers(context.Background())
    if err != nil {
        fmt.Println("Error: Failed to reset the users table")
        os.Exit(1)
    }

    err = s.DBQueries.DeleteAllPosts(context.Background())
    if err != nil {
        fmt.Println("Error: Failed to reset the posts table")
        os.Exit(1)
    }

    err = s.DBQueries.DeleteAllFeeds(context.Background())
    if err != nil {
        fmt.Println("Error: Failed to reset the feeds table")
        os.Exit(1)
    }

    fmt.Println("All tables reset successfully")
    return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) < 1 || len(cmd.Args) > 2 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	log.Printf("Collecting feeds every %s...", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}


func handlerAddFeed(s *state, cmd command, user database.User) error {
    if len(cmd.Args) < 2 {
        return fmt.Errorf("enter a feed name and URL")
    }
    feedName := cmd.Args[0]
    feedURL := cmd.Args[1]

    currentUser := s.Config.CurrentUserName
    user, err := s.DBQueries.GetUser(context.Background(), currentUser)
    if err != nil {
        return fmt.Errorf("error getting current user: %v", err)
    }

    feedID := uuid.New()
    now := time.Now()
    feed, err := s.DBQueries.CreateFeed(context.Background(), database.CreateFeedParams{
        ID:        feedID,
        CreatedAt: now,
        UpdatedAt: now,
        Name:      feedName,
        Url:       feedURL,
        UserID:    user.ID,
    })
    if err != nil {
        return fmt.Errorf("error creating feed: %v", err)
    }
    
    fmt.Printf("Feed created: %+v\n", feed)

    followID := uuid.New()
    follow, err := s.DBQueries.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
        ID:        followID,
        CreatedAt: now,
        UpdatedAt: now,
        UserID:    user.ID,
        FeedID:    feedID,
    })
    if err != nil {
        return fmt.Errorf("error following feed: %v", err)
    }

    fmt.Printf("Now following '%s' as user '%s'\n", follow.FeedName, follow.UserName)
    
    return nil
}

func handlerFeeds(s *state, cmd command) error {
    feeds, err := s.DBQueries.GetFeedsWithUsers(context.Background())
    if err != nil {
        return fmt.Errorf("error getting feeds: %v", err)
    }

    for _, feed := range feeds {
        fmt.Printf("Feed Name: %s\n", feed.FeedName)
        fmt.Printf("Feed URL: %s\n", feed.Url)
        fmt.Printf("Created by: %s\n\n", feed.UserName)
    }
    return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
    if len(cmd.Args) < 1 {
        return fmt.Errorf("enter a feed URL to follow")
    }
    feedURL := cmd.Args[0]

    currentUser := s.Config.CurrentUserName
    user, err := s.DBQueries.GetUser(context.Background(), currentUser)
    if err != nil {
        return fmt.Errorf("error getting current user: %v", err)
    }

    feed, err := s.DBQueries.GetFeedByURL(context.Background(), feedURL)
    if err != nil {
        return fmt.Errorf("feed not found with URL %s: %v", feedURL, err)
    }

    followID := uuid.New()
    now := time.Now()
    follow, err := s.DBQueries.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
        ID:        followID,
        CreatedAt: now,
        UpdatedAt: now,
        UserID:    user.ID,
        FeedID:    feed.ID,
    })
    if err != nil {
        return fmt.Errorf("error following feed: %v", err)
    }

    fmt.Printf("Now following '%s' as user '%s'\n", follow.FeedName, follow.UserName)
    return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
    currentUser := s.Config.CurrentUserName
    user, err := s.DBQueries.GetUser(context.Background(), currentUser)
    if err != nil {
        return fmt.Errorf("error getting current user: %v", err)
    }

    feedFollows, err := s.DBQueries.GetFeedFollowsForUser(context.Background(), user.ID)
    if err != nil {
        return fmt.Errorf("error getting followed feeds: %v", err)
    }

    if len(feedFollows) == 0 {
        fmt.Printf("User '%s' is not following any feeds\n", currentUser)
        return nil
    }

    fmt.Printf("Feeds followed by '%s':\n", currentUser)
    for _, follow := range feedFollows {
        fmt.Printf("* %s\n", follow.FeedName)
    }
    return nil
}


func handlerUnfollow(s *state, cmd command, user database.User) error {
    if len(cmd.Args) < 1 {
        return fmt.Errorf("enter a feed URL to unfollow")
    }
    feedURL := cmd.Args[0]

    err := s.DBQueries.DeleteFeedFollowByUserAndFeedURL(context.Background(), database.DeleteFeedFollowByUserAndFeedURLParams{
        UserID: user.ID,
        Url:    feedURL,
    })
    if err != nil {
        return fmt.Errorf("error unfollowing feed with URL %s: %v", feedURL, err)
    }

    fmt.Printf("Unfollowed feed with URL '%s' for user '%s'\n", feedURL, user.Name)
    return nil
}


func handlerBrowse(s *state, cmd command, user database.User) error {
    limit := 2
    if len(cmd.Args) > 0 {
        var err error
        limit, err = strconv.Atoi(cmd.Args[0])
        if err != nil {
            return fmt.Errorf("invalid limit: %v", err)
        }
    }

    posts, err := s.DBQueries.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
        Name:  user.Name,
        Limit: int32(limit),
    })
    if err != nil {
        return fmt.Errorf("error getting posts for user: %v", err)
    }

    for _, post := range posts {
        fmt.Printf("Title: %s\n", post.Title)
        fmt.Printf("URL: %s\n", post.Url)
        fmt.Printf("Description: %s\n", post.Description.String)
        fmt.Printf("Published At: %s\n\n", post.PublishedAt.Time)
    }

    return nil
}