package posts

import (
	"sort"
	"time"

	"github.com/mmcdole/gofeed"
)

type Post struct {
	Title   string
	Link    string
	PubDate time.Time
}

func GetPosts() ([]Post, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://xqsit94.in/rss.xml")
	if err != nil {
		return nil, err
	}

	var posts []Post
	for _, item := range feed.Items {
		posts = append(posts, Post{
			Title:   item.Title,
			Link:    item.Link,
			PubDate: *item.PublishedParsed,
		})
	}

	// Sort posts by publication date (newest first)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].PubDate.After(posts[j].PubDate)
	})

	// Return only the first 5 posts
	if len(posts) > 5 {
		posts = posts[:5]
	}

	return posts, nil
}
