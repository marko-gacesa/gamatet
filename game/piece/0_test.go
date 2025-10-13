// Copyright (c) 2025 by Marko Gaćeša

package piece

import "fmt"

const XX = true
const __ = false

func outputShapes[T interface{ def() string }](a []T) {
	for _, v := range a {
		fmt.Println(v.def())
	}
}
