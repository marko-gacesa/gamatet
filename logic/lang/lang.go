// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package lang

import (
	"slices"
	"sync"
	"sync/atomic"
)

type Lang string

var fallback *sync.Map // map[string]string
var current atomic.Pointer[sync.Map]
var supported *sync.Map // map[Lang]*sync.Map

func DefineFallback(m map[string]string) {
	if fallback != nil {
		panic("default language already set")
	}
	fallback = &sync.Map{}
	for k, v := range m {
		fallback.Store(k, v)
	}
}

func DefineFallbackFromExisting(l Lang) {
	if fallback != nil {
		panic("default language already set")
	}

	v, exists := supported.Load(l)
	if !exists {
		panic("language not defined")
	}

	fallback = v.(*sync.Map)
}

func Define(l Lang, m map[string]string) {
	if supported == nil {
		supported = &sync.Map{}
	}

	_, exists := supported.Load(l)
	if exists {
		panic("language already defined")
	}

	mm := &sync.Map{}
	for k, v := range m {
		mm.Store(k, v)
	}

	supported.Store(l, mm)
}

func Set(l Lang) {
	v, exists := supported.Load(l)
	if !exists {
		panic("language not defined")
	}

	mm := v.(*sync.Map)

	current.Store(mm)
}

func Supported() []Lang {
	var l []Lang
	supported.Range(func(k, v any) bool {
		l = append(l, k.(Lang))
		return true
	})
	slices.Sort(l)
	return l
}

func Str(key string) string {
	if curr := current.Load(); curr != nil {
		value, ok := curr.Load(key)
		if ok {
			return value.(string)
		}
	}

	if fallback != nil {
		value, ok := fallback.Load(key)
		if ok {
			return value.(string)
		}
	}

	return key
}

func StrInAll(key string) map[Lang]string {
	m := make(map[Lang]string)
	supported.Range(func(k, v any) bool {
		curr := v.(*sync.Map)
		v, ok := curr.Load(key)
		if ok {
			m[k.(Lang)] = v.(string)
		}
		return true
	})
	return m
}
