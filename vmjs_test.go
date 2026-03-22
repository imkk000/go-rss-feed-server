package main

import (
	"testing"

	"github.com/grafana/sobek"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
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
