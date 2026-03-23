package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/grafana/sobek"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
)

func TestParseHTMLFn(t *testing.T) {
	wantText := "Span Text"
	wantAttr := "color:blue"
	wantLength := int64(3)
	wantIndexes := []int64{0, 1, 2}

	vm := sobek.New()
	err := SetParseHTMLFn(vm)
	assert.NoError(t, err)

	fn, ok := sobek.AssertFunction(vm.Get("parseHTML"))
	assert.True(t, ok)

	html := `<div class="container" style="color:blue">
	<table name="test">
	<th><td>Name</td><td>Value</td></th>
	<tr><td>Item1</td><td>Value1</td></tr>
	<tr><td>Item2</td><td>Value2</td></tr>
	</table>
	<span>Span Text</span>
	</div>`
	result, err := fn(sobek.Undefined(), vm.ToValue(html))
	m := result.Export().(M)
	findFn := m["find"].(func(string) sobek.Value)
	divContainer := findFn(".container").ToObject(vm).Export().(M)
	attr := divContainer["attr"].(func(string) string)
	innerFindFn := divContainer["find"].(func(string) sobek.Value)
	span := innerFindFn("span").ToObject(vm).Export().(M)
	text := span["text"].(func() string)
	tableTr := innerFindFn("table tr").ToObject(vm).Export().(M)
	length := tableTr["length"].(int64)
	tableTrEach := tableTr["each"].(func(sobek.Callable))
	var indexes []int64
	tableTrEach(func(this sobek.Value, args ...sobek.Value) (sobek.Value, error) {
		indexes = append(indexes, args[0].ToInteger())
		return nil, nil
	})

	assert.Equal(t, wantLength, length)
	assert.Equal(t, wantIndexes, indexes)
	assert.Equal(t, wantAttr, attr("style"))
	assert.Equal(t, wantText, text())
}

func TestConvertMapFeedFn(t *testing.T) {
	wantItemImage := M{
		"title": "my image 1",
		"url":   "http://myfeed.com/1/image",
	}
	wantItem := M{
		"id":          "guid",
		"title":       "My Feed 1",
		"description": "my details 1",
		"content":     "full content 1",
		"link":        "http://myfeed.com/1",
		"author":      "kk",
		"created":     time.Now(),
		"updated":     time.Now(),
		"image":       wantItemImage,
	}
	wantResult := M{
		"title":       "My Feed",
		"link":        "http://myfeed.com",
		"author":      "kk",
		"description": "my details",
		"created":     time.Now(),
		"updated":     time.Now(),
		"items":       []any{wantItem},
	}

	vm := sobek.New()
	err := SetConvertMapFeedFn(vm)
	assert.NoError(t, err)

	fn, ok := sobek.AssertFunction(vm.Get("convertMapFeed"))
	assert.True(t, ok)

	feed := &gofeed.Feed{
		Title:           wantResult["title"].(string),
		Link:            wantResult["link"].(string),
		Authors:         []*gofeed.Person{{Name: wantResult["author"].(string)}},
		PublishedParsed: new(wantResult["created"].(time.Time)),
		UpdatedParsed:   new(wantResult["updated"].(time.Time)),
		Description:     wantResult["description"].(string),
		Items: []*gofeed.Item{{
			GUID:            wantItem["id"].(string),
			Title:           wantItem["title"].(string),
			Description:     wantItem["description"].(string),
			Content:         wantItem["content"].(string),
			Link:            wantItem["link"].(string),
			Authors:         []*gofeed.Person{{Name: wantItem["author"].(string)}},
			PublishedParsed: new(wantItem["created"].(time.Time)),
			UpdatedParsed:   new(wantItem["updated"].(time.Time)),
			Image: &gofeed.Image{
				Title: wantItemImage["title"].(string),
				URL:   wantItemImage["url"].(string),
			},
		}},
	}
	result, err := fn(sobek.Undefined(), vm.ToValue(feed))
	m := result.Export().(map[string]any)

	if assert.NoError(t, err) {
		assert.Equal(t, wantResult, m)
	}
}

func TestSetGetFn(t *testing.T) {
	wantResponseInfo := M{
		"status": http.StatusBadRequest,
		"header": http.Header{
			"Content-Length": []string{"21"},
			"Content-Type":   []string{"application/json"},
			"Date":           []string{time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")},
		},
		"body": `{"text":"hello test"}`,
	}
	wantHeaders := http.Header{
		"Accept":          {"application/json"},
		"Content-Type":    {"application/json"},
		"User-Agent":      {"unit test"},
		"Accept-Encoding": {"gzip"},
	}
	userAgent = "unit test"

	var headers http.Header
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers = r.Header

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"text":"hello test"}`))
	}))

	vm := sobek.New()
	err := SetGetFn(vm)
	assert.NoError(t, err)

	fn, ok := sobek.AssertFunction(vm.Get("get"))
	assert.True(t, ok)

	headerObj := vm.NewObject()
	headerObj.Set("Content-Type", "application/json")
	customObj := vm.NewObject()
	customObj.Set("headers", headerObj)
	result, err := fn(sobek.Undefined(), vm.ToValue(srv.URL), customObj)
	respInfo := result.Export().(map[string]any)

	assert.NoError(t, err)
	assert.Equal(t, wantHeaders, headers)
	assert.Equal(t, wantResponseInfo, respInfo)
}

func TestConvertMapString(t *testing.T) {
	tcs := []struct {
		Name   string
		Want   map[string]string
		Assert func(want map[string]string)
	}{
		{
			Name: "nil object",
			Want: map[string]string{},
			Assert: func(want map[string]string) {
				m := convertMapString(nil)

				assert.Equal(t, want, m)
			},
		},
		{
			Name: "empty object",
			Want: map[string]string{},
			Assert: func(want map[string]string) {
				vm := sobek.New()
				obj := vm.NewObject()

				m := convertMapString(obj)

				assert.Equal(t, want, m)
			},
		},
		{
			Name: "success",
			Want: map[string]string{
				"Content-Type": "application/json",
				"User-Agent":   "console",
			},
			Assert: func(want map[string]string) {
				vm := sobek.New()
				obj := vm.NewObject()
				obj.Set("Content-Type", "application/json")
				obj.Set("User-Agent", "console")

				m := convertMapString(obj)

				assert.Equal(t, want, m)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			tc.Assert(tc.Want)
		})
	}
}
