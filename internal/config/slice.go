// Copyright (c) 2025 by Marko Gaćeša

package config

func SliceFixLen[T any](a []T, desiredLen int, genFn func(idx int) T) []T {
	if len(a) > desiredLen {
		return a[:desiredLen]
	}
	for i := len(a); i < desiredLen; i++ {
		a = append(a, genFn(i))
	}
	return a[:desiredLen:desiredLen]
}
