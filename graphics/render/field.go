// Copyright (c) 2020-2024 by Marko Gaćeša

package render

import (
	"gamatet/game/block"
	"gamatet/game/field"
	"gamatet/game/piece"
	"gamatet/graphics/gtypes"
	"gamatet/graphics/render/rendercache"
	"github.com/go-gl/mathgl/mgl32"
	"math"
	"time"
)

type FieldRender struct {
	resources *Resources
}

func NewFieldRenderer(globalResources *Resources) *FieldRender {
	return &FieldRender{
		resources: globalResources,
	}
}

func (f *FieldRender) Release() {
}

func (f *FieldRender) Render(r *Renderer, model *mgl32.Mat4, renderInfo *field.RenderInfo) {
	if renderInfo == nil {
		return
	}

	t := time.Since(gtypes.Time).Seconds()

	// light intensities
	const (
		lightIntLava    = 1.2
		lightIntAcid    = 1.2
		lightIntWave    = 1.5
		lightIntRuby    = 3
		lightPowShooter = 1.5
	)

	listWall := rendercache.ModelPool.Get()
	listBack := rendercache.ModelPool.Get()
	listRock := rendercache.ModelColorPool.Get()
	listRock1 := rendercache.ModelColorPool.Get()
	listRock2 := rendercache.ModelColorPool.Get()
	listRock3 := rendercache.ModelColorPool.Get()
	listIron := rendercache.ModelPool.Get()
	listLava := rendercache.ModelColorPool.Get()
	listAcid := rendercache.ModelColorPool.Get()
	listWave := rendercache.ModelColorPool.Get()
	listRuby := rendercache.ModelColorPool.Get()
	listShad := rendercache.ModelColorPool.Get()
	lights := rendercache.PointLightPool.Get()
	defer func() {
		rendercache.ModelPool.Put(listWall)
		rendercache.ModelPool.Put(listBack)
		rendercache.ModelColorPool.Put(listRock)
		rendercache.ModelColorPool.Put(listRock1)
		rendercache.ModelColorPool.Put(listRock2)
		rendercache.ModelColorPool.Put(listRock3)
		rendercache.ModelPool.Put(listIron)
		rendercache.ModelColorPool.Put(listLava)
		rendercache.ModelColorPool.Put(listAcid)
		rendercache.ModelColorPool.Put(listWave)
		rendercache.ModelColorPool.Put(listRuby)
		rendercache.ModelColorPool.Put(listShad)
		rendercache.PointLightPool.Put(lights)
	}()

	contentWidth := renderInfo.W + 4
	contentHeight := renderInfo.H + 2

	modelFrame := model.
		Mul4(mgl32.Translate3D(-float32(contentWidth)/2+0.5, -float32(contentHeight)/2+0.5, 0))

	modelField := modelFrame.
		Mul4(mgl32.Translate3D(3, 1, 0))

	modelNextBlocks := modelFrame.
		Mul4(mgl32.Translate3D(1, 1, 1))

	// prepare the field frame

	for x := 0; x < contentWidth; x++ {
		var m mgl32.Mat4

		m = modelFrame.Mul4(mgl32.Translate3D(float32(x), float32(0), 0))
		listWall.Add(m)

		m = modelFrame.Mul4(mgl32.Translate3D(float32(x), float32(contentHeight-1), 0))
		listWall.Add(m)
	}

	for y := 1; y < contentHeight-1; y++ {
		var m mgl32.Mat4

		m = modelFrame.Mul4(mgl32.Translate3D(float32(0), float32(y), 0))
		listWall.Add(m)
		m = modelFrame.Mul4(mgl32.Translate3D(float32(1), float32(y), 0))
		listWall.Add(m)
		m = modelFrame.Mul4(mgl32.Translate3D(float32(2), float32(y), 0))
		listWall.Add(m)

		m = modelFrame.Mul4(mgl32.Translate3D(float32(contentWidth-1), float32(y), 0))
		listWall.Add(m)
	}

	for x := 0; x < renderInfo.W; x++ {
		for y := 0; y < renderInfo.H; y++ {
			m := modelField.Mul4(mgl32.Translate3D(float32(x), float32(y), float32(-1)))
			listBack.Add(m)
		}
	}

	// prepare the field's blocks

	for _, fb := range renderInfo.Blocks {
		aniMatrix, aniColor := animListUpdate(&fb.Result)

		modelFieldBlock := modelField.
			Mul4(mgl32.Translate3D(float32(fb.X), float32(fb.Y), 0)).
			Mul4(aniMatrix)

		switch fb.Type {
		case block.TypeWall:
			listWall.Add(modelFieldBlock)
		case block.TypeIron:
			listIron.Add(modelFieldBlock)
		case block.TypeRuby:
			color := colorVector(fb.Block.Color)
			listRuby.Add(modelFieldBlock, color)
			lights.AddWithModel(modelFieldBlock, color.Vec3(), lightIntRuby)
		default:
			blockColor := colorVector(fb.Block.Color)
			color := mulColor(blockColor, aniColor)
			switch fb.Hardness {
			case 0:
				listRock.Add(modelFieldBlock, color)
			case 1:
				listRock1.Add(modelFieldBlock, color)
			case 2:
				listRock2.Add(modelFieldBlock, color)
			default:
				listRock3.Add(modelFieldBlock, color)
			}
		}
	}

	// prepare the pieces

	for _, p := range renderInfo.Pieces {
		if p.Empty {
			continue
		}

		lightPower := float32(1.0)
		if p.Type == piece.TypeShooter {
			lightPower *= lightPowShooter
		}

		aniMatrix, aniColor := animListUpdate(&p.Result)

		modelPiece := modelField.
			Mul4(mgl32.Translate3D(float32(p.X), float32(p.Y), 0)).
			Mul4(mgl32.Translate3D(float32(p.DimX)/2-0.5, -float32(p.DimY)/2+0.5, 0)).
			Mul4(aniMatrix).
			Mul4(mgl32.Translate3D(-float32(p.DimX)/2+0.5, float32(p.DimY)/2-0.5, 0))

		if p.Type == piece.TypeShooter {
			listRock.Add(modelPiece.
				Mul4(mgl32.Translate3D(0, 1, 0)).
				Mul4(mgl32.Scale3D(0.8, 0.8, 0.8)),
				colorVector(p.Blocks[0].Color))
		}

		for _, pb := range p.Blocks {
			modelPieceBlock := modelPiece.
				Mul4(mgl32.Translate3D(float32(pb.X), float32(pb.Y), 0))

			switch pb.Block.Type {
			case block.TypeAcid:
				listAcid.Add(modelPieceBlock, aniColor)
				lights.AddWithModel(modelPieceBlock, colorVector(block.Acid.Color).Vec3(), lightIntAcid*lightPower)
			case block.TypeLava:
				listLava.Add(modelPieceBlock, aniColor)
				lights.AddWithModel(modelPieceBlock, colorVector(block.Lava.Color).Vec3(), lightIntLava*lightPower)
			case block.TypeWave:
				listWave.Add(modelPieceBlock, aniColor)
				lights.AddWithModel(modelPieceBlock, colorVector(block.Wave.Color).Vec3(), lightIntWave*lightPower)
			default:
				blockColor := colorVector(pb.Block.Color)
				color := mulColor(blockColor, aniColor)
				switch pb.Hardness {
				case 0:
					listRock.Add(modelPieceBlock, color)
				case 1:
					listRock1.Add(modelPieceBlock, color)
				case 2:
					listRock2.Add(modelPieceBlock, color)
				default:
					listRock3.Add(modelPieceBlock, color)
				}
			}
		}
	}

	// prepare piece shadows

	scale := float32(0.7 + 0.3*math.Sin(math.Mod(10*t, math.Pi)))

	for _, p := range renderInfo.Pieces {
		if p.Empty || !p.DrawShadow {
			continue
		}

		for _, pb := range p.Shadow.Blocks {
			modelPieceShadowBlock := modelField.
				Mul4(mgl32.Translate3D(float32(pb.X), float32(pb.Y), 0)).
				Mul4(mgl32.Scale3D(scale, scale, scale))
			blockColor := colorVector(pb.Block.Color).Mul(0.7)
			listShad.Add(modelPieceShadowBlock, blockColor)
		}
	}

	// prepare next pieces

	for _, p := range renderInfo.Pieces {
		if p.Empty {
			continue
		}

		const dirY = 1

		var y float32
		for i, nb := range p.NextBlocks {
			dim, centerX, centerY := barycenter(nb)
			dimScale := 0.3/(float32(3*i)+1.0) + 0.2
			y += dirY * dimScale * dim / 2
			modelPieceN := modelNextBlocks.
				Mul4(mgl32.Translate3D(0, y, 0)).
				Mul4(mgl32.Scale3D(dimScale, dimScale, dimScale)).
				Mul4(mgl32.HomogRotate3DX(-0.4)).
				Mul4(mgl32.HomogRotate3DZ(float32(math.Mod(t, 2*math.Pi))))
			y += dirY * dimScale * (dim/2 + 0.7)

			for _, pb := range nb {
				modelPieceBlock := modelPieceN.
					Mul4(mgl32.Translate3D(float32(pb.X)-centerX, float32(pb.Y)-centerY, 0))

				switch pb.Block.Type {
				case block.TypeAcid:
					listAcid.Add(modelPieceBlock, colorWhite)
				case block.TypeLava:
					listLava.Add(modelPieceBlock, colorWhite)
				case block.TypeWave:
					listWave.Add(modelPieceBlock, colorWhite)
				default:
					color := colorVector(pb.Block.Color)
					switch pb.Hardness {
					case 0:
						listRock.Add(modelPieceBlock, color)
					case 1:
						listRock1.Add(modelPieceBlock, color)
					case 2:
						listRock2.Add(modelPieceBlock, color)
					default:
						listRock3.Add(modelPieceBlock, color)
					}
				}
			}
		}
	}

	// render all

	r.Material(f.resources.MatRock)
	r.Geometry(f.resources.GeomCube)
	f.resources.MatRock.Lights(lights)

	f.resources.MatRock.Color(colorVector(block.Wall.Color))
	for i := range listWall {
		r.Render(&listWall[i])
	}

	r.Geometry(f.resources.GeomDentCube)

	f.resources.MatRock.Color(colorVector(block.Wall.Color).Mul(0.6))
	for i := range listBack {
		r.Render(&listBack[i])
	}

	r.Geometry(f.resources.GeomRoundedCube)

	for i := range listRock {
		f.resources.MatRock.Color(listRock[i].Color)
		r.Render(&listRock[i].Model)
	}

	if len(listRock1) > 0 {
		f.resources.MatRock.ChainTexture(f.resources.TexChain1)
		for i := range listRock1 {
			f.resources.MatRock.Color(listRock1[i].Color)
			r.Render(&listRock1[i].Model)
		}
		f.resources.MatRock.ClearChain()
	}
	if len(listRock2) > 0 {
		f.resources.MatRock.ChainTexture(f.resources.TexChain2)
		for i := range listRock2 {
			f.resources.MatRock.Color(listRock2[i].Color)
			r.Render(&listRock2[i].Model)
		}
		f.resources.MatRock.ClearChain()
	}
	if len(listRock3) > 0 {
		f.resources.MatRock.ChainTexture(f.resources.TexChain3)
		for i := range listRock3 {
			f.resources.MatRock.Color(listRock3[i].Color)
			r.Render(&listRock3[i].Model)
		}
		f.resources.MatRock.ClearChain()
	}

	if len(listRuby) > 0 {
		r.Geometry(f.resources.GeomFrame)
		f.resources.MatRock.Color(colorVector(block.Wall.Color))
		for i := range listRuby {
			r.Render(&listRuby[i].Model)
		}
	}

	if len(listIron) > 0 {
		r.Geometry(f.resources.GeomCube)
		r.Material(f.resources.MatIron)
		f.resources.MatIron.Lights(lights)
		for i := range listIron {
			r.Render(&listIron[i])
		}
	}

	if len(listShad) > 0 {
		r.Geometry(f.resources.GeomShadowFrame)
		for i := range listShad {
			f.resources.MatRock.Color(listShad[i].Color)
			r.Render(&listShad[i].Model)
		}
	}

	if len(listLava) > 0 {
		r.Geometry(f.resources.GeomRoundedCube)
		r.Material(f.resources.MatLava)
		for i := range listLava {
			f.resources.MatLava.Color(listLava[i].Color)
			r.Render(&listLava[i].Model)
		}
	}

	if len(listAcid) > 0 {
		r.Geometry(f.resources.GeomRoundedCube)
		r.Material(f.resources.MatAcid)
		for i := range listAcid {
			f.resources.MatAcid.Color(listAcid[i].Color)
			r.Render(&listAcid[i].Model)
		}
	}

	if len(listWave) > 0 {
		r.Geometry(f.resources.GeomDie)
		r.Material(f.resources.MatWave)
		for i := range listWave {
			f.resources.MatWave.Color(listWave[i].Color)
			r.Render(&listWave[i].Model)
		}
	}

	if len(listRuby) > 0 {
		r.Geometry(f.resources.GeomStar8)
		r.Material(f.resources.MatColor)
		a := float32(math.Sin(3 * t))
		a = 0.3 + 0.3*a*a
		transform :=
			mgl32.Scale3D(a, a, a).
				Mul4(mgl32.HomogRotate3DZ(float32(6.7*t + 0.2))).
				Mul4(mgl32.HomogRotate3DX(float32(5.1*t + 0.7)))
		for i := range listRuby {
			f.resources.MatColor.Color(listRuby[i].Color)
			modelRuby := listRuby[i].Model.Mul4(transform)
			r.Render(&modelRuby)
		}
	}

	r.Material(f.resources.MatText)
	f.resources.MatText.Color(mgl32.Vec4{0, 0, 0, 1})

	text := "76543210"

	modelText := modelFrame.Mul4(mgl32.Translate3D(-0.3, 5.5, 0.6))
	for _, ch := range text {
		g := f.resources.GeomChar[byte(ch)]
		r.Geometry(g)
		r.Render(&modelText)
		modelText = modelText.Mul4(mgl32.Translate3D(g.Width(), 0, 0))
	}
}

func barycenter(blocks []block.XYB) (float32, float32, float32) {
	if len(blocks) == 0 {
		return 0, 0, 0
	}

	var cx, cy int
	minX, maxX, minY, maxY := math.MaxInt, 0, math.MaxInt, 0
	for i := range blocks {
		x := blocks[i].X
		y := blocks[i].Y
		cx += x
		cy += y
		minX = min(minX, x)
		maxX = max(maxX, x)
		minY = min(minY, y)
		maxY = max(maxY, y)
	}
	dim := max(maxX-minX+1, maxY-minY+1)
	return float32(dim),
		float32(cx) / float32(len(blocks)),
		float32(cy) / float32(len(blocks))
}
