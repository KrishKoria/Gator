package main

import (
    "context"
    "encoding/xml"
    "fmt"
    "html"
    "io"
    "net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error){
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
        return nil, fmt.Errorf("failed to create request: %v", err)
    }

	req.Header.Set("User-Agent", "gator")


	client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch feed: %v", err)
    }
    defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

	body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %v", err)
    }

	var feed RSSFeed
    if err := xml.Unmarshal(body, &feed); err != nil {
        return nil, fmt.Errorf("failed to unmarshal XML: %v", err)
    }

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
    feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
    for i := range feed.Channel.Item {
        feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
        feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
    }

    return &feed, nil
}

func scrapeFeeds(ctx context.Context, s *state) error {
    feed, err := s.DBQueries.GetNextFeedToFetch(ctx)
    if err != nil {
        return fmt.Errorf("error getting next feed to fetch: %v", err)
    }

    err = s.DBQueries.MarkFeedFetched(ctx, feed.ID)
    if err != nil {
        return fmt.Errorf("error marking feed as fetched: %v", err)
    }

    rssFeed, err := fetchFeed(ctx, feed.Url)
    if err != nil {
        return fmt.Errorf("error fetching feed: %v", err)
    }

    for _, item := range rssFeed.Channel.Item {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            fmt.Printf("Title: %s\n", item.Title)
        }
    }

    return nil
}