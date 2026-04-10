package main

import (
	"fmt"
	"context"

	"github.com/MarunDArbaumont/blog-aggregator/internal/config"
	"github.com/MarunDArbaumont/blog-aggregator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	command map[string]func(*state, command) error
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

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUserByName(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}