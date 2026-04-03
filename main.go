package main

import (
	"fmt"
	"log"
	"os"

	"github.com/MarunDArbaumont/blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	mainState := &state{
		cfgPointer: cfg,
	}	

	mainCommands := commands{
		command: make(map[string]func(*state, command) error),
	}

	mainCommands.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Errorf("not enough arguments")
		os.Exit(1)
	}

	commandName := os.Args[1]
	commandArgs := os.Args[2:]

	mainCommand := command{
		name: commandName,
		args: commandArgs,
	}

	err = mainCommands.run(mainState, mainCommand)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}