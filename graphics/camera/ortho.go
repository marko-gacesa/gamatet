// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package camera

import "github.com/go-gl/mathgl/mgl32"

func (c *Camera) Orthogonal2D(w, h float32) {
	c.ResetLookAt()
	c.Projection(mgl32.Ortho2D(-w/2, w/2, -h/2, h/2))
}
