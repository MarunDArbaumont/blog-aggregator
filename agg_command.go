package main

import (
	"fmt"
	"time"
	"context"
	"database/sql"

	"github.com/MarunDArbaumont/blog-aggregator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("Give a time duration")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("Error while parsing duration: %v\n", err)
	}

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

	return nil
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	argsMarkedFeed := database.MarkFeedFetchedParams{
		UpdatedAt: time.Now(),
		LastFetchedAt: sql.NullTime{
			Time: time.Now(),
			Valid: true,
		},
		ID: nextFeed.ID,
	}
	err = s.db.MarkFeedFetched(context.Background(), argsMarkedFeed)
	if err != nil {
		return err
	}

	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("Error while fetching data: %v\n", err)
	}

	for _, item := range feed.Channel.Item {
		fmt.Printf("%v\n", item.Title)
	}

	return nil
}