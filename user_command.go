package main

import (
	"fmt"
	"time"
	"context"
	
	"github.com/google/uuid"
	"github.com/MarunDArbaumont/blog-aggregator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("No username given")
	}

	_, err := s.db.GetUserByName(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User set to %v", s.cfg.CurrentUserName)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("No username given")
	}

	newUserParams := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.args[0],
	}
	_, err := s.db.CreateUser(context.Background(), newUserParams)
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User %v has been created. You are logged in as %v\n", cmd.args[0], s.cfg.CurrentUserName)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("This command shouldn't have args")
	}

	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("All users have been deleted from the users table")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("This command shouldn't have args")
	}

	listUsers, err := s.db.ListUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range listUsers {
		if user == s.cfg.CurrentUserName {
			fmt.Printf("* %v (current)\n", user)
			continue
		}
		fmt.Printf("* %v\n", user)
	}
	return nil
}