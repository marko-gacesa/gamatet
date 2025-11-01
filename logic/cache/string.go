// Copyright (c) 2025 by Marko Gaćeša

package cache

import (
	"time"
)

type String[T any] struct {
	getFn    func() T
	equalFn  func(T, T) bool
	formatFn func(T) string
	cacheDur time.Duration

	valuesOld  T
	valuesStr  string
	valuesTime time.Time
}

func NewString[T any](getFn func() T, equalFn func(T, T) bool, formatFn func(T) string, cacheDur time.Duration) *String[T] {
	return &String[T]{
		getFn:    getFn,
		equalFn:  equalFn,
		formatFn: formatFn,
		cacheDur: cacheDur,
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

	equal := c.equalFn(c.valuesOld, value)

	if !equal {
		c.valuesOld = value
		c.valuesStr = c.formatFn(value)
	}

	return c.valuesStr
}
