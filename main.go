package main

import (
	// "fmt"
	// "log"
	"context"
	"log"
	"os"
	"database/sql"

	//"github.com/bootdotdev/gator/internal/config"
	"gator/internal/config"
	"gator/internal/database"
	_ "github.com/lib/pq"
)

/* Before we can worry about command handlers, we need to think about how we will give our handlers 
access to the application state (later the database connection, but, for now, the config file). */
// Create a state struct that holds a pointer to a config:
type state struct {
	// Open a connection to the database, and store it in the state struct:
	db  *database.Queries
	cfg *config.Config
}

func main() {
	// Call config.Read() to load your ~/.gatorconfig.json into a Config struct:
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// In main(), load in your database URL to the config struct and sql.Open() a connection to your database:
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()
	// Use your generated database package to create a new *database.Queries, and store it in 
	// your state struct:
	dbQueries := database.New(db)

	// n the main function, remove the manual update of the config file. Instead, simply 
	// read the config file, and store the config in a new instance of the state struct:
	programState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}
	// Create a new instance of the commands struct with an initialized map of handler functions:
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	// Register a handler function for the login command:
	cmds.register("login", handlerLogin)
	// Create a register handler and register it with the commands:
	cmds.register("register", handlerRegister)
	// Add a new command called reset that calls the query:
	cmds.register("reset", handlerReset)
	// Add a new command called users that calls GetUsers and prints all the users to the console:
	cmds.register("users", handlerListUsers)
	// Add an agg command:
	cmds.register("agg", handlerAgg)
	// Add a new command called addfeed:
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	// Add a new feeds handler. It takes no arguments and prints all the feeds in the database 
	// to the console:
	cmds.register("feeds", handlerListFeeds)
	/* Add a follow command. It takes a single url argument and creates a new feed follow record 
	for the current user. It should print the name of the feed and the current user once the record 
	is created (which the query we just made should support). You'll need a query to look up feeds 
	by URL */
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	// Add a following command. It should print all the names of the feeds the current user is following:
	cmds.register("following", middlewareLoggedIn(handlerListFeedFollows))
	// Add a new unfollow command that accepts a feed's URL as an argument and unfollows it for 
	// the current user. This is, of course, a "logged in" command - use the new middleware:
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	// Add the browse command. It should take an optional "limit" parameter
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))
	/* If there are fewer than 2 arguments, print an error message to the terminal and exit. 
	Why two? The first argument is automatically the program name, which we ignore, and we 
	require a command name */
	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}
	// Use os.Args to get the command-line arguments passed in by the user:
	// You'll need to split the command-line arguments into the command name and the arguments 
	// slice to create a command instance:
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	//  Use the commands.run method to run the given command and print any errors returned:
	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}

/* Create logged-in middleware. It will allow us to change the function signature of our handlers 
that require a logged in user to accept a user as an argument and DRY up our code. Here's the 
function signature of my middleware: 
You'll notice it's a higher order function that takes a handler of the "logged in" type and returns a "normal" handler 
that we can register.*/
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}