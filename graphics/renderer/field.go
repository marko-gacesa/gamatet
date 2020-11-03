// Copyright (c) 2020 by Marko Gaćeša

package renderer

import (
	"gamatet/game/block"
	"gamatet/game/field"
	"gamatet/graphics/geometry"
	"gamatet/graphics/gtypes"
	"gamatet/graphics/material"
	"gamatet/graphics/rendercache"
	"github.com/go-gl/mathgl/mgl32"
	"math"
	"time"
)

type FieldRender struct {
	rock material.Rock
	lava material.Lava
	acid material.Acid
}

func NewFieldRender(
	rock material.Rock, lava material.Lava, acid material.Acid,
) FieldRender {
	return FieldRender{
		rock: rock,
		lava: lava,
		acid: acid,
	}
}

func (f FieldRender) Render(r *Renderer, model *mgl32.Mat4, renderInfo *field.RenderInfo) {
	if renderInfo == nil {
		return
	}

	t := time.Since(gtypes.Time).Seconds()

	// light intensities
	const (
		liPiece = 0.4
		liLava  = 1.5
		liAcid  = 1.5
		liRuby  = 3
	)

	listWall := rendercache.ModelPool.Get()
	listBack := rendercache.ModelPool.Get()
	listRock := rendercache.ModelColorPool.Get()
	listLava := rendercache.ModelColorPool.Get()
	listAcid := rendercache.ModelColorPool.Get()
	listRuby := rendercache.ModelPool.Get()
	listShad := rendercache.ModelColorPool.Get()
	lights := rendercache.PointLightPool.Get()
	defer func() {
		rendercache.ModelPool.Put(listWall)
		rendercache.ModelPool.Put(listBack)
		rendercache.ModelColorPool.Put(listRock)
		rendercache.ModelColorPool.Put(listLava)
		rendercache.ModelColorPool.Put(listAcid)
		rendercache.ModelPool.Put(listRuby)
		rendercache.ModelColorPool.Put(listShad)
		rendercache.PointLightPool.Put(lights)
	}()

	dx := float32(renderInfo.W) / 2
	dy := float32(renderInfo.H) / 2
	modelField := model.
		Mul4(mgl32.Translate3D(-dx, -dy, 0))

	// prepare the field frame

	for x := -1; x <= renderInfo.W; x++ {
		var m mgl32.Mat4

		m = modelField.Mul4(mgl32.Translate3D(float32(x), float32(-1), 0))
		listWall.Add(m)

		m = modelField.Mul4(mgl32.Translate3D(float32(x), float32(renderInfo.H), 0))
		listWall.Add(m)
	}

	for y := 0; y < renderInfo.H; y++ {
		var m mgl32.Mat4

		m = modelField.Mul4(mgl32.Translate3D(float32(-1), float32(y), 0))
		listWall.Add(m)

		m = modelField.Mul4(mgl32.Translate3D(float32(renderInfo.W), float32(y), 0))
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
		case block.TypeRuby:
			listRuby.Add(modelFieldBlock)
			lights.AddWithModel(modelFieldBlock, colorVector(block.Ruby.Color).Vec3(), liRuby)
		default:
			blockColor := colorVector(fb.Block.Color)
			color := mulColor(blockColor, aniColor)
			listRock.Add(modelFieldBlock, color)
		}
	}

	// prepare the pieces

	for _, p := range renderInfo.Pieces {
		if p.Empty {
			continue
		}

		aniMatrix, aniColor := animListUpdate(&p.Result)

		modelPiece := modelField.
			Mul4(mgl32.Translate3D(float32(p.X), float32(p.Y), 0)).
			Mul4(mgl32.Translate3D(float32(p.DimX)/2-0.5, -float32(p.DimY)/2+0.5, 0)).
			Mul4(aniMatrix).
			Mul4(mgl32.Translate3D(-float32(p.DimX)/2+0.5, float32(p.DimY)/2-0.5, 0))

		for _, pb := range p.Blocks {
			modelPieceBlock := modelPiece.
				Mul4(mgl32.Translate3D(float32(pb.X), float32(pb.Y), 0))

			switch pb.Block.Type {
			case block.TypeAcid:
				listAcid.Add(modelPieceBlock, aniColor)
				lights.AddWithModel(modelPieceBlock, colorVector(block.Acid.Color).Vec3(), liAcid)
			case block.TypeLava:
				listLava.Add(modelPieceBlock, aniColor)
				lights.AddWithModel(modelPieceBlock, colorVector(block.Acid.Color).Vec3(), liLava)
			default:
				blockColor := colorVector(pb.Block.Color)
				color := mulColor(blockColor, aniColor)
				listRock.Add(modelPieceBlock, color)
				lights.AddWithModel(modelPieceBlock, blockColor.Vec3(), liPiece)
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

	for pIdx, p := range renderInfo.Pieces {
		if p.Empty {
			continue
		}

		var dirY float32
		var modelPiece mgl32.Mat4
		switch pIdx {
		case 2:
			dirY = 1
			modelPiece = modelField.
				Mul4(mgl32.Translate3D(-2.5, 0, 0))
		case 1:
			dirY = 1
			modelPiece = modelField.
				Mul4(mgl32.Translate3D(float32(renderInfo.W)+2.5, 0, 0))
		case 0:
			dirY = -1
			modelPiece = modelField.
				Mul4(mgl32.Translate3D(-2.5, float32(renderInfo.H-2), 0))
		case 3:
			dirY = -1
			modelPiece = modelField.
				Mul4(mgl32.Translate3D(float32(renderInfo.W)+2.5, float32(renderInfo.H-2), 0))
		default:
			break
		}

		var y float32
		for i, nb := range p.NextBlocks {
			_ = i
			dim, centerX, centerY := center(nb)
			scale = 1 / (1 + 0.6*float32(i*i))
			dimScale := 0.5*scale + 0.3
			y += dirY * dimScale * dim / 2
			modelPieceN := modelPiece.
				Mul4(mgl32.Translate3D(0, y, 0)).
				Mul4(mgl32.Scale3D(dimScale, dimScale, dimScale)).
				Mul4(mgl32.HomogRotate3DX(-0.4)).
				Mul4(mgl32.HomogRotate3DZ(float32(math.Mod(2*t, 2*math.Pi))))
			y += dirY * dimScale * (dim/2 + 0.7)
			for _, pb := range nb {
				modelPieceBlock := modelPieceN.
					Mul4(mgl32.Translate3D(float32(pb.X)-centerX, float32(pb.Y)-centerY, 0))

				switch pb.Block.Type {
				case block.TypeAcid:
					blockColor := mgl32.Vec4{scale, scale, scale, 1}
					listAcid.Add(modelPieceBlock, blockColor)
				case block.TypeLava:
					blockColor := mgl32.Vec4{scale, scale, scale, 1}
					listLava.Add(modelPieceBlock, blockColor)
				default:
					blockColor := colorVector(pb.Block.Color).Mul(scale)
					listRock.Add(modelPieceBlock, blockColor)
				}
			}
		}
	}

	// render all

	r.Geometry(geometry.Cube)
	r.Material(f.rock)
	f.rock.Lights(lights)

	f.rock.Color(colorVector(block.Wall.Color))
	for i := range listWall {
		r.Render(&listWall[i])
	}

	r.Geometry(geometry.DentCube)

	f.rock.Color(colorVector(block.Wall.Color).Mul(0.6))
	for i := range listBack {
		r.Render(&listBack[i])
	}

	r.Geometry(geometry.RoundedCube)

	for i := range listRock {
		f.rock.Color(listRock[i].Color)
		r.Render(&listRock[i].Model)
	}

	if len(listShad) > 0 {
		r.Geometry(geometry.Frame)
		for i := range listShad {
			f.rock.Color(listShad[i].Color)
			r.Render(&listShad[i].Model)
		}
	}

	if len(listLava) > 0 {
		r.Geometry(geometry.RoundedCube)
		r.Material(f.lava)
		for i := range listLava {
			f.lava.Color(listLava[i].Color)
			r.Render(&listLava[i].Model)
		}
	}

	if len(listAcid) > 0 {
		r.Geometry(geometry.RoundedCube)
		r.Material(f.acid)
		for i := range listAcid {
			f.lava.Color(listAcid[i].Color)
			r.Render(&listAcid[i].Model)
		}
	}

	if len(listRuby) > 0 {
		r.Geometry(geometry.Gem)
		r.Material(f.lava) // TODO: change material
		f.lava.Color(colorVector(block.Ruby.Color))
		for i := range listRuby {
			r.Render(&listRuby[i])
		}
	}
}

func center(blocks []block.XYB) (float32, float32, float32) {
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
