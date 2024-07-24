// Copyright (c) 2023,2024 by Marko Gaćeša

package render

import (
	"gamatet/graphics/geometry"
	"gamatet/graphics/material"
	"gamatet/graphics/textcanvas"
	"gamatet/graphics/texture"
	"github.com/golang/freetype/truetype"
	"image/color"
)

type Resources struct {
	texManager *texture.Manager

	TextCanvas  *textcanvas.TextCanvas
	TextRectMap map[byte]textcanvas.RectUV

	TexRock   uint32
	TexChain1 uint32
	TexChain2 uint32
	TexChain3 uint32
	TexText   uint32

	MatTexUV material.Material
	MatNorm  material.Material
	MatRock  *material.Rock
	MatIron  *material.Iron
	MatLava  *material.Lava
	MatAcid  *material.Acid
	MatWave  *material.Curl
	MatColor *material.Color
	MatText  *material.Text

	GeomSquare      geometry.Geometry
	GeomSquare0     geometry.Geometry
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

func GenerateResources(manager *texture.Manager, font *truetype.Font) *Resources {
	const seed = 345

	rock := texture.GrayTex(seed)
	link := texture.Link(seed)

	textRectMap := make(map[byte]textcanvas.RectUV)

	face := textcanvas.NewFace(font, 32, 72)
	canvas := textcanvas.NewTextCanvas(512)
	for key := byte(33); key < 127; key++ {
		textRectMap[key] = canvas.TextUV(string(rune(key)), face, color.White, false)
	}

	texRock := texture.Instance.Bind(rock)
	texText := texture.Instance.Bind(canvas.Image())

	return &Resources{
		texManager: manager,

		TextCanvas:  canvas,
		TextRectMap: textRectMap,

		TexRock:   texRock,
		TexChain1: texture.Instance.Bind(texture.Chain1(link)),
		TexChain2: texture.Instance.Bind(texture.Chain2(link)),
		TexChain3: texture.Instance.Bind(texture.Chain3(link)),
		TexText:   texText,

		MatTexUV: material.TexUV(),
		MatNorm:  material.Normal(),
		MatRock:  material.NewRock(texRock),
		MatIron:  material.NewIron(texRock),
		MatLava:  material.NewLava(texRock),
		MatAcid:  material.NewAcid(texRock),
		MatWave:  material.NewCurl(texRock),
		MatColor: material.NewColor(),
		MatText:  material.NewText(texText),

		GeomSquare:      geometry.NewSquare(),
		GeomSquare0:     geometry.NewSquare0(),
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

func (r Resources) Release() {
	r.GeomSquare.Delete()
	r.GeomSquare0.Delete()
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

	r.GeomSquare = nil
	r.GeomSquare0 = nil
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
	r.MatText.Delete()

	texture.Instance.Delete(r.TexRock)
	texture.Instance.Delete(r.TexChain1)
	texture.Instance.Delete(r.TexChain2)
	texture.Instance.Delete(r.TexChain3)
	texture.Instance.Delete(r.TexText)

	r.TextCanvas.Clear()

	for k := range r.TextRectMap {
		delete(r.TextRectMap, k)
	}
}
