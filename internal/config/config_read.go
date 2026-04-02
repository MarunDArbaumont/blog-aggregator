package config

import (
	"fmt"
	"os"
)

func Read(filepath string) (Config, error) {
	homeDir := os.UserHomeDir
	// file, err := os.Open(homeDir + filepath)
	fmt.Printf("This is your home dir: %v", homeDir)
	return Config{}, nil
}