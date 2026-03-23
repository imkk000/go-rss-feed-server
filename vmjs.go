package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/grafana/sobek"
	"github.com/mmcdole/gofeed"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func NewVM() (*sobek.Runtime, error) {
	funcs := []NewSobekFn{
		SetExitFn,
		SetConsoleFn,
		SetGetFn,
		SetFetchFeedFn,
		SetConvertMapFeedFn,
		SetParseHTMLFn,
		SetParseMarkdownAwesomeFn,
	}
	vm := sobek.New()
	for _, fn := range funcs {
		if err := fn(vm); err != nil {
			return nil, err
		}
	}

	return vm, nil
}

func parseAwesomeLink(text string) []string {
	// Handle: [title](link) - description
	// Or: [title](link)
	if matches := regexp.MustCompile(`\[(.+?)\]\((.+?)\)(?:\s*-\s*(.+))?`).FindStringSubmatch(text); matches != nil {
		desc := matches[3]
		if desc == "" {
			desc = matches[1] // Use title if no description
		}
		return []string{matches[1], matches[2], desc}
	}
	return nil
}

func SetParseMarkdownAwesomeFn(vm *sobek.Runtime) error {
	fn := func(source []byte) ([]Awesome, error) {
		md := goldmark.New()
		doc := md.Parser().Parse(text.NewReader(source))

		var awesomes []Awesome
		var currentTopic []string

		ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
			if !entering {
				return ast.WalkContinue, nil
			}

			// Track heading hierarchy
			if h, ok := n.(*ast.Heading); ok {
				text := string(h.Lines().Value(source))
				level := h.Level - 2
				if level >= 0 && level < 3 {
					currentTopic = currentTopic[:level]
					currentTopic = append(currentTopic, text)
				}
			}

			// Parse list items
			if item, ok := n.(*ast.ListItem); ok {
				if block, ok := item.FirstChild().(*ast.TextBlock); ok && block.ChildCount() > 1 {
					text := string(block.Lines().Value(source))
					if parts := parseAwesomeLink(text); parts != nil {
						awesomes = append(awesomes, Awesome{
							Topic:       strings.Join(currentTopic, " > "),
							Title:       parts[0],
							Link:        parts[1],
							Description: parts[2],
						})
					}
				}
			}

			return ast.WalkContinue, nil
		})

		return awesomes, nil
	}
	if err := vm.Set("parseMarkdownAwesome", fn); err != nil {
		return fmt.Errorf("set parse markdown awesome func: %w", err)
	}

	return nil
}

func createSelectionWrapper(vm *sobek.Runtime, sel *goquery.Selection) sobek.Value {
	obj := vm.NewObject()
	obj.Set("length", sel.Length())
	obj.Set("find", func(selector string) sobek.Value {
		return createSelectionWrapper(vm, sel.Find(selector))
	})
	obj.Set("text", func() string {
		return sel.Text()
	})
	obj.Set("attr", func(attr string) string {
		val, _ := sel.Attr(attr)
		return val
	})
	obj.Set("each", func(callback sobek.Callable) {
		sel.Each(func(i int, s *goquery.Selection) {
			callback(sobek.Undefined(), vm.ToValue(i), createSelectionWrapper(vm, s))
		})
	})

	return obj
}

func SetParseHTMLFn(vm *sobek.Runtime) error {
	fn := func(html string) (sobek.Value, error) {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			return nil, fmt.Errorf("new document: %w", err)
		}

		return createSelectionWrapper(vm, doc.Selection), nil
	}
	if err := vm.Set("parseHTML", fn); err != nil {
		return fmt.Errorf("set parse html func: %w", err)
	}

	return nil
}

func SetConvertMapFeedFn(vm *sobek.Runtime) error {
	fn := func(feed *gofeed.Feed) (sobek.Value, error) {
		result := M{
			"title":       feed.Title,
			"link":        feed.Link,
			"author":      getVal(getElm(feed.Authors, 0)).Name,
			"description": feed.Description,
			"created":     getVal(feed.PublishedParsed),
			"updated":     getVal(feed.UpdatedParsed),
		}
		if len(feed.Items) > 0 {
			items := make([]any, len(feed.Items))
			for i, item := range feed.Items {
				m := M{
					"id":          item.GUID,
					"title":       item.Title,
					"description": item.Description,
					"content":     item.Content,
					"link":        item.Link,
					"author":      getVal(getElm(item.Authors, 0)).Name,
					"created":     getVal(item.PublishedParsed),
					"updated":     getVal(item.UpdatedParsed),
				}
				if item.Image != nil {
					m["image"] = M{
						"title": item.Image.Title,
						"url":   item.Image.URL,
					}
				}
				items[i] = m
			}

			result["items"] = items
		}

		return vm.ToValue(result), nil
	}
	if err := vm.Set("convertMapFeed", fn); err != nil {
		return fmt.Errorf("set convert feed func: %w", err)
	}

	return nil
}

func SetFetchFeedFn(vm *sobek.Runtime) error {
	fn := func(url string) (sobek.Value, error) {
		feed, err := feedParser.ParseURL(url)
		if err != nil {
			return nil, fmt.Errorf("fetch feed: %w", err)
		}

		return vm.ToValue(feed), nil
	}
	if err := vm.Set("fetchFeed", fn); err != nil {
		return fmt.Errorf("set fetch feed func: %w", err)
	}

	return nil
}

func SetGetFn(vm *sobek.Runtime) error {
	fn := func(url string, options *sobek.Object) (sobek.Value, error) {
		var headers map[string]string
		if options != nil {
			if obj := options.Get("headers"); !sobek.IsNull(obj) {
				headers = convertMapString(obj.ToObject(vm))
			}
		}

		r := httpClient.
			R().
			SetHeader("User-Agent", userAgent).
			SetHeaders(headers)

		resp, err := r.Get(url)
		if err != nil {
			return nil, fmt.Errorf("do request: %w", err)
		}

		return vm.ToValue(M{
			"status": resp.StatusCode(),
			"header": resp.Header(),
			"body":   string(resp.Body()),
		}), nil
	}
	if err := vm.Set("get", fn); err != nil {
		return fmt.Errorf("set get func: %w", err)
	}

	return nil
}

func SetConsoleFn(vm *sobek.Runtime) error {
	console := vm.NewObject()
	fn := func(v any) {
		fmt.Println(v)
	}
	if err := console.Set("log", fn); err != nil {
		return fmt.Errorf("new console log: %w", err)
	}
	if err := vm.Set("console", console); err != nil {
		return fmt.Errorf("set console: %w", err)
	}

	return nil
}

func SetExitFn(vm *sobek.Runtime) error {
	fn := func(msg string) {
		panic(vm.ToValue(msg))
	}
	if err := vm.Set("exit", fn); err != nil {
		return fmt.Errorf("set func: %w", err)
	}

	return nil
}

func convertMapString(obj *sobek.Object) map[string]string {
	result := make(map[string]string)
	if obj == nil {
		return result
	}
	for _, key := range obj.Keys() {
		result[key] = obj.Get(key).String()
	}

	return result
}

type (
	NewSobekFn func(vm *sobek.Runtime) error
	M          = map[string]any
)

type Awesome struct {
	Topic       string
	Title       string
	Link        string
	Description string
}
