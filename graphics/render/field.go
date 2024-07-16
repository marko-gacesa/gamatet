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

var (
	colorWall = colorVector(block.Wall.Color)
	colorBack = colorVector(block.Wall.Color).Mul(0.6)
	colorLava = colorVector(block.Lava.Color)
	colorAcid = colorVector(block.Acid.Color)
	colorCurl = colorVector(block.Curl.Color)
	colorWave = colorVector(block.Wave.Color)
	colorBomb = colorVector(block.Bomb.Color)
)

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
		lightIntGoal    = 3
		lightPowShooter = 1.5
	)

	var colorsBack [field.MaxWidth]mgl32.Vec4
	var listsBack [field.MaxWidth]rendercache.Models
	for i := 0; i < renderInfo.W; i++ {
		listsBack[i] = rendercache.ModelPool.Get()
	}

	listWall := rendercache.ModelPool.Get()
	listRock := rendercache.ModelColorPool.Get()
	listRock1 := rendercache.ModelColorPool.Get()
	listRock2 := rendercache.ModelColorPool.Get()
	listRock3 := rendercache.ModelColorPool.Get()
	listRuby := rendercache.ModelColorPool.Get()
	listRuby1 := rendercache.ModelColorPool.Get()
	listRuby2 := rendercache.ModelColorPool.Get()
	listRuby3 := rendercache.ModelColorPool.Get()
	listIron := rendercache.ModelPool.Get()
	listLava := rendercache.ModelColorPool.Get()
	listAcid := rendercache.ModelColorPool.Get()
	listWave := rendercache.ModelColorPool.Get()
	listGoal := rendercache.ModelColorPool.Get()
	listFrame := rendercache.ModelPool.Get()
	listBomb := rendercache.ModelPool.Get()
	listShad := rendercache.ModelColorPool.Get()
	lights := rendercache.PointLightPool.Get()
	defer func() {
		for i := 0; i < renderInfo.W; i++ {
			rendercache.ModelPool.Put(listsBack[i])
		}
		rendercache.ModelPool.Put(listWall)
		rendercache.ModelColorPool.Put(listRock)
		rendercache.ModelColorPool.Put(listRock1)
		rendercache.ModelColorPool.Put(listRock2)
		rendercache.ModelColorPool.Put(listRock3)
		rendercache.ModelColorPool.Put(listRuby)
		rendercache.ModelColorPool.Put(listRuby1)
		rendercache.ModelColorPool.Put(listRuby2)
		rendercache.ModelColorPool.Put(listRuby3)
		rendercache.ModelPool.Put(listIron)
		rendercache.ModelColorPool.Put(listLava)
		rendercache.ModelColorPool.Put(listAcid)
		rendercache.ModelColorPool.Put(listWave)
		rendercache.ModelColorPool.Put(listGoal)
		rendercache.ModelPool.Put(listFrame)
		rendercache.ModelPool.Put(listBomb)
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

	// prepare the background

	for x := 0; x < renderInfo.W; x++ {
		colorCol := colorBack
		for pIdx := range renderInfo.Pieces {
			if !renderInfo.Pieces[pIdx].Empty {
				within := x >= renderInfo.Pieces[pIdx].Limits.Min && x <= renderInfo.Pieces[pIdx].Limits.Max
				shadow := x >= renderInfo.Pieces[pIdx].Shadow.ColL && x < renderInfo.Pieces[pIdx].Shadow.ColR
				if within {
					//colorCol = colorCol.Mul(1.1)
				}
				if shadow {
					colorCol = colorCol.Mul(1.1)
				}
			}
		}

		colorsBack[x] = colorCol

		for y := 0; y < renderInfo.H; y++ {
			m := modelField.Mul4(mgl32.Translate3D(float32(x), float32(y), float32(-1)))
			listsBack[x].Add(m)
		}
	}

	// prepare the field's blocks

	for _, fb := range renderInfo.Blocks {
		aniMatrix, aniColor := animListUpdate(&fb.Result)

		modelFieldBlockBase := modelField.Mul4(mgl32.Translate3D(float32(fb.X), float32(fb.Y), 0))
		modelFieldBlock := modelFieldBlockBase.Mul4(aniMatrix)

		switch fb.Type {
		case block.TypeWall:
			listWall.Add(modelFieldBlock)
		case block.TypeIron:
			listIron.Add(modelFieldBlock)
		case block.TypeRuby:
			blockColor := colorVector(fb.Block.Color)
			color := mulColor(blockColor, aniColor)
			switch fb.Hardness {
			case 0:
				listRuby.Add(modelFieldBlock, color)
			case 1:
				listRuby1.Add(modelFieldBlock, color)
			case 2:
				listRuby2.Add(modelFieldBlock, color)
			default:
				listRuby3.Add(modelFieldBlock, color)
			}
		case block.TypeAcid:
			listAcid.Add(modelFieldBlock, aniColor)
			lights.AddWithModel(modelFieldBlock, colorAcid.Vec3(), lightIntAcid*fb.Result.SX)
		case block.TypeLava:
			listLava.Add(modelFieldBlock, aniColor)
			lights.AddWithModel(modelFieldBlock, colorLava.Vec3(), lightIntLava*fb.Result.SX)
		case block.TypeCurl:
			listWave.Add(modelFieldBlock, mulColor(colorCurl, aniColor))
			lights.AddWithModel(modelFieldBlock, colorCurl.Vec3(), lightIntWave*fb.Result.SX)
		case block.TypeWave:
			listWave.Add(modelFieldBlock, mulColor(colorWave, aniColor))
			lights.AddWithModel(modelFieldBlock, colorWave.Vec3(), lightIntWave*fb.Result.SX)
		case block.TypeBomb:
			listBomb.Add(modelFieldBlock)
		case block.TypeGoal:
			listFrame.Add(modelFieldBlockBase)
			color := colorVector(fb.Block.Color)
			modelGoal := modelFieldBlock.Mul4(mgl32.Scale3D(0.6, 0.6, 0.6))
			listGoal.Add(modelGoal, color)
			lights.AddWithModel(modelFieldBlock, color.Vec3(), lightIntGoal)
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
				lights.AddWithModel(modelPieceBlock, colorAcid.Vec3(), lightIntAcid*lightPower)
			case block.TypeLava:
				listLava.Add(modelPieceBlock, aniColor)
				lights.AddWithModel(modelPieceBlock, colorLava.Vec3(), lightIntLava*lightPower)
			case block.TypeCurl:
				listWave.Add(modelPieceBlock, mulColor(colorCurl, aniColor))
				lights.AddWithModel(modelPieceBlock, colorCurl.Vec3(), lightIntWave*lightPower)
			case block.TypeWave:
				listWave.Add(modelPieceBlock, mulColor(colorWave, aniColor))
				lights.AddWithModel(modelPieceBlock, colorWave.Vec3(), lightIntWave*lightPower)
			case block.TypeBomb:
				listBomb.Add(modelPieceBlock)
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
				case block.TypeCurl:
					listWave.Add(modelPieceBlock, colorCurl)
				case block.TypeWave:
					listWave.Add(modelPieceBlock, colorWave)
				case block.TypeBomb:
					listBomb.Add(modelPieceBlock)
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

	f.resources.MatRock.Color(colorWall)
	for i := range listWall {
		r.Render(&listWall[i])
	}

	r.Geometry(f.resources.GeomDentCube)

	f.resources.MatRock.Color(colorBack)
	for i := range listsBack {
		if listsBack[i] == nil {
			break
		}
		f.resources.MatRock.Color(colorsBack[i])
		for j := range listsBack[i] {
			r.Render(&listsBack[i][j])
		}
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
		r.Geometry(f.resources.GeomGem)
		for i := range listRuby {
			f.resources.MatRock.Color(listRuby[i].Color)
			r.Render(&listRuby[i].Model)
		}
	}

	if len(listRuby1) > 0 {
		r.Geometry(f.resources.GeomGem)
		f.resources.MatRock.ChainTexture(f.resources.TexChain1)
		for i := range listRuby1 {
			f.resources.MatRock.Color(listRuby1[i].Color)
			r.Render(&listRuby1[i].Model)
		}
		f.resources.MatRock.ClearChain()
	}
	if len(listRuby2) > 0 {
		r.Geometry(f.resources.GeomGem)
		f.resources.MatRock.ChainTexture(f.resources.TexChain2)
		for i := range listRuby2 {
			f.resources.MatRock.Color(listRuby2[i].Color)
			r.Render(&listRuby2[i].Model)
		}
		f.resources.MatRock.ClearChain()
	}
	if len(listRuby3) > 0 {
		r.Geometry(f.resources.GeomGem)
		f.resources.MatRock.ChainTexture(f.resources.TexChain3)
		for i := range listRuby3 {
			f.resources.MatRock.Color(listRuby3[i].Color)
			r.Render(&listRuby3[i].Model)
		}
		f.resources.MatRock.ClearChain()
	}

	if len(listFrame) > 0 {
		r.Geometry(f.resources.GeomFrame)
		f.resources.MatRock.Color(colorWall)
		for i := range listFrame {
			r.Render(&listFrame[i])
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

	if len(listBomb) > 0 {
		r.Geometry(f.resources.GeomStar6)
		r.Material(f.resources.MatRock)
		f.resources.MatRock.Color(colorBomb)
		transform :=
			mgl32.HomogRotate3DZ(float32(6.7*t + 0.2)).
				Mul4(mgl32.HomogRotate3DX(float32(5.1*t + 0.7)))
		for i := range listBomb {
			modelBomb := listBomb[i].Mul4(transform)
			r.Render(&modelBomb)
		}
	}

	if len(listGoal) > 0 {
		r.Geometry(f.resources.GeomStar8)
		r.Material(f.resources.MatColor)
		for i := range listGoal {
			f.resources.MatColor.Color(listGoal[i].Color)
			r.Render(&listGoal[i].Model)
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
