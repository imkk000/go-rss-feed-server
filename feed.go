package main

import (
	"time"

	"github.com/gorilla/feeds"
)

func convertToFeeds(m map[string]any) (feed *feeds.Feed, err error) {
	feed = &feeds.Feed{
		Title:   getValOrDefault[string](m["title"]),
		Link:    &feeds.Link{Href: getValOrDefault[string](m["link"])},
		Author:  &feeds.Author{Name: getValOrDefault[string](m["author"])},
		Created: getValOrDefault(m["created"], time.Now()),
		Updated: getValOrDefault(m["updated"], time.Now()),
	}
	items := getValOrDefault[[]any](m["items"])
	feed.Items = make([]*feeds.Item, len(items))
	for i, item := range items {
		m := getValOrDefault[map[string]any](item)
		feed.Items[i] = &feeds.Item{
			Title:   getValOrDefault[string](m["title"]),
			Link:    &feeds.Link{Href: getValOrDefault[string](m["link"])},
			Author:  &feeds.Author{Name: getValOrDefault[string](m["author"])},
			Created: getValOrDefault(m["created"], time.Now()),
			Updated: getValOrDefault(m["updated"], time.Now()),
		}
	}

	return feed, nil
}

func getValOrDefault[T any](val any, def ...T) (r T) {
	v, ok := val.(T)
	if !ok {
		if def != nil {
			r = def[0]
		}
		return r
	}
	return v
}
