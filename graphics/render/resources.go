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

	TextCanvas *textcanvas.TextCanvas

	TexRock   uint32
	TexChain1 uint32
	TexChain2 uint32
	TexChain3 uint32
	TexText   uint32

	MatTexUV material.Material
	MatNorm  material.Material
	MatRock  *material.Rock
	MatLava  *material.Lava
	MatAcid  *material.Acid
	MatWave  *material.Curl
	MatText  *material.Text

	GeomCube        geometry.Geometry
	GeomDentCube    geometry.Geometry
	GeomFrame       geometry.Geometry
	GeomRoundedCube geometry.Geometry
	GeomGem         geometry.Geometry
	GeomDie         geometry.Geometry
	GeomStar6       geometry.Geometry
	GeomStar8       geometry.Geometry
	GeomSphere      geometry.Geometry
	GeomChar        map[byte]geometry.Text
	GeomText        map[string]geometry.Text
}

func GenerateResources(manager *texture.Manager, font *truetype.Font) Resources {
	const seed = 345

	rock := texture.GrayTex(seed)
	link := texture.Link(seed)

	face := textcanvas.NewFace(font, 32, 72)
	canvas := textcanvas.NewTextCanvas(512)
	geomChar := make(map[byte]geometry.Text)
	geomText := make(map[string]geometry.Text)
	for key := byte(33); key < 127; key++ {
		value := canvas.TextFloat32(string(rune(key)), face, color.RGBA{255, 255, 255, 128}, false)
		geomChar[key] = geometry.NewTextWithHeightAndScale(value, 1, 0.75)
	}
	for _, key := range []string{"ij"} {
		value := canvas.TextFloat32(key, face, color.RGBA{255, 255, 255, 128}, false)
		geomText[key] = geometry.NewTextWithHeightAndScale(value, 1, 0.75)
	}

	texRock := texture.Instance.Bind(rock)
	texText := texture.Instance.Bind(canvas.Image())

	return Resources{
		texManager: manager,

		TextCanvas: canvas,

		TexRock:   texRock,
		TexChain1: texture.Instance.Bind(texture.Chain1(link)),
		TexChain2: texture.Instance.Bind(texture.Chain2(link)),
		TexChain3: texture.Instance.Bind(texture.Chain3(link)),
		TexText:   texText,

		MatTexUV: material.TexUV(),
		MatNorm:  material.Normal(),
		MatRock:  material.NewRock(texRock),
		MatLava:  material.NewLava(texRock),
		MatAcid:  material.NewAcid(texRock),
		MatWave:  material.NewCurl(texRock),
		MatText:  material.NewText(texText),

		GeomCube:        geometry.MakeCubeGeometry(geometry.CubeSideSimple),
		GeomDentCube:    geometry.MakeCubeGeometry(geometry.CubeSideDent),
		GeomFrame:       geometry.MakeCubeGeometry(geometry.CubeSideFrame),
		GeomRoundedCube: geometry.MakeCubeGeometry(geometry.CubeSideRounded),
		GeomGem:         geometry.MakeCubeGeometry(geometry.CubeSideTruncated),
		GeomDie:         geometry.MakeCubeGeometry(geometry.CubeSideDie),
		GeomStar6:       geometry.MakeCubeGeometry(geometry.CubeSideHexagonalStar(0.5)),
		GeomStar8:       geometry.MakeCubeGeometry(geometry.CubeSideOctagonalStar(1)),
		GeomSphere:      geometry.MakeSphereGeometry(0.55, 16, 8),
		GeomChar:        geomChar,
		GeomText:        geomText,
	}
}

func (r Resources) Release() {
	r.GeomCube.Delete()
	r.GeomDentCube.Delete()
	r.GeomFrame.Delete()
	r.GeomRoundedCube.Delete()
	r.GeomGem.Delete()
	r.GeomDie.Delete()
	r.GeomStar6.Delete()
	r.GeomStar8.Delete()
	r.GeomSphere.Delete()
	for key, g := range r.GeomChar {
		g.Delete()
		delete(r.GeomChar, key)
	}
	for key, g := range r.GeomText {
		g.Delete()
		delete(r.GeomText, key)
	}

	r.GeomCube = nil
	r.GeomDentCube = nil
	r.GeomFrame = nil
	r.GeomRoundedCube = nil
	r.GeomGem = nil
	r.GeomDie = nil
	r.GeomStar6 = nil
	r.GeomStar8 = nil
	r.GeomSphere = nil

	r.MatTexUV.Delete()
	r.MatNorm.Delete()
	r.MatRock.Delete()
	r.MatLava.Delete()
	r.MatAcid.Delete()
	r.MatWave.Delete()
	r.MatText.Delete()

	texture.Instance.Delete(r.TexRock)
	texture.Instance.Delete(r.TexChain1)
	texture.Instance.Delete(r.TexChain2)
	texture.Instance.Delete(r.TexChain3)
	texture.Instance.Delete(r.TexText)

	r.TextCanvas.Clear()
}
