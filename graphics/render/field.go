// Copyright (c) 2020-2024 by Marko Gaćeša

package render

import (
	"context"
	"gamatet/game/block"
	"gamatet/game/core"
	"gamatet/game/field"
	"gamatet/game/piece"
	"gamatet/graphics/render/rendercache"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"math"
	"slices"
	"strconv"
	"time"
)

var (
	colorWall = colorVector(block.Wall.Color)
	colorBack = colorVector(block.Wall.Color).Mul(0.6)
	colorLava = colorVector(block.Lava.Color)
	colorAcid = colorVector(block.Acid.Color)
	colorCurl = colorVector(block.Curl.Color)
	colorWave = colorVector(block.Wave.Color)
	colorBomb = colorVector(block.Bomb.Color)

	colorPlayerBack = []mgl32.Vec4{
		{0.0, 0.0, 0.4, 1.0},
		{0.4, 0.0, 0.0, 1.0},
		{0.0, 0.4, 0.0, 1.0},
		{0.4, 0.5, 0.0, 1.0},
	}
	colorPlayer = []mgl32.Vec4{
		{0.0, 0.5, 1.0, 0.8},
		{1.0, 0.2, 0.2, 0.8},
		{0.0, 1.0, 0.0, 0.8},
		{1.0, 1.0, 0.0, 0.8},
	}
	colorLabel = mgl32.Vec4{1, 1, 1, 0.5}
)

const widthPad = 5

var t0 = time.Now()

func GetExtendedContent(w, h int, infoPos []piece.DisplayPosition) (int, int) {
	w += 2 // left frame, right frame
	h += 1 // bottom frame
	var hasLeft, hasRight bool
	for _, p := range infoPos {
		hasLeft = hasLeft || p == piece.DisplayPositionTopLeft || p == piece.DisplayPositionBottomLeft
		hasRight = hasRight || p == piece.DisplayPositionTopRight || p == piece.DisplayPositionBottomRight
	}
	if hasLeft {
		w += widthPad
	}
	if hasRight {
		w += widthPad
	}
	return w, h
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
	listRock  rendercache.ModelColorValueList[int]
	listRuby  rendercache.ModelColorValueList[int]
	listLava  rendercache.ModelColorList
	listAcid  rendercache.ModelColorList
	listWave  rendercache.ModelColorList
	listGoal  rendercache.ModelColorList
	listShad  rendercache.ModelColorList
	listAmmo  rendercache.ModelColorValueList[int]
	lights    rendercache.PointLights
	listStr   rendercache.ModelColorValueList[string]
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
	f.listRock = rendercache.ModelColorIntPool.Get()
	f.listRuby = rendercache.ModelColorIntPool.Get()
	f.listLava = rendercache.ModelColorPool.Get()
	f.listAcid = rendercache.ModelColorPool.Get()
	f.listWave = rendercache.ModelColorPool.Get()
	f.listGoal = rendercache.ModelColorPool.Get()
	f.listShad = rendercache.ModelColorPool.Get()
	f.listAmmo = rendercache.ModelColorIntPool.Get()
	f.lights = rendercache.PointLightPool.Get()
	f.listStr = rendercache.ModelColorStringPool.Get()
}

func (f *Field) postRender() {
	for i, w := 0, f.w; i < w; i++ {
		rendercache.ModelPool.Put(f.listsBack[i])
	}

	rendercache.ModelPool.Put(f.listWall)
	rendercache.ModelPool.Put(f.listIron)
	rendercache.ModelPool.Put(f.listFrame)
	rendercache.ModelPool.Put(f.listBomb)
	rendercache.ModelColorIntPool.Put(f.listRock)
	rendercache.ModelColorIntPool.Put(f.listRuby)
	rendercache.ModelColorPool.Put(f.listLava)
	rendercache.ModelColorPool.Put(f.listAcid)
	rendercache.ModelColorPool.Put(f.listWave)
	rendercache.ModelColorPool.Put(f.listGoal)
	rendercache.ModelColorPool.Put(f.listShad)
	rendercache.ModelColorIntPool.Put(f.listAmmo)
	rendercache.PointLightPool.Put(f.lights)
	rendercache.ModelColorStringPool.Put(f.listStr)
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

	var infoPositions [field.MaxPieces]piece.DisplayPosition
	for i := range renderInfo.Pieces {
		infoPositions[i] = renderInfo.Pieces[i].Position
	}

	hasLeftPad := slices.Contains(infoPositions[:], piece.DisplayPositionTopLeft) || slices.Contains(infoPositions[:], piece.DisplayPositionBottomLeft)
	hasRightPad := slices.Contains(infoPositions[:], piece.DisplayPositionTopRight) || slices.Contains(infoPositions[:], piece.DisplayPositionBottomRight)

	contentWidth, contentHeight := renderInfo.W+2, renderInfo.H+1
	if hasLeftPad {
		contentWidth += widthPad
	}
	if hasRightPad {
		contentWidth += widthPad
	}

	var modelFrame, modelField mgl32.Mat4

	modelFrame = model.
		Mul4(mgl32.Translate3D(-float32(contentWidth)/2+0.5, -float32(contentHeight)/2+0.5, 0))

	if hasLeftPad {
		modelField = modelFrame.
			Mul4(mgl32.Translate3D(1+widthPad, 1, 0))
	} else {
		modelField = modelFrame.
			Mul4(mgl32.Translate3D(1, 1, 0))
	}

	pulse := float32(0.7 + 0.3*math.Sin(math.Mod(10*f.t, math.Pi)))

	// prepare the field frame

	for x := 0; x < contentWidth; x++ {
		m := modelFrame.Mul4(mgl32.Translate3D(float32(x), float32(0), 0))
		f.listWall.Add(m)
	}

	for y := 1; y < contentHeight; y++ {
		var m mgl32.Mat4

		m = modelFrame.Mul4(mgl32.Translate3D(float32(0), float32(y), 0))
		f.listWall.Add(m)

		if hasLeftPad {
			for i := 0; i < widthPad; i++ {
				m = modelFrame.Mul4(mgl32.Translate3D(float32(1+i), float32(y), 0))
				f.listWall.Add(m)
			}
		}
		if hasRightPad {
			for i := 0; i < widthPad; i++ {
				m = modelFrame.Mul4(mgl32.Translate3D(float32(contentWidth-2-i), float32(y), 0))
				f.listWall.Add(m)
			}
		}

		m = modelFrame.Mul4(mgl32.Translate3D(float32(contentWidth-1), float32(y), 0))
		f.listWall.Add(m)
	}

	// prepare the background

	for x := 0; x < renderInfo.W; x++ {
		colorCol := colorBack
		for pIdx := range renderInfo.Pieces {
			p := &renderInfo.Pieces[pIdx]

			if p.Position == piece.DisplayPositionOff {
				continue
			}

			within := p.IsLimited && x >= p.Limits.Min && x <= p.Limits.Max
			if within {
				colorCol = colorCol.Add(colorPlayerBack[pIdx])
			}

			shadow := !p.PieceEmpty && x >= renderInfo.Pieces[pIdx].Shadow.ColL && x < renderInfo.Pieces[pIdx].Shadow.ColR
			if shadow {
				colorCol = colorCol.Mul(1.1)
			}
		}

		f.colorsBack[x] = colorCol

		for y := 0; y < renderInfo.H; y++ {
			m := modelField.Mul4(mgl32.Translate3D(float32(x), float32(y), float32(-1)))
			f.listsBack[x].Add(m)
		}
	}

	// pause

	if renderInfo.Paused {
		const text = "PAUSE"

		w, h := f.text.Dim(text)
		w, h = w*pulse, h*pulse

		modelPause := modelField.
			Mul4(mgl32.Translate3D(float32(renderInfo.W)/2.0-w/2-0.5, float32(renderInfo.H)/2-0.5, 0)).
			Mul4(mgl32.Scale3D(pulse, pulse, 1))

		f.listStr.Add(modelPause, mgl32.Vec4{1, 1, 1, 1}, text)
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
		if p.PieceEmpty {
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

	for _, p := range renderInfo.Pieces {
		if p.PieceEmpty || !p.DrawShadow {
			continue
		}

		for _, pb := range p.Shadow.Blocks {
			modelPieceShadowBlock := modelField.
				Mul4(mgl32.Translate3D(float32(pb.X), float32(pb.Y), 0)).
				Mul4(mgl32.Scale3D(pulse, pulse, pulse))
			blockColor := colorVector(pb.Block.Color).Mul(0.7)
			f.listShad.Add(modelPieceShadowBlock, blockColor)
		}
	}

	// prepare player info strings and next pieces

	for idx, p := range renderInfo.Pieces {
		var modelInfo mgl32.Mat4
		var hDir float32

		colorText := colorPlayer[idx]

		const edgeOffset = 0.75

		switch p.Position {
		case piece.DisplayPositionTopLeft:
			modelInfo = modelField.Mul4(mgl32.Translate3D(-(1+widthPad)-0.5, float32(renderInfo.H)-edgeOffset-0.5, 0.5))
			hDir = -1
		case piece.DisplayPositionTopRight:
			modelInfo = modelField.Mul4(mgl32.Translate3D(float32(renderInfo.W)-0.5, float32(renderInfo.H)-edgeOffset-0.5, 0.5))
			hDir = -1
		case piece.DisplayPositionBottomLeft:
			modelInfo = modelField.Mul4(mgl32.Translate3D(-(1+widthPad)-0.5, edgeOffset-1-0.5, 0.5))
			hDir = 1
		case piece.DisplayPositionBottomRight:
			modelInfo = modelField.Mul4(mgl32.Translate3D(float32(renderInfo.W)-0.5, edgeOffset-1-0.5, 0.5))
			hDir = 1
		default:
			continue
		}

		//f.listStr.Add(modelInfo, mgl32.Vec4{1, 1, 1, 1}, "X")
		//f.listRock.Add(
		//	modelInfo.
		//		Mul4(mgl32.Translate3D(0, 0, 0.5)).
		//		Mul4(mgl32.Scale3D(pulse*0.3, pulse*0.2, pulse*0.3)),
		//	mgl32.Vec4{1, 1, 1, 1}, 0)
		//f.listRock.Add(
		//	modelInfo.
		//		Mul4(mgl32.Translate3D(widthPad+1, 0, 0.5)).
		//		Mul4(mgl32.Scale3D(pulse*0.3, pulse*0.2, pulse*0.3)),
		//	mgl32.Vec4{1, 1, 1, 1}, 0)

		f.printValue(&modelInfo, colorLabel, colorText, "PLAYER", p.PieceTextData.Name, hDir)
		f.printValue(&modelInfo, colorLabel, colorText, "SCORE", p.PieceTextData.Score, hDir)
		f.printValue(&modelInfo, colorLabel, colorText, "PIECE", p.PieceTextData.PieceNum, hDir)
		f.printValue(&modelInfo, colorLabel, colorText, "STATE", strconv.Itoa(int(p.State)), hDir)

		modelInfo = modelInfo.Mul4(mgl32.Translate3D(0, 0.5*hDir, 0))
		f.printText(&modelInfo, colorLabel, "NEXT", hDir)

		modelNextBlocks := modelInfo.
			Mul4(mgl32.Translate3D((1.0+widthPad)/2, 0.3*hDir, 0.5))

		var y float32
		for i, nb := range p.NextBlocks {
			dim, centerX, centerY := barycenter(nb)
			dimScale := 0.3/(float32(3*i)+1.0) + 0.2
			y += hDir * dimScale * dim / 2
			modelPieceN := modelNextBlocks.
				Mul4(mgl32.Translate3D(0, y, 0)).
				Mul4(mgl32.Scale3D(dimScale, dimScale, dimScale)).
				Mul4(mgl32.HomogRotate3DX(-0.4)).
				Mul4(mgl32.HomogRotate3DZ(float32(math.Mod(f.t, 2*math.Pi))))
			y += hDir * dimScale * (dim/2 + 0.7)

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

	for i := range f.listAmmo {
		f.text.Rune(r, f.listAmmo[i].Model, f.listAmmo[i].Color, '0'+rune(f.listAmmo[i].Value))
	}

	gl.Disable(gl.DEPTH_TEST)

	for i := range f.listStr {
		f.text.String(r, f.listStr[i].Model, f.listStr[i].Color, f.listStr[i].Value)
	}

	gl.Enable(gl.DEPTH_TEST)
}

func (f *Field) printValue(modelInfo *mgl32.Mat4, colorLabel, colorText mgl32.Vec4, title, value string, hDir float32) {
	if value == "" {
		return
	}

	const scaleTitle = 0.5
	const scaleValue = 0.8
	const padding = 0.2

	var wt, ht, wv, hv float32
	wt, ht = f.text.Dim(title)
	wv, hv = f.text.Dim(value)
	wt, ht = scaleTitle*wt, scaleTitle*ht
	wv, hv = scaleValue*wv, scaleValue*hv

	// so that the value is always below the title
	var yt, yv float32
	if hDir > 0 {
		yt = hDir * (hv + ht*0.5)
		yv = hDir * hv * 0.5
	} else {
		yt = hDir * ht * 0.5
		yv = hDir * (ht + hv*0.5)
	}

	modelTitle := modelInfo.
		Mul4(mgl32.Translate3D((1.0+widthPad)/2.0-wt/2, yt, 0)).
		Mul4(mgl32.Scale3D(scaleTitle, scaleTitle, 1.0))
	f.listStr.Add(modelTitle, colorLabel, title)

	modelValue := modelInfo.
		Mul4(mgl32.Translate3D((1.0+widthPad)/2.0-wv/2, yv, 0)).
		Mul4(mgl32.Scale3D(scaleValue, scaleValue, 1.0))
	f.listStr.Add(modelValue, colorText, value)

	*modelInfo = modelInfo.Mul4(mgl32.Translate3D(0, hDir*(ht+hv+padding), 0))
}

func (f *Field) printText(modelInfo *mgl32.Mat4, colorText mgl32.Vec4, s string, hDir float32) {
	const scale = 0.6
	const padding = 0.2

	var w, h float32
	w, h = f.text.Dim(s)
	w, h = scale*w, scale*h

	y := hDir * h * 0.5

	modelValue := modelInfo.
		Mul4(mgl32.Translate3D((1.0+widthPad)/2.0-w/2, y, 0)).
		Mul4(mgl32.Scale3D(scale, scale, 1.0))
	f.listStr.Add(modelValue, colorText, s)

	*modelInfo = modelInfo.Mul4(mgl32.Translate3D(0, hDir*(h+padding), 0))
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
