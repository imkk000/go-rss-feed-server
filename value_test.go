package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetValOrDefault(t *testing.T) {
	tcs := []struct {
		Name   string
		Want   int
		Assert func(want int)
	}{
		{
			Name: "get invalid type with default",
			Want: 20,
			Assert: func(want int) {
				v := getValOrDefault(uint64(10), 20)

				assert.Equal(t, want, v)
			},
		},
		{
			Name: "get invalid type no default",
			Want: 0,
			Assert: func(want int) {
				v := getValOrDefault[int](uint64(10))

				assert.Equal(t, want, v)
			},
		},
		{
			Name: "get val",
			Want: 10,
			Assert: func(want int) {
				v := getValOrDefault[int](int(10))

				assert.Equal(t, want, v)
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

func TestGetVal(t *testing.T) {
	tcs := []struct {
		Name   string
		Want   int
		Assert func(want int)
	}{
		{
			Name: "get nil",
			Want: 0,
			Assert: func(want int) {
				v := getVal[int](nil)

				assert.Equal(t, want, v)
			},
		},
		{
			Name: "get 10",
			Want: 10,
			Assert: func(want int) {
				v := getVal(new(int(10)))

				assert.Equal(t, want, v)
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

func TestGetElm(t *testing.T) {
	tcs := []struct {
		Name   string
		Want   int
		Assert func(want int)
	}{
		{
			Name: "get from nil",
			Want: 0,
			Assert: func(want int) {
				v := getElm[int](nil, 0)

				assert.Equal(t, want, v)
			},
		},
		{
			Name: "get from empty",
			Want: 0,
			Assert: func(want int) {
				v := getElm([]int{}, 0)

				assert.Equal(t, want, v)
			},
		},
		{
			Name: "get index greater than length",
			Want: 0,
			Assert: func(want int) {
				v := getElm([]int{10, 20}, 2)

				assert.Equal(t, want, v)
			},
		},
		{
			Name: "get index 0",
			Want: 10,
			Assert: func(want int) {
				v := getElm([]int{10, 20}, 0)

				assert.Equal(t, want, v)
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
