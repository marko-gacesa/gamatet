// Copyright (c) 2023-2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package render

import (
	"github.com/marko-gacesa/gamatet/graphics/geometry"
	"github.com/marko-gacesa/gamatet/graphics/material"
	"github.com/marko-gacesa/gamatet/graphics/texture"
)

type FieldResources struct {
	texManager *texture.Manager

	TexRock   uint32
	TexChain1 uint32
	TexChain2 uint32
	TexChain3 uint32

	MatTexUV material.Material
	MatNorm  material.Material
	MatRock  *material.Rock
	MatIron  *material.Iron
	MatLava  *material.Lava
	MatAcid  *material.Acid
	MatWave  *material.Curl
	MatColor *material.Color

	GeomSquareBack  geometry.Geometry
	GeomCube        geometry.Geometry
	GeomCubeDent    geometry.Geometry
	GeomFrame       geometry.Geometry
	GeomFrameThin   geometry.Geometry
	GeomRoundedCube geometry.Geometry
	GeomGem         geometry.Geometry
	GeomDie         geometry.Geometry
	GeomStar6       geometry.Geometry
	GeomStar8       geometry.Geometry
	GeomSphere      geometry.Geometry
}

func GenerateFieldResources(manager *texture.Manager) *FieldResources {
	const seed = 345

	rockImage := texture.GrayTex(seed)
	linkImage := texture.Link(seed)

	texRock := manager.Bind(rockImage)

	return &FieldResources{
		texManager: manager,

		TexRock:   texRock,
		TexChain1: manager.Bind(texture.Chain1(linkImage)),
		TexChain2: manager.Bind(texture.Chain2(linkImage)),
		TexChain3: manager.Bind(texture.Chain3(linkImage)),

		MatTexUV: material.TexUV(),
		MatNorm:  material.Normal(),
		MatRock:  material.NewRock(texRock),
		MatIron:  material.NewIron(texRock),
		MatLava:  material.NewLava(texRock),
		MatAcid:  material.NewAcid(texRock),
		MatWave:  material.NewCurl(texRock),
		MatColor: material.NewColor(),

		GeomSquareBack:  geometry.MakeSquareGeometry(geometry.CubeSideDent),
		GeomCube:        geometry.MakeCubeGeometry(geometry.CubeSideSimple),
		GeomCubeDent:    geometry.MakeCubeGeometry(geometry.CubeSideDent),
		GeomFrame:       geometry.MakeCubeGeometry(geometry.CubeSideFrame(1, 0.15)),
		GeomFrameThin:   geometry.MakeCubeGeometry(geometry.CubeSideFrame(0.9, 0.08)),
		GeomRoundedCube: geometry.MakeCubeGeometry(geometry.CubeSideRounded),
		GeomGem:         geometry.MakeCubeGeometry(geometry.CubeSideTruncated),
		GeomDie:         geometry.MakeCubeGeometry(geometry.CubeSideDie),
		GeomStar6:       geometry.MakeCubeGeometry(geometry.CubeSideHexagonalStar(0.25)),
		GeomStar8:       geometry.MakeCubeGeometry(geometry.CubeSideOctagonalStar(1)),
		GeomSphere:      geometry.MakeSphereGeometry(0.55, 16, 8),
	}
}

func (r FieldResources) Release() {
	r.GeomSquareBack.Delete()
	r.GeomCube.Delete()
	r.GeomCubeDent.Delete()
	r.GeomFrame.Delete()
	r.GeomFrameThin.Delete()
	r.GeomRoundedCube.Delete()
	r.GeomGem.Delete()
	r.GeomDie.Delete()
	r.GeomStar6.Delete()
	r.GeomStar8.Delete()
	r.GeomSphere.Delete()

	r.GeomSquareBack = nil
	r.GeomCube = nil
	r.GeomCubeDent = nil
	r.GeomFrame = nil
	r.GeomFrameThin = nil
	r.GeomRoundedCube = nil
	r.GeomGem = nil
	r.GeomDie = nil
	r.GeomStar6 = nil
	r.GeomStar8 = nil
	r.GeomSphere = nil

	r.MatTexUV.Delete()
	r.MatNorm.Delete()
	r.MatRock.Delete()
	r.MatIron.Delete()
	r.MatLava.Delete()
	r.MatAcid.Delete()
	r.MatWave.Delete()
	r.MatColor.Delete()

	r.texManager.Delete(r.TexRock)
	r.texManager.Delete(r.TexChain1)
	r.texManager.Delete(r.TexChain2)
	r.texManager.Delete(r.TexChain3)
}
