// Copyright (c) 2020 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package material

import "github.com/go-gl/mathgl/mgl32"

type PointLight struct {
	Position  mgl32.Vec3
	Color     mgl32.Vec3
	Intensity float32
}
