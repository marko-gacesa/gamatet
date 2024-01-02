// Copyright (c) 2023 by Marko Gaćeša

package render

type Scene struct {
	Objects []SceneObject
}

type SceneObject interface {
	StartPrepare()
	EndPrepare()
	Render()
}
