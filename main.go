package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog/log"
)

var (
	configFilename = "config.yaml"
	userAgent      = ""
	addr           = "127.0.0.1:9000"
	httpClient     = resty.New()
	feedParser     = gofeed.NewParser()
)

func main() {
	if err := LoadConfigProgram(configFilename); err != nil {
		log.Fatal().Err(err).Msg("load configuration")
	}

	e := echo.New()
	e.Use(
		middleware.Recover(),
		middleware.RequestLogger(),
	)
	e.GET("reload", func(c *echo.Context) error {
		if err := LoadConfigProgram(configFilename); err != nil {
			return echo.ErrInternalServerError.Wrap(err)
		}

		return c.NoContent(http.StatusOK)
	})
	e.GET("feed/:name", func(c *echo.Context) error {
		name := c.Param("name")
		locker.Lock()
		program, ok := programs[name]
		locker.Unlock()
		if !ok {
			return echo.ErrNotFound
		}
		vm, err := NewVM()
		if err != nil {
			return echo.ErrInternalServerError.Wrap(err)
		}
		result, err := vm.RunProgram(program)
		if err != nil {
			return echo.ErrInternalServerError.Wrap(err)
		}
		data := result.Export()
		if data == nil {
			return echo.ErrInternalServerError.Wrap(errors.New("invalid data"))
		}
		items, ok := data.(map[string]any)
		if !ok {
			return echo.ErrInternalServerError.Wrap(errors.New("empty items"))
		}
		feed, err := convertToFeeds(items)
		if err != nil {
			return echo.ErrInternalServerError.Wrap(err)
		}

		atom, err := feed.ToAtom()
		if err != nil {
			return echo.ErrInternalServerError.Wrap(err)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationXMLCharsetUTF8)
		return c.String(http.StatusOK, atom)
	})

	sc := echo.StartConfig{
		Address:    addr,
		HideBanner: true,
		HidePort:   true,
	}
	if err := sc.Start(context.Background(), e); err != nil {
		log.Fatal().Err(err).Msgf("run server at %s", addr)
	}
}
