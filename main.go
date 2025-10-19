package main

import (
	// "fmt"
	// "log"
	"log"
	"os"

	//"github.com/bootdotdev/gator/internal/config"
	"gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	// Call config.Read() to load your ~/.gatorconfig.json into a Config struct:
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	
	programState := &state{
		cfg: &cfg,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
