package main

import (
	"testing"
	"time"

	"github.com/gorilla/feeds"
	"github.com/stretchr/testify/assert"
)

func TestConvertFeed(t *testing.T) {
	want := &feeds.Feed{
		Title:   "My Feed",
		Link:    &feeds.Link{Href: "http://myfeed.com"},
		Author:  &feeds.Author{Name: "kk"},
		Created: time.Now(),
		Updated: time.Now(),
		Items: []*feeds.Item{{
			Title:   "Item 1",
			Link:    &feeds.Link{Href: "http://myfeed.com/1"},
			Author:  &feeds.Author{Name: "kk"},
			Created: time.Now(),
			Updated: time.Now(),
		}},
	}

	actual, err := convertToFeeds(map[string]any{
		"title":   want.Title,
		"link":    want.Link.Href,
		"author":  want.Author.Name,
		"created": want.Created,
		"updated": want.Updated,
		"items": []any{
			map[string]any{
				"title":   want.Items[0].Title,
				"link":    want.Items[0].Link.Href,
				"author":  want.Items[0].Author.Name,
				"created": want.Items[0].Created,
				"updated": want.Items[0].Updated,
			},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, want, actual)
}
