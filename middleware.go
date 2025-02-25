package main

import (
	"context"
	"fmt"
	"github.com/KrishKoria/Gator/internal/database"
)
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
    return func(s *state, cmd command) error {
        currentUser := s.Config.CurrentUserName
        user, err := s.DBQueries.GetUser(context.Background(), currentUser)
        if err != nil {
            return fmt.Errorf("error getting current user: %v", err)
        }
        return handler(s, cmd, user)
    }
}