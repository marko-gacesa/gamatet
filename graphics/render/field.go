// Copyright (c) 2020-2024 by Marko Gaćeša

package render

import (
	"context"
	"gamatet/game/block"
	"gamatet/game/core"
	"gamatet/game/field"
	"gamatet/game/piece"
	"gamatet/graphics/render/rendercache"
	"github.com/go-gl/mathgl/mgl32"
	"math"
	"time"
)

var (
	colorWall      = colorVector(block.Wall.Color)
	colorBack      = colorVector(block.Wall.Color).Mul(0.6)
	colorLava      = colorVector(block.Lava.Color)
	colorAcid      = colorVector(block.Acid.Color)
	colorCurl      = colorVector(block.Curl.Color)
	colorWave      = colorVector(block.Wave.Color)
	colorBomb      = colorVector(block.Bomb.Color)
	colorBackMulti = []mgl32.Vec4{
		colorVector(0x0000FFFF).Mul(0.6),
		colorVector(0xFF0000FF).Mul(0.6),
		colorVector(0x00FF00FF).Mul(0.6),
		colorVector(0xFF8080FF).Mul(0.6),
	}
)

var t0 = time.Now()

func GetExtendedContent(w, h int) (int, int) {
	return w + 4, // left frame, game info (2), right frame
		h + 2 // top frame, bottom frame
}

type Field struct {
	model     mgl32.Mat4
	resources *FieldResources
	text      *Text

	renderRequesterFieldIdx int
	renderRequester         core.RenderRequester

	renderInfo    *field.RenderInfo
	renderInfoCh  chan *field.RenderInfo
	prepareDoneCh chan struct{}

	t          float64 // time
	w          int     // field width
	colorsBack [field.MaxWidth]mgl32.Vec4
	listsBack  [field.MaxWidth]rendercache.Models

	listWall  rendercache.Models
	listIron  rendercache.Models
	listFrame rendercache.Models
	listBomb  rendercache.Models
	listRock  rendercache.ModelColorValueList
	listRuby  rendercache.ModelColorValueList
	listLava  rendercache.ModelColorList
	listAcid  rendercache.ModelColorList
	listWave  rendercache.ModelColorList
	listGoal  rendercache.ModelColorList
	listShad  rendercache.ModelColorList
	listAmmo  rendercache.ModelColorValueList
	lights    rendercache.PointLights
}

func NewField(
	model mgl32.Mat4,
	resources *FieldResources,
	text *Text,
	renderRequesterFieldIdx int,
	renderRequester core.RenderRequester,
) *Field {
	return &Field{
		model:                   model,
		resources:               resources,
		text:                    text,
		renderRequesterFieldIdx: renderRequesterFieldIdx,
		renderRequester:         renderRequester,
		renderInfoCh:            make(chan *field.RenderInfo),
		prepareDoneCh:           make(chan struct{}),
	}
}

func (f *Field) Prepare(ctx context.Context, now time.Time) {
	f.renderRequester.RenderRequest(ctx, f.renderRequesterFieldIdx, now, f.renderInfoCh)
	go func() {
		defer func() { f.prepareDoneCh <- struct{}{} }()
		select {
		case <-ctx.Done():
		case renderInfo := <-f.renderInfoCh:
			if renderInfo == nil {
				return
			}
			f.preRender(renderInfo, now)
			f.prepareModels(renderInfo)
			field.ReturnRenderInfo(f.renderInfo)
		}
	}()
}

func (f *Field) Render(r *Renderer) {
	<-f.prepareDoneCh
	f.renderAll(r)
	f.postRender()
}

func (f *Field) preRender(renderInfo *field.RenderInfo, now time.Time) {
	if renderInfo == nil {
		return
	}

	f.t = now.Sub(t0).Seconds()

	w := renderInfo.W
	for i := 0; i < w; i++ {
		f.listsBack[i] = rendercache.ModelPool.Get()
	}
	f.w = w

	f.listWall = rendercache.ModelPool.Get()
	f.listIron = rendercache.ModelPool.Get()
	f.listFrame = rendercache.ModelPool.Get()
	f.listBomb = rendercache.ModelPool.Get()
	f.listRock = rendercache.ModelColorValuePool.Get()
	f.listRuby = rendercache.ModelColorValuePool.Get()
	f.listLava = rendercache.ModelColorPool.Get()
	f.listAcid = rendercache.ModelColorPool.Get()
	f.listWave = rendercache.ModelColorPool.Get()
	f.listGoal = rendercache.ModelColorPool.Get()
	f.listShad = rendercache.ModelColorPool.Get()
	f.listAmmo = rendercache.ModelColorValuePool.Get()
	f.lights = rendercache.PointLightPool.Get()
}

func (f *Field) postRender() {
	for i, w := 0, f.w; i < w; i++ {
		rendercache.ModelPool.Put(f.listsBack[i])
	}

	rendercache.ModelPool.Put(f.listWall)
	rendercache.ModelPool.Put(f.listIron)
	rendercache.ModelPool.Put(f.listFrame)
	rendercache.ModelPool.Put(f.listBomb)
	rendercache.ModelColorValuePool.Put(f.listRock)
	rendercache.ModelColorValuePool.Put(f.listRuby)
	rendercache.ModelColorPool.Put(f.listLava)
	rendercache.ModelColorPool.Put(f.listAcid)
	rendercache.ModelColorPool.Put(f.listWave)
	rendercache.ModelColorPool.Put(f.listGoal)
	rendercache.ModelColorPool.Put(f.listShad)
	rendercache.ModelColorValuePool.Put(f.listAmmo)
	rendercache.PointLightPool.Put(f.lights)
}

func (f *Field) prepareModels(renderInfo *field.RenderInfo) {
	model := &f.model

	// light intensities
	const (
		lightIntLava    = 1.2
		lightIntAcid    = 1.2
		lightIntWave    = 1.5
		lightIntGoal    = 3
		lightPowShooter = 1.5
	)

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
		f.listWall.Add(m)

		m = modelFrame.Mul4(mgl32.Translate3D(float32(x), float32(contentHeight-1), 0))
		f.listWall.Add(m)
	}

	for y := 1; y < contentHeight-1; y++ {
		var m mgl32.Mat4

		m = modelFrame.Mul4(mgl32.Translate3D(float32(0), float32(y), 0))
		f.listWall.Add(m)
		m = modelFrame.Mul4(mgl32.Translate3D(float32(1), float32(y), 0))
		f.listWall.Add(m)
		m = modelFrame.Mul4(mgl32.Translate3D(float32(2), float32(y), 0))
		f.listWall.Add(m)

		m = modelFrame.Mul4(mgl32.Translate3D(float32(contentWidth-1), float32(y), 0))
		f.listWall.Add(m)
	}

	// prepare the background

	for x := 0; x < renderInfo.W; x++ {
		colorCol := colorBack
		for pIdx := range renderInfo.Pieces {
			if !renderInfo.Pieces[pIdx].Empty {
				within := x >= renderInfo.Pieces[pIdx].Limits.Min && x <= renderInfo.Pieces[pIdx].Limits.Max
				shadow := x >= renderInfo.Pieces[pIdx].Shadow.ColL && x < renderInfo.Pieces[pIdx].Shadow.ColR
				if within {
					colorCol = colorCol.Add(colorBackMulti[pIdx])
				}
				if shadow {
					colorCol = colorCol.Mul(1.1)
				}
			}
		}

		f.colorsBack[x] = colorCol

		for y := 0; y < renderInfo.H; y++ {
			m := modelField.Mul4(mgl32.Translate3D(float32(x), float32(y), float32(-1)))
			f.listsBack[x].Add(m)
		}
	}

	// prepare the field's blocks

	for _, fb := range renderInfo.Blocks {
		aniMatrix, aniColor := animListUpdate(&fb.Result)

		modelFieldBlockBase := modelField.Mul4(mgl32.Translate3D(float32(fb.X), float32(fb.Y), 0))
		modelFieldBlock := modelFieldBlockBase.Mul4(aniMatrix)

		switch fb.Type {
		case block.TypeWall:
			f.listWall.Add(modelFieldBlock)
		case block.TypeIron:
			f.listIron.Add(modelFieldBlock)
		case block.TypeRuby:
			blockColor := colorVector(fb.Block.Color)
			color := mulColor(blockColor, aniColor)
			f.listRuby.Add(modelFieldBlock, color, int(fb.Hardness))
		case block.TypeAcid:
			f.listAcid.Add(modelFieldBlock, aniColor)
			f.lights.AddWithModel(modelFieldBlock, colorAcid.Vec3(), lightIntAcid*fb.Result.SX)
		case block.TypeLava:
			f.listLava.Add(modelFieldBlock, aniColor)
			f.lights.AddWithModel(modelFieldBlock, colorLava.Vec3(), lightIntLava*fb.Result.SX)
		case block.TypeCurl:
			f.listWave.Add(modelFieldBlock, mulColor(colorCurl, aniColor))
			f.lights.AddWithModel(modelFieldBlock, colorCurl.Vec3(), lightIntWave*fb.Result.SX)
		case block.TypeWave:
			f.listWave.Add(modelFieldBlock, mulColor(colorWave, aniColor))
			f.lights.AddWithModel(modelFieldBlock, colorWave.Vec3(), lightIntWave*fb.Result.SX)
		case block.TypeBomb:
			f.listBomb.Add(modelFieldBlock)
		case block.TypeGoal:
			f.listFrame.Add(modelFieldBlockBase)
			color := colorVector(fb.Block.Color)
			modelGoal := modelFieldBlock.Mul4(mgl32.Scale3D(0.6, 0.6, 0.6))
			f.listGoal.Add(modelGoal, color)
			f.lights.AddWithModel(modelFieldBlock, color.Vec3(), lightIntGoal)
		default:
			blockColor := colorVector(fb.Block.Color)
			color := mulColor(blockColor, aniColor)
			f.listRock.Add(modelFieldBlock, color, int(fb.Hardness))
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
			modelAmmo := modelPiece.
				Mul4(mgl32.Translate3D(0, 1.1, 0)).
				Mul4(mgl32.Scale3D(0.8, 0.8, 0.8))
			f.listAmmo.Add(
				modelAmmo,
				colorVector(p.Blocks[0].Color),
				p.ActCount)
		}

		for _, pb := range p.Blocks {
			modelPieceBlock := modelPiece.
				Mul4(mgl32.Translate3D(float32(pb.X), float32(pb.Y), 0))

			switch pb.Block.Type {
			case block.TypeAcid:
				f.listAcid.Add(modelPieceBlock, aniColor)
				f.lights.AddWithModel(modelPieceBlock, colorAcid.Vec3(), lightIntAcid*lightPower)
			case block.TypeLava:
				f.listLava.Add(modelPieceBlock, aniColor)
				f.lights.AddWithModel(modelPieceBlock, colorLava.Vec3(), lightIntLava*lightPower)
			case block.TypeCurl:
				f.listWave.Add(modelPieceBlock, mulColor(colorCurl, aniColor))
				f.lights.AddWithModel(modelPieceBlock, colorCurl.Vec3(), lightIntWave*lightPower)
			case block.TypeWave:
				f.listWave.Add(modelPieceBlock, mulColor(colorWave, aniColor))
				f.lights.AddWithModel(modelPieceBlock, colorWave.Vec3(), lightIntWave*lightPower)
			case block.TypeBomb:
				f.listBomb.Add(modelPieceBlock)
			default:
				blockColor := colorVector(pb.Block.Color)
				color := mulColor(blockColor, aniColor)
				f.listRock.Add(modelPieceBlock, color, int(pb.Hardness))
			}
		}
	}

	// prepare piece shadows

	scale := float32(0.7 + 0.3*math.Sin(math.Mod(10*f.t, math.Pi)))

	for _, p := range renderInfo.Pieces {
		if p.Empty || !p.DrawShadow {
			continue
		}

		for _, pb := range p.Shadow.Blocks {
			modelPieceShadowBlock := modelField.
				Mul4(mgl32.Translate3D(float32(pb.X), float32(pb.Y), 0)).
				Mul4(mgl32.Scale3D(scale, scale, scale))
			blockColor := colorVector(pb.Block.Color).Mul(0.7)
			f.listShad.Add(modelPieceShadowBlock, blockColor)
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
				Mul4(mgl32.HomogRotate3DZ(float32(math.Mod(f.t, 2*math.Pi))))
			y += dirY * dimScale * (dim/2 + 0.7)

			for _, pb := range nb {
				modelPieceBlock := modelPieceN.
					Mul4(mgl32.Translate3D(float32(pb.X)-centerX, float32(pb.Y)-centerY, 0))

				switch pb.Block.Type {
				case block.TypeAcid:
					f.listAcid.Add(modelPieceBlock, colorWhite)
				case block.TypeLava:
					f.listLava.Add(modelPieceBlock, colorWhite)
				case block.TypeCurl:
					f.listWave.Add(modelPieceBlock, colorCurl)
				case block.TypeWave:
					f.listWave.Add(modelPieceBlock, colorWave)
				case block.TypeBomb:
					f.listBomb.Add(modelPieceBlock)
				default:
					color := colorVector(pb.Block.Color)
					f.listRock.Add(modelPieceBlock, color, int(pb.Hardness))
				}
			}
		}
	}

	// sort

	f.listRock.OrderByValue()
	f.listRuby.OrderByValue()
}

func (f *Field) renderAll(r *Renderer) {
	r.Material(f.resources.MatRock)
	r.Geometry(f.resources.GeomCube)
	f.resources.MatRock.Lights(f.lights)

	f.resources.MatRock.Color(colorWall)
	for i := range f.listWall {
		r.Render(&f.listWall[i])
	}

	r.Geometry(f.resources.GeomSquareBack)

	f.resources.MatRock.Color(colorBack)
	for i := range f.listsBack {
		if f.listsBack[i] == nil {
			break
		}
		f.resources.MatRock.Color(f.colorsBack[i])
		for j := range f.listsBack[i] {
			r.Render(&f.listsBack[i][j])
		}
	}

	if len(f.listRock) > 0 {
		r.Geometry(f.resources.GeomRoundedCube)

		oldValue := -1
		for i := range f.listRock {
			if value := f.listRock[i].Value; value != oldValue {
				switch value {
				case 0:
					f.resources.MatRock.ClearChain()
				case 1:
					f.resources.MatRock.ChainTexture(f.resources.TexChain1)
				case 2:
					f.resources.MatRock.ChainTexture(f.resources.TexChain2)
				default:
					f.resources.MatRock.ChainTexture(f.resources.TexChain3)
				}
				oldValue = value
			}

			f.resources.MatRock.Color(f.listRock[i].Color)
			r.Render(&f.listRock[i].Model)
		}
		if oldValue != 0 {
			f.resources.MatRock.ClearChain()
		}
	}

	if len(f.listRuby) > 0 {
		r.Geometry(f.resources.GeomGem)

		oldValue := -1
		for i := range f.listRuby {
			if value := f.listRuby[i].Value; value != oldValue {
				switch value {
				case 0:
					f.resources.MatRock.ClearChain()
				case 1:
					f.resources.MatRock.ChainTexture(f.resources.TexChain1)
				case 2:
					f.resources.MatRock.ChainTexture(f.resources.TexChain2)
				default:
					f.resources.MatRock.ChainTexture(f.resources.TexChain3)
				}
				oldValue = value
			}

			f.resources.MatRock.Color(f.listRuby[i].Color)
			r.Render(&f.listRuby[i].Model)
		}
		if oldValue != 0 {
			f.resources.MatRock.ClearChain()
		}
	}

	if len(f.listFrame) > 0 {
		r.Geometry(f.resources.GeomFrame)
		f.resources.MatRock.Color(colorWall)
		for i := range f.listFrame {
			r.Render(&f.listFrame[i])
		}
	}

	if len(f.listShad) > 0 {
		r.Geometry(f.resources.GeomFrameThin)
		for i := range f.listShad {
			f.resources.MatRock.Color(f.listShad[i].Color)
			r.Render(&f.listShad[i].Model)
		}
	}

	if len(f.listIron) > 0 {
		r.Geometry(f.resources.GeomCube)
		r.Material(f.resources.MatIron)
		f.resources.MatIron.Lights(f.lights)
		for i := range f.listIron {
			r.Render(&f.listIron[i])
		}
	}

	if len(f.listLava) > 0 {
		r.Geometry(f.resources.GeomRoundedCube)
		r.Material(f.resources.MatLava)
		for i := range f.listLava {
			f.resources.MatLava.Color(f.listLava[i].Color)
			r.Render(&f.listLava[i].Model)
		}
	}

	if len(f.listAcid) > 0 {
		r.Geometry(f.resources.GeomRoundedCube)
		r.Material(f.resources.MatAcid)
		for i := range f.listAcid {
			f.resources.MatAcid.Color(f.listAcid[i].Color)
			r.Render(&f.listAcid[i].Model)
		}
	}

	if len(f.listWave) > 0 {
		r.Geometry(f.resources.GeomDie)
		r.Material(f.resources.MatWave)
		for i := range f.listWave {
			f.resources.MatWave.Color(f.listWave[i].Color)
			r.Render(&f.listWave[i].Model)
		}
	}

	if len(f.listBomb) > 0 {
		r.Geometry(f.resources.GeomStar6)
		r.Material(f.resources.MatRock)
		f.resources.MatRock.Color(colorBomb)
		transform :=
			mgl32.HomogRotate3DZ(float32(6.7*f.t + 0.2)).
				Mul4(mgl32.HomogRotate3DX(float32(5.1*f.t + 0.7)))
		for i := range f.listBomb {
			modelBomb := f.listBomb[i].Mul4(transform)
			r.Render(&modelBomb)
		}
	}

	if len(f.listGoal) > 0 {
		r.Geometry(f.resources.GeomStar8)
		r.Material(f.resources.MatColor)
		for i := range f.listGoal {
			f.resources.MatColor.Color(f.listGoal[i].Color)
			r.Render(&f.listGoal[i].Model)
		}
	}

	if len(f.listAmmo) > 0 {
		for i := range f.listAmmo {
			f.text.Rune(r, f.listAmmo[i].Model, f.listAmmo[i].Color, '0'+rune(f.listAmmo[i].Value))
		}
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
