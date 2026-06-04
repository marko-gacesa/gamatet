// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package rendercache

import (
	"runtime"
	"testing"
)

func TestPointLightsPool(t *testing.T) {
	var mem1, mem2 runtime.MemStats

	f := func(name string) {
		runtime.ReadMemStats(&mem2)
		if allocations := mem2.Mallocs - mem1.Mallocs; allocations > 0 {
			t.Errorf("Expected no allocations at %s, but got %d", name, allocations)
		}
		mem1 = mem2
	}

	a := PointLightPool.Get()
	PointLightPool.Put(a)

	runtime.ReadMemStats(&mem1)

	a = PointLightPool.Get()
	f("1")
	PointLightPool.Put(a)
	f("2")
	a = PointLightPool.Get()
	f("3")
	PointLightPool.Put(a)
	f("4")
}
