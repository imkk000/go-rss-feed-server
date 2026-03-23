package main

import (
	"testing"
	"time"

	"github.com/gorilla/feeds"
	"github.com/stretchr/testify/assert"
)

func TestConvertFeed(t *testing.T) {
	want := &feeds.Feed{
		Title:       "My Feed",
		Link:        &feeds.Link{Href: "http://myfeed.com"},
		Author:      &feeds.Author{Name: "kk"},
		Description: "text description",
		Created:     time.Now(),
		Updated:     time.Now(),
		Image: &feeds.Image{
			Title: "image title",
			Url:   "http://myfeed.com/image",
		},
		Items: []*feeds.Item{{
			Id:          "tag:item_1",
			Title:       "Item 1",
			Description: "Item details",
			Content:     "no content",
			Link:        &feeds.Link{Href: "http://myfeed.com/1"},
			Author:      &feeds.Author{Name: "kk"},
			Created:     time.Now(),
			Updated:     time.Now(),
			Enclosure: &feeds.Enclosure{
				Url: "http://myfeed.com/1/image",
			},
		}},
	}

	actual, err := convertToFeeds(map[string]any{
		"title":       want.Title,
		"link":        want.Link.Href,
		"author":      want.Author.Name,
		"description": want.Description,
		"created":     want.Created,
		"updated":     want.Updated,
		"image": M{
			"title": want.Image.Title,
			"url":   want.Image.Url,
		},
		"items": []any{
			map[string]any{
				"id":          want.Items[0].Id,
				"title":       want.Items[0].Title,
				"description": want.Items[0].Description,
				"content":     want.Items[0].Content,
				"link":        want.Items[0].Link.Href,
				"author":      want.Items[0].Author.Name,
				"created":     want.Items[0].Created,
				"updated":     want.Items[0].Updated,
				"image": M{
					"url": want.Items[0].Enclosure.Url,
				},
			},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, want, actual)
}
