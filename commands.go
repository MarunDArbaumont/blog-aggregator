package main

import (
	"fmt"
	"github.com/MarunDArbaumont/blog-aggregator/internal/config"
)

type state struct {
	cfgPointer config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	command map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("No username given")
	}

	err := s.cfgPointer.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User set to %v", s.cfgPointer.CurrentUserName)
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	commandCallback, exists := c.command[cmd.name]
	if exists {
		return commandCallback(s, cmd)
	}
	return fmt.Errorf("%v is not a command\n",cmd.name)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.command[name] = f
}