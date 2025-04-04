// Copyright (c) 2025 by Marko Gaćeša

package menu

func WithDisabled(fn func() bool) func(Item) {
	return func(item Item) {
		item.b().disabledFn = fn
	}
}

func WithVisible(fn func() bool) func(Item) {
	return func(item Item) {
		item.b().visibleFn = fn
	}
}

func applyOptions(item Item, options ...func(Item)) {
	for _, opt := range options {
		opt(item)
	}
}
