package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

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
        fmt.Println("Error: Failed to reset the database")
        os.Exit(1)
    }

    fmt.Println("Database reset successfully")
    return nil
}

func handlerAgg(s *state, cmd command) error {
    ctx := context.Background()
    feedURL := "https://www.wagslane.dev/index.xml"
    feed, err := fetchFeed(ctx, feedURL)
    if err != nil {
        return fmt.Errorf("error fetching feed: %v", err)
    }

    fmt.Printf("Fetched Feed: %+v\n", feed)
    return nil
}


func handlerAddFeed(s *state, cmd command) error {
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
