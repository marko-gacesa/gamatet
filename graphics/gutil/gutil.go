// Copyright (c) 2020 by Marko Gaćeša

package gutil

func CeilPow2(i int) int {
	i--
	i |= i >> 1
	i |= i >> 2
	i |= i >> 4
	i |= i >> 8
	i |= i >> 16
	return i + 1
}

func IsPow2(i int) bool { return (i & (i - 1)) == 0 }
