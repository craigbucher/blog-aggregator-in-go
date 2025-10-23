package main

import (
	// "fmt"
	// "log"
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
