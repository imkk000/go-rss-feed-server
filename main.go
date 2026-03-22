package main

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/rs/zerolog/log"
)

var (
	configFilename = "config.yaml"
	userAgent      = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36"
	addr           = "127.0.0.1:9000"
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

	if err := e.Start(addr); err != nil {
		log.Fatal().Err(err).Msgf("run server at %s", addr)
	}
}
