package main

import (
    "context"
    "database/sql"
    "encoding/xml"
    "fmt"
    "html"
    "io"
    "net/http"
    "time"
    "github.com/google/uuid"
    "github.com/KrishKoria/Gator/internal/database"
    "github.com/lib/pq"
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
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

func parsePubDate(pubDate string) (time.Time, error) {
    layouts := []string{
        time.RFC1123Z,
        time.RFC1123,
        time.RFC3339,
        "Mon, 02 Jan 2006 15:04:05 MST",
    }

    var t time.Time
    var err error
    for _, layout := range layouts {
        t, err = time.Parse(layout, pubDate)
        if err == nil {
            return t, nil
        }
    }
    return t, fmt.Errorf("could not parse pubDate: %v", pubDate)
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
            pubDate, err := parsePubDate(item.PubDate)
            if err != nil {
                fmt.Printf("Error parsing pubDate: %v\n", err)
                continue
            }

            postID := uuid.New()
            now := time.Now()
            _, err = s.DBQueries.CreatePost(ctx, database.CreatePostParams{
                ID:          postID,
                CreatedAt:   now,
                UpdatedAt:   now,
                Title:       item.Title,
                Url:         item.Link,
                Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
                PublishedAt: sql.NullTime{Time: pubDate, Valid: !pubDate.IsZero()},
                FeedID:      feed.ID,
            })
            if err != nil {
                if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
                    continue
                }
                fmt.Printf("Error creating post: %v\n", err)
            }
        }
    }

    return nil
}