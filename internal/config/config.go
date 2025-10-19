package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)
// declares a Go package-level constant string named configFileName with value ".gatorconfig.json". 
// You use it anywhere you need the config’s filename (e.g., building ~/<name> paths) so it’s 
// centralized and not hard-coded in multiple places:
const configFileName = ".gatorconfig.json"
// define a struct type that mirrors your JSON file’s shape:
type Config struct {
	DBURL           string `json:"db_url"`				// DBURL holds the database URL, maps to JSON key db_url
	CurrentUserName string `json:"current_user_name"`	// holds the logged-in user, maps to JSON key current_user_name
}
// method on Config that updates and persists the current user:
func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName		//  (note the pointer receiver, so it mutates the original)
	return write(*cfg)		// Calls write(*cfg) to serialize the updated struct and save it to the config file
}
// read and decode the JSON config from disk into a Config:
func Read() (Config, error) {
	fullPath, err := getConfigFilePath()	// Get the file path
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(fullPath)	// open the file
	if err != nil {
		return Config{}, err
	}
	defer file.Close()				// defer close

	decoder := json.NewDecoder(file)	
	cfg := Config{}
	err = decoder.Decode(&cfg)		// // Decode JSON
	if err != nil {
		return Config{}, err
	}
	// On success, return the populated cfg:
	return cfg, nil
}
//  builds the absolute path to your config file:
func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()	// Get the user’s home directory
	if err != nil {
		return "", err
	}
	// Join it with configFileName (e.g., ".gatorconfig.json") using filepath.Join:
	fullPath := filepath.Join(home, configFileName)
	// Returns the full path or an error if HOME couldn’t be determined:
	return fullPath, nil
}
// write the Config to the JSON file:
func write(cfg Config) error {
	fullPath, err := getConfigFilePath()	// Build the file path
	if err != nil {
		return err
	}

	file, err := os.Create(fullPath)	// Creates/truncates the file
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)	// JSON-encodes cfg to that file
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}
	// Return any error encountered; otherwise nil:
	return nil
}
