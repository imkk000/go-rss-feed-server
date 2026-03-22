package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigProgram(t *testing.T) {
	var noErr error

	tcs := []struct {
		Name    string
		WantErr error
		Assert  func(wantErr error)
	}{
		{
			Name:    "failed",
			WantErr: errors.New("SyntaxError: myfeed: Line 1:1 Unexpected token ."),
			Assert: func(wantErr error) {
				err := LoadConfigProgram("./testdata/failed.yaml")
				err = errors.Unwrap(err)

				assert.EqualError(t, err, wantErr.Error())
			},
		},
		{
			Name:    "success",
			WantErr: noErr,
			Assert: func(wantErr error) {
				err := LoadConfigProgram("./testdata/success.yaml")

				assert.ErrorIs(t, err, wantErr)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			tc.Assert(tc.WantErr)
		})
	}
}
