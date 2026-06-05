// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package profiling

import (
	"fmt"
	"os"
	"runtime"
)

func RuntimeInfo() {
	fmt.Fprintf(os.Stderr, ""+
		"Runtime.Version=%s CPUs=%d GOMAXPROCS=%d GOOS=%s GOARCH=%s\n",
		runtime.Version(), runtime.NumCPU(), runtime.GOMAXPROCS(0), runtime.GOOS, runtime.GOARCH,
	)
}

func RuntimeStats() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	fmt.Fprintf(os.Stderr, ""+
		"Goroutines=%d CgoCalls=%d\n"+
		"Heap: InUse=%d [%dK] Sys=%d [%dK]\n"+
		"Mallocs=%d Frees=%d Objects=%d\n"+
		"GC: Num=%d NumForced=%d TotalPause=%dms\n",
		runtime.NumGoroutine(), runtime.NumCgoCall(),
		mem.HeapAlloc, mem.HeapAlloc/1024, mem.Sys, mem.Sys/1024,
		mem.Mallocs, mem.Frees, mem.HeapObjects,
		mem.NumGC, mem.NumForcedGC, mem.PauseTotalNs/1000000,
	)
}

func DumpStack() {
	buf := make([]byte, 64<<10)
	numBytes := runtime.Stack(buf, true)
	os.Stderr.Write(buf[:numBytes])
}
