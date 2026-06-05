// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package profiling

import (
	"runtime"
	"strconv"

	"github.com/marko-gacesa/gamatet/logic/cache"
)

type Mallocs struct {
	cached *cache.String[uint64]
}

func NewMallocs() Mallocs {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return Mallocs{
		cached: cache.NewString[uint64](
			func() uint64 {
				var memNew runtime.MemStats
				runtime.ReadMemStats(&memNew)
				value := memNew.Mallocs - mem.Mallocs
				mem = memNew
				return value
			},
			func(v1 *uint64, v2 uint64) bool {
				equal := *v1 == v2
				*v1 = v2
				return equal
			},
			func(mallocs uint64) string { return "mallocs=" + strconv.FormatInt(int64(mallocs), 10) },
			0,
		),
	}
}

func (f Mallocs) String() string { return f.cached.String() }
