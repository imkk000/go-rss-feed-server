package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sync"
	"text/template"

	"github.com/goccy/go-yaml"
	"github.com/grafana/sobek"
)

type RssFeedConfig struct {
	UserAgent       string                    `yaml:"user_agent"`
	FeedItemsConfig map[string]FeedItemConfig `yaml:"feeds"`
}

type FeedItemConfig struct {
	File   string `yaml:"file,omitempty"`
	URL    string `yaml:"url,omitempty"`
	Script string `yaml:"script,omitempty"`
}

var (
	programs map[string]*sobek.Program
	locker   = new(sync.Mutex)
)

func LoadConfigProgram(configFilename string) (err error) {
	fs, err := os.Open(configFilename)
	if err != nil {
		return fmt.Errorf("read configuration file: %w", err)
	}
	defer func() {
		if ierr := fs.Close(); ierr != nil {
			ierr = fmt.Errorf("close file: %w", ierr)
			err = errors.Join(err, ierr)
		}
	}()

	var cfg RssFeedConfig
	if err := yaml.NewDecoder(fs).Decode(&cfg); err != nil {
		return fmt.Errorf("decode configuration: %w", err)
	}
	userAgent = cfg.UserAgent
	feedParser.UserAgent = cfg.UserAgent

	locker.Lock()
	defer locker.Unlock()
	programs = make(map[string]*sobek.Program)
	for name, c := range cfg.FeedItemsConfig {
		var content string

		if len(c.File) > 0 {
			raw, err := os.ReadFile(c.File)
			if err != nil {
				return fmt.Errorf("read %s: %w", c.File, err)
			}
			content = string(raw)
		} else if len(c.URL) > 0 {
			t, err := template.New("content").Parse(scriptTemplate)
			if err != nil {
				return fmt.Errorf("new template %s: %w", name, err)
			}
			var buf bytes.Buffer
			if err := t.Execute(&buf, c); err != nil {
				return fmt.Errorf("execute template %s: %w", name, err)
			}
			content = buf.String()
		}

		program, err := sobek.Compile(name, content, false)
		if err != nil {
			return fmt.Errorf("compile %s: %w", name, err)
		}
		programs[name] = program
	}

	return nil
}

const scriptTemplate = `
const url = "{{.URL}}";
let feeds = fetchFeed(url);

// inject here
{{.Script}};

convertMapFeed(feeds);
`
