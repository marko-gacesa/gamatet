// Copyright (c) 2020-2025 by Marko Gaćeša

package render

import (
	"math"
	"slices"
	"sync"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/game/core"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/piece"
	"github.com/marko-gacesa/gamatet/game/setup"
	"github.com/marko-gacesa/gamatet/graphics/render/rendercache"
)

var (
	colorWall = colorVector(block.Wall.Color)
	colorBack = colorVector(block.Wall.Color).Mul(0.6)
	colorLava = colorVector(block.Lava.Color)
	colorAcid = colorVector(block.Acid.Color)
	colorCurl = colorVector(block.Curl.Color)
	colorWave = colorVector(block.Wave.Color)
	colorBomb = colorVector(block.Bomb.Color)

	colorText  = mgl32.Vec4{1, 1, 1, 0.8}
	colorLabel = mgl32.Vec4{1, 1, 1, 0.5}

	colorPlayer     = [setup.MaxPlayers]mgl32.Vec4{}
	colorPlayerBack = [setup.MaxPlayers]mgl32.Vec4{}
)

func init() {
	for i := range setup.MaxPlayers {
		colorPlayer[i] = mgl32.Vec3(setup.ColorRGB[i][:]).Vec4(1.0)
		colorPlayerBack[i] = mgl32.Vec3(setup.ColorBackRGB[i][:]).Vec4(1.0)
	}
}

const (
	scaleLabel = 0.6
	scaleText  = 0.8

	paddingText = 0.2
)

const sidePanelBlockWidth = 5

var t0 = time.Now()

func GetExtendedContent(w, h int, infoPos [field.MaxPieces]DisplayPosition) (int, int) {
	w += 2 // left frame, right frame
	h += 2 // top, bottom frame
	var hasLeft, hasRight bool
	for _, p := range infoPos {
		hasLeft = hasLeft || p == DisplayPositionTopLeft || p == DisplayPositionBottomLeft
		hasRight = hasRight || p == DisplayPositionTopRight || p == DisplayPositionBottomRight
	}
	if hasLeft {
		w += sidePanelBlockWidth
	}
	if hasRight {
		w += sidePanelBlockWidth
	}
	return w, h
}

type FieldStrings struct {
	TitlePanel struct {
		Score  string
		Blocks string
	}
	SidePanel struct {
		Player string
		Score  string
		Piece  string
		Next   string
	}
	Message struct {
		GameOver   string
		Victory    string
		Defeat     string
		Pause      string
		Suspended  string
		ServerLost string
	}
}

type Field struct {
	model     mgl32.Mat4
	resources *FieldResources
	text      *Text

	str FieldStrings

	renderRequesterFieldIdx int
	renderRequester         core.RenderRequester
	preferredSide           PreferredSide

	renderInfo   *field.RenderInfo
	renderInfoWG sync.WaitGroup

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
	str FieldStrings,
	renderRequesterFieldIdx int,
	renderRequester core.RenderRequester,
	preferredSide PreferredSide,
) *Field {
	return &Field{
		model:                   model,
		resources:               resources,
		text:                    text,
		str:                     str,
		renderRequesterFieldIdx: renderRequesterFieldIdx,
		renderRequester:         renderRequester,
		preferredSide:           preferredSide,
	}
}

func (f *Field) Prepare(now time.Time) {
	renderInfoCh := make(chan *field.RenderInfo)
	f.renderRequester.RenderRequest(f.renderRequesterFieldIdx, now, renderInfoCh)
	f.renderInfoWG.Add(1)
	go func() {
		defer f.renderInfoWG.Done()

		renderInfo := <-renderInfoCh
		if renderInfo == nil {
			return
		}

		f.preRender(renderInfo, now)
		f.prepareModels(renderInfo)
		field.ReturnRenderInfo(f.renderInfo)
	}()
}

func (f *Field) Render(r *Renderer) {
	f.renderInfoWG.Wait()
	f.renderAll(r)
	f.postRender()
}

func (f *Field) preRender(renderInfo *field.RenderInfo, now time.Time) {
	if renderInfo == nil {
		return
	}

	f.t = now.Sub(t0).Seconds()

	w := renderInfo.W
	for i := range w {
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
	w := f.w
	for i := range w {
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
	// light intensities
	const (
		lightIntLava    = 1.2
		lightIntAcid    = 1.2
		lightIntWave    = 1.5
		lightIntGoal    = 3
		lightPowShooter = 1.5
	)

	infoPositions := f.preferredSide.PieceCorners(renderInfo.PieceCount)

	hasLeftPad := slices.Contains(infoPositions[:], DisplayPositionTopLeft) || slices.Contains(infoPositions[:], DisplayPositionBottomLeft)
	hasRightPad := slices.Contains(infoPositions[:], DisplayPositionTopRight) || slices.Contains(infoPositions[:], DisplayPositionBottomRight)

	contentWidth, contentHeight := renderInfo.W+2, renderInfo.H+2
	if hasLeftPad {
		contentWidth += sidePanelBlockWidth
	}
	if hasRightPad {
		contentWidth += sidePanelBlockWidth
	}

	aniMatrixField, _ := animListUpdate(&renderInfo.Result)
	model := f.model.Mul4(aniMatrixField)

	var modelFrame, modelField mgl32.Mat4

	modelFrame = model.
		Mul4(mgl32.Translate3D(-float32(contentWidth)/2+0.5, -float32(contentHeight)/2+0.5, 0))

	if hasLeftPad {
		modelField = modelFrame.
			Mul4(mgl32.Translate3D(1+sidePanelBlockWidth, 1, 0))
	} else {
		modelField = modelFrame.
			Mul4(mgl32.Translate3D(1, 1, 0))
	}

	pulse := float32(0.7 + 0.3*math.Sin(math.Mod(10*f.t, math.Pi)))

	// prepare the field frame

	for x := range contentWidth {
		f.listWall.Add(modelFrame.Mul4(mgl32.Translate3D(float32(x), float32(0), 0)))
		f.listWall.Add(modelFrame.Mul4(mgl32.Translate3D(float32(x), float32(contentHeight-1), 0)))
	}

	modelTitleLeft := modelField.Mul4(mgl32.Translate3D(-0.5, float32(contentHeight-2), 0.5))
	modelTitleRight := modelTitleLeft.Mul4(mgl32.Translate3D(float32(f.w), 0, 0))
	//d := f.printText(
	//	modelTitleLeft,
	//	colorLabel,
	//	scaleLabel,
	//	f.str.TitlePanel.Score)
	//f.printText(
	//	modelTitleLeft.Mul4(mgl32.Translate3D(d+0.25, 0, 0)),
	//	colorText,
	//	scaleText,
	//	"00234234")
	d := f.printTextRight(
		modelTitleRight,
		colorText,
		scaleText,
		renderInfo.TextData.BlocksRemoved)
	f.printTextRight(
		modelTitleRight.Mul4(mgl32.Translate3D(-d-0.25, 0, 0)),
		colorLabel,
		scaleLabel,
		f.str.TitlePanel.Blocks)

	for y := 1; y < contentHeight-1; y++ {
		var m mgl32.Mat4

		m = modelFrame.Mul4(mgl32.Translate3D(float32(0), float32(y), 0))
		f.listWall.Add(m)

		if hasLeftPad {
			for i := range sidePanelBlockWidth {
				m = modelFrame.Mul4(mgl32.Translate3D(float32(1+i), float32(y), 0))
				f.listWall.Add(m)
			}
		}
		if hasRightPad {
			for i := range sidePanelBlockWidth {
				m = modelFrame.Mul4(mgl32.Translate3D(float32(contentWidth-2-i), float32(y), 0))
				f.listWall.Add(m)
			}
		}

		m = modelFrame.Mul4(mgl32.Translate3D(float32(contentWidth-1), float32(y), 0))
		f.listWall.Add(m)
	}

	// prepare the background

	for x := range renderInfo.W {
		colorCol := colorBack
		for pIdx := range renderInfo.Pieces {
			p := &renderInfo.Pieces[pIdx]

			if infoPositions[pIdx] == DisplayPositionOff {
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

		for y := range renderInfo.H {
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

	for idx := range renderInfo.Pieces {
		var modelInfo mgl32.Mat4
		var hDir float32

		const edgeOffset = 0.75

		switch infoPositions[idx] {
		case DisplayPositionTopLeft:
			modelInfo = modelField.Mul4(mgl32.Translate3D(-(1+sidePanelBlockWidth)-0.5, float32(renderInfo.H)-edgeOffset+1-0.5, 0.5))
			hDir = -1
		case DisplayPositionTopRight:
			modelInfo = modelField.Mul4(mgl32.Translate3D(float32(renderInfo.W)-0.5, float32(renderInfo.H)-edgeOffset+1-0.5, 0.5))
			hDir = -1
		case DisplayPositionBottomLeft:
			modelInfo = modelField.Mul4(mgl32.Translate3D(-(1+sidePanelBlockWidth)-0.5, edgeOffset-1-0.5, 0.5))
			hDir = 1
		case DisplayPositionBottomRight:
			modelInfo = modelField.Mul4(mgl32.Translate3D(float32(renderInfo.W)-0.5, edgeOffset-1-0.5, 0.5))
			hDir = 1
		default:
			continue
		}

		// Render text position helpers
		//f.listStr.Add(modelInfo, mgl32.Vec4{1, 1, 1, 1}, "X")
		//f.listRock.Add(
		//	modelInfo.
		//		Mul4(mgl32.Translate3D(0, 0, 0.5)).
		//		Mul4(mgl32.Scale3D(pulse*0.3, pulse*0.2, pulse*0.3)),
		//	mgl32.Vec4{1, 1, 1, 1}, 0)
		//f.listRock.Add(
		//	modelInfo.
		//		Mul4(mgl32.Translate3D(sidePanelBlockWidth+1, 0, 0.5)).
		//		Mul4(mgl32.Scale3D(pulse*0.3, pulse*0.2, pulse*0.3)),
		//	mgl32.Vec4{1, 1, 1, 1}, 0)

		p := &renderInfo.Pieces[idx]
		colorPlayerText := colorPlayer[p.PlayerIndex]

		f.printLabelAndValue(&modelInfo, colorLabel, colorPlayerText, f.str.SidePanel.Player, p.PieceTextData.Name, hDir)
		f.printLabelAndValue(&modelInfo, colorLabel, colorPlayerText, f.str.SidePanel.Score, p.PieceTextData.Score, hDir)
		f.printLabelAndValue(&modelInfo, colorLabel, colorPlayerText, f.str.SidePanel.Piece, p.PieceTextData.PieceNum, hDir)

		if len(p.NextPieces[0].Blocks) == 0 {
			continue
		}

		modelInfo = modelInfo.Mul4(mgl32.Translate3D(0, 0.5*hDir, 0))
		f.printLabel(&modelInfo, colorLabel, f.str.SidePanel.Next, hDir)

		modelNextBlocks := modelInfo.
			Mul4(mgl32.Translate3D((1.0+sidePanelBlockWidth)/2, 0.3*hDir, 0.5))

		var y float32
		for i, np := range p.NextPieces {
			dim, centerX, centerY := barycenter(np.Blocks)
			dimScale := 0.3/(float32(3*i)+1.0) + 0.2

			y += hDir * dimScale * dim / 2

			modelPieceN := modelNextBlocks.
				Mul4(mgl32.Translate3D(0, y, 0)).
				Mul4(mgl32.Scale3D(dimScale, dimScale, dimScale)).
				Mul4(mgl32.HomogRotate3DX(-0.4))

			switch np.Type {
			case piece.TypeFlipV:
				modelPieceN = modelPieceN.Mul4(mgl32.HomogRotate3DX(float32(math.Mod(2*f.t, 2*math.Pi))))
			case piece.TypeFlipH:
				modelPieceN = modelPieceN.Mul4(mgl32.HomogRotate3DY(float32(math.Mod(2*f.t, 2*math.Pi))))
			case piece.TypeRotation:
				t := f.t
				if p.DirectionCW {
					t = -t
				}
				modelPieceN = modelPieceN.Mul4(mgl32.HomogRotate3DZ(float32(math.Mod(t, 2*math.Pi))))
			case piece.TypeShooter:
				t := math.Sin(10 * f.t)
				t = 0.25*t*t + 0.75
				s := float32(t)
				modelPieceN = modelPieceN.Mul4(mgl32.Scale3D(s, s, s))
			}

			y += hDir * dimScale * (dim/2 + 0.7)

			for _, pb := range np.Blocks {
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

	// text

	var message string

	switch renderInfo.Mode {
	case field.ModeNormal:
	case field.ModeGameOver:
		message = f.str.Message.GameOver
	case field.ModeVictory:
		message = f.str.Message.Victory
	case field.ModeDefeat:
		message = f.str.Message.Defeat
	case field.ModePause:
		message = f.str.Message.Pause
	case field.ModeSuspended:
		message = f.str.Message.Suspended
	case field.ModeServerLost:
		message = f.str.Message.ServerLost
	}

	if message != "" {
		w, h := f.text.Dim(message)

		modelMessage := modelField.
			Mul4(mgl32.Translate3D((float32(renderInfo.W-1)-w*pulse)/2.0, (float32(renderInfo.H-2)+h*pulse)/2.0, 0)).
			Mul4(mgl32.Scale3D(pulse, pulse, 1))
		f.listStr.Add(modelMessage, mgl32.Vec4{1, 1, 1, 1}, message)
	}

	if (renderInfo.Mode == field.ModePause || renderInfo.Mode == field.ModeSuspended) && renderInfo.TextData.Latencies != "" {
		const scale = 0.5
		_, h := f.text.Dim(renderInfo.TextData.Latencies)
		modelLatencies := modelField.
			Mul4(mgl32.Translate3D(0, h*scale, 0)).
			Mul4(mgl32.Scale3D(scale, scale, 1))
		f.listStr.Add(modelLatencies, mgl32.Vec4{0.6, 0.6, 0.6, 1}, renderInfo.TextData.Latencies)
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

func (f *Field) printLabelAndValue(modelInfo *mgl32.Mat4, colorLabel, colorText mgl32.Vec4, label, value string, hDir float32) {
	if value == "" {
		return
	}

	const scaleTitle = 0.5
	const scaleValue = 0.8
	const padding = 0.2

	var wt, ht, wv, hv float32
	wt, ht = f.text.Dim(label)
	wv, hv = f.text.Dim(value)
	wt, ht = scaleTitle*wt, scaleTitle*ht
	wv, hv = scaleValue*wv, scaleValue*hv

	// so that the value is always below the label
	var yt, yv float32
	if hDir > 0 {
		yt = hDir * (hv + ht*0.5)
		yv = hDir * hv * 0.5
	} else {
		yt = hDir * ht * 0.5
		yv = hDir * (ht + hv*0.5)
	}

	modelTitle := modelInfo.
		Mul4(mgl32.Translate3D((1.0+sidePanelBlockWidth)/2.0-wt/2, yt, 0)).
		Mul4(mgl32.Scale3D(scaleTitle, scaleTitle, 1.0))
	f.listStr.Add(modelTitle, colorLabel, label)

	modelValue := modelInfo.
		Mul4(mgl32.Translate3D((1.0+sidePanelBlockWidth)/2.0-wv/2, yv, 0)).
		Mul4(mgl32.Scale3D(scaleValue, scaleValue, 1.0))
	f.listStr.Add(modelValue, colorText, value)

	*modelInfo = modelInfo.Mul4(mgl32.Translate3D(0, hDir*(ht+hv+padding), 0))
}

func (f *Field) printLabel(modelInfo *mgl32.Mat4, colorText mgl32.Vec4, s string, hDir float32) {
	var w, h float32
	w, h = f.text.Dim(s)
	w, h = scaleLabel*w, scaleLabel*h

	y := hDir * h * 0.5

	modelValue := modelInfo.
		Mul4(mgl32.Translate3D((1.0+sidePanelBlockWidth)/2.0-w/2, y, 0)).
		Mul4(mgl32.Scale3D(scaleLabel, scaleLabel, 1.0))
	f.listStr.Add(modelValue, colorText, s)

	*modelInfo = modelInfo.Mul4(mgl32.Translate3D(0, hDir*(h+paddingText), 0))
}

func (f *Field) printText(m mgl32.Mat4, colorText mgl32.Vec4, scaleText float32, text string) float32 {
	w, _ := f.text.Dim(text)
	w *= scaleText
	f.listStr.Add(m.Mul4(mgl32.Scale3D(scaleText, scaleText, 1.0)), colorText, text)
	return w
}

func (f *Field) printTextRight(m mgl32.Mat4, colorText mgl32.Vec4, scaleText float32, text string) float32 {
	w, _ := f.text.Dim(text)
	w *= scaleText
	m = m.Mul4(mgl32.Translate3D(-w, 0, 0))
	f.listStr.Add(m.Mul4(mgl32.Scale3D(scaleText, scaleText, 1.0)), colorText, text)
	return w
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
