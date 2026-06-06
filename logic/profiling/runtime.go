// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package profiling

import (
	"fmt"
	"runtime"

	"github.com/marko-gacesa/gamatet/logic/cache"
)

type Runtime struct {
	cached *cache.String[runtimeData]
}

type runtimeData struct {
	NumGoroutine int
	NumGC        uint32
	Sys          uint64
	HeapObjects  uint64
	HeapAlloc    uint64
	MallocDelta  uint64
}

func NewRuntime() Runtime {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return Runtime{
		cached: cache.NewString[runtimeData](
			func() runtimeData {
				var memNew runtime.MemStats
				runtime.ReadMemStats(&memNew)

				mallocDelta := memNew.Mallocs - mem.Mallocs
				mem = memNew

				return runtimeData{
					NumGoroutine: runtime.NumGoroutine(),
					NumGC:        mem.NumGC,
					Sys:          mem.Sys,
					HeapObjects:  mem.HeapObjects,
					HeapAlloc:    mem.HeapAlloc,
					MallocDelta:  mallocDelta,
				}
			},
			func(v1 *runtimeData, v2 runtimeData) bool {
				equal := *v1 == v2
				*v1 = v2
				return equal
			},
			func(data runtimeData) string {
				return fmt.Sprintf(""+
					"Goroutines=%d Mem=%dK NumGC=%d\n"+
					"HeapObjects=%d\tHeapInUse=%dK\n"+
					"MallocDelta=%d\n",
					data.NumGoroutine, data.Sys/1024, data.NumGC,
					data.HeapObjects, data.HeapAlloc/1024,
					data.MallocDelta,
				)
			},
			0,
		),
	}
}

func (f Runtime) String() string { return f.cached.String() }
