// Copyright (c) 2025 by Marko Gaćeša

package cache

import (
	"time"
)

type String[T any] struct {
	// getFn is the function used the fetch the relevant value.
	getFn func() T

	// cmpAndLoadFn should compare and load the curr value to the prev if the values are different.
	// It should return true if the values were equal and false if they were not.
	cmpAndLoadFn func(prev *T, curr T) bool

	// formatFn is the function that converts the value to a string.
	formatFn func(T) string

	// cacheDur is duration to use already generated value (cached value) before fetching it again.
	cacheDur time.Duration

	valuesOld  T
	valuesStr  string
	valuesTime time.Time
}

func NewString[T any](
	getFn func() T,
	cmpAndLoadFn func(prev *T, curr T) bool,
	formatFn func(T) string,
	cacheDur time.Duration,
) *String[T] {
	return &String[T]{
		getFn:        getFn,
		cmpAndLoadFn: cmpAndLoadFn,
		formatFn:     formatFn,
		cacheDur:     cacheDur,
	}
}

func (c *String[T]) String() string {
	if c == nil || c.getFn == nil {
		return ""
	}

	if c.cacheDur > 0 {
		t := time.Now()
		if t.Sub(c.valuesTime) < c.cacheDur {
			return c.valuesStr
		}
		c.valuesTime = t
	}

	value := c.getFn()

	if equal := c.cmpAndLoadFn(&c.valuesOld, value); !equal {
		c.valuesStr = c.formatFn(value)
	}

	return c.valuesStr
}
