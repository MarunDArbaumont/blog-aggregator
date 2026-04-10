package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/MarunDArbaumont/blog-aggregator/internal/database"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("This command takes two args")
	}

	newFeedParams := database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.args[0],
		Url: cmd.args[1],
		UserID: user.ID,
	}

	addedFeed, err := s.db.CreateFeed(context.Background(), newFeedParams)
	if err != nil {
		return err
	}

	newFeedFollowParams := database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
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

	if len(listFeeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
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

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("This command shouldn have an arg")
	}

	wantedFeed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	newFeedFollowParams := database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: wantedFeed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), newFeedFollowParams)
	if err != nil {
		return err
	}

	fmt.Printf("%v: %v\n", feedFollow.FeedName, feedFollow.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("This command shouldn't have args")
	}

	feedFollowsForUser, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return err
	}

	for _, feedFollow := range feedFollowsForUser {
		fmt.Printf("* %v: %v\n", feedFollow.UserName, feedFollow.FeedName)
	}

	return nil
}

func handlerUnFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("This command shouldn't have args")
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	deleteFeedFollowArgs := database.DeleteFeedFollowParams{
		FeedID: feed.ID,
		UserID: user.ID,
	}

	err = s.db.DeleteFeedFollow(context.Background(), deleteFeedFollowArgs)
	if err != nil {
		return err
	}
	fmt.Printf("You have unfollowed %v\n", feed.Name)

	return nil
}
