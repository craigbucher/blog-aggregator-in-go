package main

import (
	"context"
	"fmt"
)
// Add a new command called reset that calls the query. Report back to the user about whether 
// or not it was successful with an appropriate exit code:
func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't delete users: %w", err)
	}
	fmt.Println("Database reset successfully!")
	return nil
}
