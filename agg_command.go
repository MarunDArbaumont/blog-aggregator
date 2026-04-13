package main

import (
	"fmt"
	"time"
	"context"
	"database/sql"
	"strconv"
	"log"

	"github.com/MarunDArbaumont/blog-aggregator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
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

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	var err error
	if len(cmd.args) > 0 {
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("error with %v: %v", limit, err)
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: int32(limit),
	})
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	for _, post := range posts {
		fmt.Println("__________________________________")
		fmt.Printf("Title: %v\n", post.Title)
		fmt.Printf("Url: %v\n", post.Url)
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

	layouts := []string{
		"Mon, 02 Jan 2006 15:04:05 -0700",
    	"2006-01-02T15:04:05Z",
		"01/02 03:04:05PM '06 -0700",
	}

	var pubDate time.Time

	for _, item := range feed.Channel.Item {
		for _, layout := range layouts {
			pubDate, err = time.Parse(layout, item.PubDate)
			if err == nil {
				break
			}
		}
		post, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Url: item.Link,
			Title: item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid: true,
			},
			PublishedAt: sql.NullTime{
				Time: pubDate,
				Valid: true,
			},
			FeedID: nextFeed.ID,
		})
		if err != nil {
			if pgError, ok := err.(*pq.Error); ok {
				if pgError.Code != "23505" {
					log.Printf("errror: %v", err)
				}
			} else {
				log.Printf("non-pg error: %v", err)
			}
		}
		fmt.Printf("%v: has been added to the database", post.Title)
	}

	return nil
}