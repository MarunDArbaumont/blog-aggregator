package main

import (
	"fmt"
	"time"
	"context"
	"github.com/google/uuid"

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

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("This command shouldn't have args")
	}

	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Printf("Feed: %+v\n", feed)

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("This command takes two args")
	}

	currentUser, err := s.db.GetUserByName(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	newFeedParams := database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.args[0],
		Url: cmd.args[1],
		UserID: currentUser.ID,
	}

	addedFeed, err := s.db.CreateFeed(context.Background(), newFeedParams)
	if err != nil {
		return err
	}

	fmt.Printf("Feed %+v has been added to database\n", addedFeed)

	newFeedFollowParams := database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: currentUser.ID,
		FeedID: addedFeed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), newFeedFollowParams)
	if err != nil {
		return err
	}

	fmt.Printf("You know follow the created feed: %v\n", addedFeed.Name)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("This command shouldn't have args")
	}

	listFeeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range listFeeds {
		feedUser, err := s.db.GetUserId(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("* %v: %v (user: %v)\n", feed.Name, feed.Url, feedUser)
	}
	return nil
}

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("This command shouldn have an arg")
	}

	wantedFeed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	currentUser, err := s.db.GetUserByName(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	newFeedFollowParams := database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: currentUser.ID,
		FeedID: wantedFeed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), newFeedFollowParams)
	if err != nil {
		return err
	}

	fmt.Printf("%v: %v\n", feedFollow.FeedName, feedFollow.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("This command shouldn't have args")
	}

	feedFollowsForUser, err := s.db.GetFeedFollowsForUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	for _, feedFollow := range feedFollowsForUser {
		fmt.Printf("* %v: %v\n", feedFollow.UserName, feedFollow.FeedName)
	}

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