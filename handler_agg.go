package main

import (
	"context"
	"fmt"
)

/* Add an agg command. Later this will be our long-running aggregator service. For now, we'll 
just use it to fetch a single feed and ensure our parsing works. It should fetch the feed found 
at https://www.wagslane.dev/index.xml and print the entire struct to the console */
func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}
	fmt.Printf("Feed: %+v\n", feed)
	return nil
}