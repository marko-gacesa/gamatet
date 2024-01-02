// Copyright (c) 2023 by Marko Gaćeša

package render

import (
	"gamatet/graphics/geometry"
	"gamatet/graphics/material"
	"gamatet/graphics/texture"
)

var Resources struct {
	TexRock uint32

	MatTexUV material.Material
	MatNorm  material.Material
	MatRock  material.Rock
	MatLava  material.Lava
	MatAcid  material.Acid

	GeomCube        geometry.Geometry
	GeomDentCube    geometry.Geometry
	GeomFrame       geometry.Geometry
	GeomRoundedCube geometry.Geometry
	GeomGem         geometry.Geometry
	GeomDie         geometry.Geometry
	GeomStar6       geometry.Geometry
	GeomStar8       geometry.Geometry
	GeomSphere      geometry.Geometry
}

func GenerateResources() {
	texture.Instance = texture.Init()

	//imgR := texture.SymbolCache.Symbol('R')
	img := texture.GrayTex(345)

	Resources.TexRock = texture.Instance.Bind(img)

	Resources.MatTexUV = material.TexUV()
	Resources.MatNorm = material.Normal()
	Resources.MatRock = material.NewRock(Resources.TexRock)
	Resources.MatLava = material.NewLava(Resources.TexRock)
	Resources.MatAcid = material.NewAcid(Resources.TexRock)

	Resources.GeomCube = geometry.MakeCubeGeometry(geometry.CubeSideSimple)
	Resources.GeomDentCube = geometry.MakeCubeGeometry(geometry.CubeSideDent)
	Resources.GeomFrame = geometry.MakeCubeGeometry(geometry.CubeSideFrame)
	Resources.GeomRoundedCube = geometry.MakeCubeGeometry(geometry.CubeSideRounded)
	Resources.GeomGem = geometry.MakeCubeGeometry(geometry.CubeSideTruncated)
	Resources.GeomDie = geometry.MakeCubeGeometry(geometry.CubeSideDie)
	Resources.GeomStar6 = geometry.MakeCubeGeometry(geometry.CubeSideHexagonalStar(0.5))
	Resources.GeomStar8 = geometry.MakeCubeGeometry(geometry.CubeSideOctagonalStar(1))
	Resources.GeomSphere = geometry.MakeSphereGeometry(0.55, 16, 8)
}

func ReleaseResources() {
	Resources.GeomCube.Delete()
	Resources.GeomDentCube.Delete()
	Resources.GeomFrame.Delete()
	Resources.GeomRoundedCube.Delete()
	Resources.GeomGem.Delete()
	Resources.GeomDie.Delete()
	Resources.GeomStar6.Delete()
	Resources.GeomStar8.Delete()
	Resources.GeomSphere.Delete()

	Resources.GeomCube = nil
	Resources.GeomDentCube = nil
	Resources.GeomFrame = nil
	Resources.GeomRoundedCube = nil
	Resources.GeomGem = nil
	Resources.GeomDie = nil
	Resources.GeomStar6 = nil
	Resources.GeomStar8 = nil
	Resources.GeomSphere = nil

	Resources.MatTexUV.Delete()
	Resources.MatNorm.Delete()
	Resources.MatRock.Delete()
	Resources.MatLava.Delete()
	Resources.MatAcid.Delete()

	texture.Instance.DeleteAll()
}
