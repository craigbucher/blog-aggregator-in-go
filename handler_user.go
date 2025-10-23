package main

import (
	"context"
	"fmt"
	"time"

	"gator/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <name>", cmd.Name)
	}

	name := cmd.Args[0]
	// Pass context.Background() to the query to create an empty Context argument:
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		// Use the uuid.New() function to generate a new UUID for the user:
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),	// created_at and updated_at should be the current time
		UpdatedAt: time.Now().UTC(),	// created_at and updated_at should be the current time
		Name:      name,	// Use the provided name
	})
	// Exit with code 1 if a user with that name already exists:
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	// Print a message that the user was created, and log the user's data to the console 
	// for your own debugging:
	fmt.Println("User created successfully:")
	printUser(user)
	return nil
}

// Create a login handler function: func handlerLogin(s *state, cmd command) error. 
// This will be the function signature of all command handlers:
func handlerLogin(s *state, cmd command) error {
	// If the command's arg's slice is empty, return an error; the login handler expects 
	// a single argument, the username:
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]
	// Update the login command handler to error (and exit with code 1) if the given 
	// username doesn't exist in the database:
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}

	// Use the state's access to the config struct to set the user to the given username. 
	// Remember to return any errors:
	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	// Print a message to the terminal that the user has been set:
	fmt.Println("User switched successfully!")
	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}