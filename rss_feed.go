package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
	"time"
)

// We directly unmarshal the XML document into structs like this:
type RSSFeed struct {
	// RSSFeed has a single field, Channel:
	Channel struct {
		// Channel is an anonymous nested struct with fields Title, Link, Description, and Item:
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		// Item []RSSItem means multiple <item> elements become a slice of RSSItem:
		Item        []RSSItem `xml:"item"`
	// The tag on Channel, xml:"channel", tells the decoder to look for the <channel> element 
	// and fill that nested struct:
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// Write a func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) function. It 
// should fetch a feed from the given URL, and, assuming that nothing goes wrong, return a 
// filled-out RSSFeed struct:
func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// construct a new http.Client value using a composite literal:
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}
	// construct an HTTP request object:
		// Use ctx to control cancellation/timeouts
		// Method is "GET"
		// URL is feedURL
		// Body is nil (no request body)
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	// set the User-Agent header to gator in the request with request.Header.Set. 
	// This is a common practice to identify your program to the server:
	req.Header.Set("User-Agent", "gator")
	// send the HTTP request and return the server’s response:
		// httpClient.Do(req) performs the network call
		// resp is an *http.Response with status, headers, and Body (an io.ReadCloser)
		// err is non-nil if the request failed (DNS, connect, timeout, etc.)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	// You should defer resp.Body.Close() after checking err:
	defer resp.Body.Close()

	// read the entire HTTP response body into memory:
		// io.ReadAll consumes resp.Body (an io.Reader) and returns a byte slice dat
		// After this, you should close resp.Body (if you haven’t already deferred Close)
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// declare a new RSSFeed variable
	// When you declare a variable this way without explicitly initializing it, Go automatically 
	// sets it to the zero value for that type. For a struct like RSSFeed, the zero value means:
		// All string fields will be empty strings ("")
		// All int fields will be 0
		// All slice fields will be nil
		// All nested structs will also be zero-valued
	var rssFeed RSSFeed
	// after this declaration, you have an empty RSSFeed struct that you can then populate - 
	// for example, by unmarshaling XML data into it:
	err = xml.Unmarshal(dat, &rssFeed)	//  (works the same as json.Unmarshal)
	if err != nil {
		return nil, err
	}

	// Use the html.UnescapeString function to decode escaped HTML entities (like &ldquo;). 
	// You'll need to run the Title and Description fields (of both the entire channel as well 
	// as the items) through this function:
		// &ldquo; becomes "
		// &rdquo; becomes "
		// &amp; becomes &
		// &#39; becomes '
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i, item := range rssFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		rssFeed.Channel.Item[i] = item
	}

	return &rssFeed, nil
}
