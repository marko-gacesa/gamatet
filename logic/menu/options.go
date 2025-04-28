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

func WithLabelFn(fn func() string) func(Item) {
	return func(item Item) {
		item.b().labelFn = fn
	}
}

func WithDescriptionFn(fn func() string) func(Item) {
	return func(item Item) {
		item.b().descriptionFn = fn
	}
}

func applyOptions(item Item, options ...func(Item)) {
	for _, opt := range options {
		opt(item)
	}
}
