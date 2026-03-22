package main

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/goccy/go-yaml"
	"github.com/grafana/sobek"
)

type RssFeedConfig struct {
	FeedItemsConfig map[string]FeedItemConfig `yaml:"feeds"`
}

type FeedItemConfig struct {
	File string `yaml:"file,omitempty"`
}

var (
	programs map[string]*sobek.Program
	locker   = new(sync.Mutex)
)

func LoadConfigProgram(configFilename string) (err error) {
	fs, err := os.OpenInRoot(".", configFilename)
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

	locker.Lock()
	defer locker.Unlock()
	programs = make(map[string]*sobek.Program)
	for name, c := range cfg.FeedItemsConfig {
		content, err := os.ReadFile(c.File)
		if err != nil {
			return fmt.Errorf("read %s: %w", c.File, err)
		}
		program, err := sobek.Compile(name, string(content), false)
		if err != nil {
			return fmt.Errorf("compile %s: %w", name, err)
		}
		programs[name] = program
	}

	return nil
}
