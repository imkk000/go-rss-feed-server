package main

import (
	"time"

	"github.com/gorilla/feeds"
)

func convertToFeeds(m map[string]any) (feed *feeds.Feed, err error) {
	feed = &feeds.Feed{
		Title:       getValOrDefault[string](m["title"]),
		Link:        &feeds.Link{Href: getValOrDefault[string](m["link"])},
		Description: getValOrDefault[string](m["description"]),
		Author:      &feeds.Author{Name: getValOrDefault[string](m["author"])},
		Created:     getValOrDefault(m["created"], time.Now()),
		Updated:     getValOrDefault(m["updated"], time.Now()),
	}
	if m["image"] != nil {
		image := getValOrDefault[M](m["image"])
		feed.Image = &feeds.Image{
			Title: getValOrDefault[string](image["title"]),
			Url:   getValOrDefault[string](image["url"]),
		}
	}
	items := getValOrDefault[[]any](m["items"])
	feed.Items = make([]*feeds.Item, len(items))
	for i, item := range items {
		m := getValOrDefault[map[string]any](item)
		feed.Items[i] = &feeds.Item{
			Id:          getValOrDefault[string](m["id"]),
			Title:       getValOrDefault[string](m["title"]),
			Description: getValOrDefault[string](m["description"]),
			Content:     getValOrDefault[string](m["content"]),
			Link:        &feeds.Link{Href: getValOrDefault[string](m["link"])},
			Author:      &feeds.Author{Name: getValOrDefault[string](m["author"])},
			Created:     getValOrDefault(m["created"], time.Now()),
			Updated:     getValOrDefault(m["updated"], time.Now()),
		}
		if m["image"] != nil {
			image := getValOrDefault[M](m["image"])
			feed.Items[i].Enclosure = &feeds.Enclosure{
				Url: getValOrDefault[string](image["url"]),
			}
		}
	}

	return feed, nil
}
