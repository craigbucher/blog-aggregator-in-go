package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"gator/internal/database"
	"github.com/google/uuid"
)
// Update the agg command to now take a single argument: time_between_reqs:
/* time_between_reqs is a duration string, like 1s, 1m, 1h, etc. I used the time.ParseDuration 
function to parse it into a time.Duration value */
func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) < 1 || len(cmd.Args) > 2 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}
	// It should print a message like Collecting feeds every 1m0s when it starts:
	log.Printf("Collecting feeds every %s...", timeBetweenRequests)
	// Use a time.Ticker to run your scrapeFeeds function once every time_between_reqs. I used a 
	// for loop to ensure that it runs immediately (I don't like waiting) and then every time the 
	// ticker ticks:
	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}
// Write an aggregation function, I called mine scrapeFeeds
func scrapeFeeds(s *state) {
	// Get the next feed to fetch from the DB:
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Couldn't get next feeds to fetch", err)
		return
	}
	log.Println("Found a feed to fetch!")
	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	// Mark the feed as fetched:
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}
	// Fetch the feed using the URL:
	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}
	// Update your scraper to save posts. Instead of printing out the 
	// titles of the posts, save them to the database!:
	for _, item := range feedData.Channel.Item {
		// Make sure that you're parsing the "published at" time properly 
		// from the feeds. Sometimes they might be in a different format than 
		// you expect, so you might need to handle that:
		// (You may have to manually convert the data into database/sql types)
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			FeedID:    feed.ID,
			Title:     item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			Url:         item.Link,
			PublishedAt: publishedAt,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}