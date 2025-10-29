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

// list all users and mark the current one:
// a handler function that uses shared state and the given command, does some work
//  and reports success/failure via an error
func handlerListUsers(s *state, cmd command) error {
	// Call the database query GetUsers to fetch all users:
	// (context.Background() supplies a base context to the DB call)
	/* A base context is a root context you start from when you don’t have an existing one to 
	derive. In Go, context.Background() returns an empty, never-canceled context with no values 
	or deadline. It’s commonly used at the top level (e.g., main, init, tests) to kick off 
	operations, and then you derive child contexts with timeouts, cancellation, or values using 
	context.WithCancel, context.WithTimeout, or context.WithValue. */
	users, err := s.db.GetUsers(context.Background())
	// If the query fails, return a wrapped error:
	// (%w wraps the original error so callers can unwrap it)
	if err != nil {
		return fmt.Errorf("couldn't list users: %w", err)
	}
	// Iterates over the users:
	for _, user := range users {
		// If a user’s name matches the current user from config (s.cfg.CurrentUserName), 
		// prints “* name (current)”:
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %v (current)\n", user.Name)
			// skips the non-current print for that user after printing the “(current)” line:
			continue
		}
		// Otherwise prints “* name”:
		fmt.Printf("* %v\n", user.Name)
	}
	// Returns nil on success:
	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}