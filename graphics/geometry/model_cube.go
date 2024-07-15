// Copyright (c) 2020-2024 by Marko Gaćeša

package geometry

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

func CubeSideSimple(model mgl32.Mat3, v *[]blockVertex) {
	halfSide(model, func(model, texture mgl32.Mat3) {
		n := model.Mul3x1(mgl32.Vec3{0.0, 0.0, 1.0})

		p0 := model.Mul3x1(mgl32.Vec3{-0.5, -0.5, 0.5})
		p1 := model.Mul3x1(mgl32.Vec3{-0.5, 0.5, 0.5})
		p2 := model.Mul3x1(mgl32.Vec3{0.5, 0.5, 0.5})

		v0 := gen(p0, n, texture.Mul3x1(mgl32.Vec3{0, 1, 1}))
		v1 := gen(p1, n, texture.Mul3x1(mgl32.Vec3{0, 0, 1}))
		v2 := gen(p2, n, texture.Mul3x1(mgl32.Vec3{1, 0, 1}))

		*v = append(*v, v0, v1, v2)
	})
}

func CubeSideDent(model mgl32.Mat3, v *[]blockVertex) {
	quarterSide(model, func(model, texture mgl32.Mat3) {
		nn := model.Mul3x1(mgl32.Vec3{0.0, 0.0, 1.0}).Normalize()
		n0 := model.Mul3x1(mgl32.Vec3{0.5, -0.5, 1.0}).Normalize()
		n1 := model.Mul3x1(mgl32.Vec3{0.5, 0.5, 1.0}).Normalize()
		nc := model.Mul3x1(mgl32.Vec3{-0.0, 0.0, 1.0}).Normalize()

		pOuter0 := model.Mul3x1(mgl32.Vec3{-0.5, -0.5, 0.5})
		pOuter1 := model.Mul3x1(mgl32.Vec3{-0.5, 0.5, 0.5})
		pInner0 := model.Mul3x1(mgl32.Vec3{-0.3, -0.3, 0.45})
		pInner1 := model.Mul3x1(mgl32.Vec3{-0.3, 0.3, 0.45})
		pCenter := model.Mul3x1(mgl32.Vec3{0.0, 0.0, 0.45})

		uvOuter0 := texture.Mul3x1(mgl32.Vec3{0, 1, 1})
		uvOuter1 := texture.Mul3x1(mgl32.Vec3{0, 0, 1})
		uvInner0 := texture.Mul3x1(mgl32.Vec3{0.2, 0.8, 1})
		uvInner1 := texture.Mul3x1(mgl32.Vec3{0.2, 0.2, 1})
		uvCenter := texture.Mul3x1(mgl32.Vec3{0.5, 0.5, 1})

		vOuter0 := gen(pOuter0, n0, uvOuter0)
		vInner0 := gen(pInner0, nn, uvInner0)
		vOuter1 := gen(pOuter1, n1, uvOuter1)
		vInner1 := gen(pInner1, nn, uvInner1)
		vCenter := gen(pCenter, nc, uvCenter)

		// 3 triangles
		*v = append(*v, vOuter0, vOuter1, vInner1)
		*v = append(*v, vInner1, vInner0, vOuter0)
		*v = append(*v, vInner0, vInner1, vCenter)
	})
}

func CubeSideFrame(size, thickness float32) func(mgl32.Mat3, *[]blockVertex) {
	halfA := size / 2 // halfA is the cube's side length (max is 0.5 for the full sized cube)
	d := thickness
	return func(side mgl32.Mat3, v *[]blockVertex) {
		quarterSide(side, func(model, texture mgl32.Mat3) {
			nn := model.Mul3x1(mgl32.Vec3{0.0, 0.0, 1.0})
			nf := model.Mul3x1(mgl32.Vec3{1.0, 0.0, 0.0})

			pOuter0 := model.Mul3x1(mgl32.Vec3{-halfA, -halfA, halfA})
			pOuter1 := model.Mul3x1(mgl32.Vec3{-halfA, halfA, halfA})
			pInner0 := model.Mul3x1(mgl32.Vec3{-halfA + d, -halfA + d, halfA})
			pInner1 := model.Mul3x1(mgl32.Vec3{-halfA + d, halfA - d, halfA})
			pLower0 := model.Mul3x1(mgl32.Vec3{-halfA + d, -halfA + d, halfA - d})
			pLower1 := model.Mul3x1(mgl32.Vec3{-halfA + d, halfA - d, halfA - d})

			uvOuter0 := texture.Mul3x1(mgl32.Vec3{0, 1, 1})
			uvOuter1 := texture.Mul3x1(mgl32.Vec3{0, 0, 1})
			uvInner0 := texture.Mul3x1(mgl32.Vec3{d, 1 - d, 1})
			uvInner1 := texture.Mul3x1(mgl32.Vec3{d, d, 1})
			uvLower0 := texture.Mul3x1(mgl32.Vec3{d + d, 1 - d, 1})
			uvLower1 := texture.Mul3x1(mgl32.Vec3{d + d, d, 1})

			vOuter0 := gen(pOuter0, nn, uvOuter0)
			vInner0 := gen(pInner0, nn, uvInner0)
			vOuter1 := gen(pOuter1, nn, uvOuter1)
			vInner1 := gen(pInner1, nn, uvInner1)

			// 2 triangles
			*v = append(*v, vOuter0, vOuter1, vInner1)
			*v = append(*v, vInner1, vInner0, vOuter0)

			vInner0.setN(nf)
			vLower0 := gen(pLower0, nf, uvLower0)
			vInner1.setN(nf)
			vLower1 := gen(pLower1, nf, uvLower1)

			// 2 triangles
			*v = append(*v, vInner0, vInner1, vLower1)
			*v = append(*v, vLower1, vLower0, vInner0)
		})
	}
}

func CubeSideRounded(side mgl32.Mat3, v *[]blockVertex) {
	const rund = 0.16         // 0.16
	const alpha = math.Pi / 6 // pi/6

	sq3 := math.Sqrt(3)
	w := float32(rund * (1 - (1 / sq3)))

	tg := float32(math.Tan(alpha))
	d := rund * tg / (1 + tg)
	l := rund / (1 + tg) // l+d = w

	quarterSide(side, func(model, texture mgl32.Mat3) {
		nEdge0 := model.Mul3x1(mgl32.Vec3{-1, 0, 1}).Normalize()
		nEdge1 := model.Mul3x1(mgl32.Vec3{-1, 0, 1}).Normalize()
		nCent := model.Mul3x1(mgl32.Vec3{0, 0, 1}).Normalize()
		nVert0 := model.Mul3x1(mgl32.Vec3{-1, -1, 1}).Normalize()
		nVert1 := model.Mul3x1(mgl32.Vec3{-1, 1, 1}).Normalize()

		pEdge0 := model.Mul3x1(mgl32.Vec3{-0.5 + d, -0.5 + rund, 0.5 - d})
		pEdge1 := model.Mul3x1(mgl32.Vec3{-0.5 + d, 0.5 - rund, 0.5 - d})
		pTop0 := model.Mul3x1(mgl32.Vec3{-0.5 + rund, -0.5 + rund, 0.5})
		pTop1 := model.Mul3x1(mgl32.Vec3{-0.5 + rund, 0.5 - rund, 0.5})
		pCent := model.Mul3x1(mgl32.Vec3{0.0, 0.0, 0.5})
		pVert0 := model.Mul3x1(mgl32.Vec3{-0.5 + w, -0.5 + w, 0.5 - w})
		pVert1 := model.Mul3x1(mgl32.Vec3{-0.5 + w, 0.5 - w, 0.5 - w})

		uvEdge0 := texture.Mul3x1(mgl32.Vec3{0, 1 - rund, 1})
		uvEdge1 := texture.Mul3x1(mgl32.Vec3{0, rund, 1})
		uvTop0 := texture.Mul3x1(mgl32.Vec3{0 + rund, 1 - rund, 1})
		uvTop1 := texture.Mul3x1(mgl32.Vec3{0 + rund, rund, 1})
		uvCent := texture.Mul3x1(mgl32.Vec3{0.5, 0.5, 1})
		uvVert0 := texture.Mul3x1(mgl32.Vec3{0.0, 1.0, 1})
		uvVert1 := texture.Mul3x1(mgl32.Vec3{0.0, 0.0, 1})

		vEdge0 := gen(pEdge0, nEdge0, uvEdge0)
		vEdge1 := gen(pEdge1, nEdge1, uvEdge1)
		vTop0 := gen(pTop0, nCent, uvTop0)
		vTop1 := gen(pTop1, nCent, uvTop1)
		vCent := gen(pCent, nCent, uvCent)
		vVert0 := gen(pVert0, nVert0, uvVert0)
		vVert1 := gen(pVert1, nVert1, uvVert1)

		*v = append(*v, vEdge0, vEdge1, vTop1)
		*v = append(*v, vTop1, vTop0, vEdge0)
		*v = append(*v, vTop0, vTop1, vCent)

		if d >= rund || l <= 0.000001 {
			return
		}

		*v = append(*v, vVert0, vEdge0, vTop0)
		*v = append(*v, vVert1, vTop1, vEdge1)
	})
}

func CubeSideTruncated(side mgl32.Mat3, v *[]blockVertex) {
	//const edge = 0.2 // pretty
	const edge = 1 / (2 + math.Sqrt2) // symmetric (Rhombicuboctahedron)
	const vert = edge * 2 / 3

	quarterSide(side, func(model, texture mgl32.Mat3) {
		nSlope := model.Mul3x1(mgl32.Vec3{-1, 0, 1}).Normalize()
		nTop := model.Mul3x1(mgl32.Vec3{0, 0, 1})
		nVert0 := model.Mul3x1(mgl32.Vec3{-1, -1, 1}).Normalize()
		nVert1 := model.Mul3x1(mgl32.Vec3{-1, 1, 1}).Normalize()

		pEdge0 := model.Mul3x1(mgl32.Vec3{-0.5 + edge/2, -0.5 + edge, 0.5 - edge/2})
		pEdge1 := model.Mul3x1(mgl32.Vec3{-0.5 + edge/2, 0.5 - edge, 0.5 - edge/2})
		pTop0 := model.Mul3x1(mgl32.Vec3{-0.5 + edge, -0.5 + edge, 0.5})
		pTop1 := model.Mul3x1(mgl32.Vec3{-0.5 + edge, 0.5 - edge, 0.5})
		pCent := model.Mul3x1(mgl32.Vec3{0.0, 0.0, 0.5})
		pVert0 := model.Mul3x1(mgl32.Vec3{-0.5 + vert, -0.5 + vert, 0.5 - vert})
		pVert1 := model.Mul3x1(mgl32.Vec3{-0.5 + vert, 0.5 - vert, 0.5 - vert})

		uvEdge0 := texture.Mul3x1(mgl32.Vec3{0, 1 - edge, 1})
		uvEdge1 := texture.Mul3x1(mgl32.Vec3{0, edge, 1})
		uvTop0 := texture.Mul3x1(mgl32.Vec3{0 + edge, 1 - edge, 1})
		uvTop1 := texture.Mul3x1(mgl32.Vec3{0 + edge, edge, 1})
		uvCent := texture.Mul3x1(mgl32.Vec3{0.5, 0.5, 1})
		uvVert0 := texture.Mul3x1(mgl32.Vec3{0.0, 1.0, 1})
		uvVert1 := texture.Mul3x1(mgl32.Vec3{0.0, 0.0, 1})

		vEdge0 := gen(pEdge0, nSlope, uvEdge0)
		vEdge1 := gen(pEdge1, nSlope, uvEdge1)
		vTop0 := gen(pTop0, nSlope, uvTop0)
		vTop1 := gen(pTop1, nSlope, uvTop1)

		*v = append(*v, vEdge0, vEdge1, vTop1)
		*v = append(*v, vTop1, vTop0, vEdge0)

		vTop0.setN(nTop)
		vTop1.setN(nTop)
		vCent := gen(pCent, nTop, uvCent)

		*v = append(*v, vTop0, vTop1, vCent)

		vVert0 := gen(pVert0, nVert0, uvVert0)
		vVert1 := gen(pVert1, nVert1, uvVert1)

		vEdge0.setN(nVert0)
		vTop0.setN(nVert0)
		vEdge1.setN(nVert1)
		vTop1.setN(nVert1)

		*v = append(*v, vVert0, vEdge0, vTop0)
		*v = append(*v, vVert1, vTop1, vEdge1)
	})
}

func CubeSideDie(side mgl32.Mat3, v *[]blockVertex) {
	const halfSinPi4 = math.Sqrt2 / 4

	const x2 = halfSinPi4
	const y2 = halfSinPi4

	var x1 = float32(0.5 * math.Cos(math.Pi/8))
	var y1 = float32(0.5 * math.Sin(math.Pi/8))

	var a = float32(1 / math.Sqrt(6))

	quarterSide(side, func(model, texture mgl32.Mat3) {
		nTop := model.Mul3x1(mgl32.Vec3{0, 0, 1})

		p2L := model.Mul3x1(mgl32.Vec3{-x2, y2, 0.5})
		p1L := model.Mul3x1(mgl32.Vec3{-x1, y1, 0.5})
		p00 := model.Mul3x1(mgl32.Vec3{-0.5, 0, 0.5})
		p1R := model.Mul3x1(mgl32.Vec3{-x1, -y1, 0.5})
		p2R := model.Mul3x1(mgl32.Vec3{-x2, -y2, 0.5})
		p0c := model.Mul3x1(mgl32.Vec3{0, 0, 0.5})

		uv2L := texture.Mul3x1(mgl32.Vec3{0.5 - x2, 0.5 - y2, 1})
		uv1L := texture.Mul3x1(mgl32.Vec3{0.5 - x1, 0.5 - y1, 1})
		uv00 := texture.Mul3x1(mgl32.Vec3{0, 0.5, 1})
		uv1R := texture.Mul3x1(mgl32.Vec3{0.5 - x1, 0.5 + y1, 1})
		uv2R := texture.Mul3x1(mgl32.Vec3{0.5 - x2, 0.5 + y2, 1})
		uv0c := texture.Mul3x1(mgl32.Vec3{0.5, 0.5, 1})

		v2L := gen(p2L, nTop, uv2L) //   / v2L
		v1L := gen(p1L, nTop, uv1L) //  / v1L
		v00 := gen(p00, nTop, uv00) // | v00  v0c
		v1R := gen(p1R, nTop, uv1R) //  \ v1R
		v2R := gen(p2R, nTop, uv2R) //   \ v2R
		v0c := gen(p0c, nTop, uv0c) //

		*v = append(*v, v0c, v2R, v1R)
		*v = append(*v, v0c, v1R, v00)
		*v = append(*v, v0c, v00, v1L)
		*v = append(*v, v0c, v1L, v2L)

		_ = v2L
		_ = v1L
		_ = v2R
		_ = v1R
		_ = v00
		_ = v0c

		v2R.setN(p2R.Normalize())
		v1R.setN(p1R.Normalize())
		v00.setN(p00.Normalize())
		v1L.setN(p1L.Normalize())
		v2L.setN(p2L.Normalize())

		pVL := model.Mul3x1(mgl32.Vec3{-a, a, a})
		pVR := model.Mul3x1(mgl32.Vec3{-a, -a, a})

		pHL := pVL.Add(p00).Normalize().Mul(0.5 * math.Sqrt2)
		pHR := pVR.Add(p00).Normalize().Mul(0.5 * math.Sqrt2)

		uvVL := texture.Mul3x1(mgl32.Vec3{0, 0, 1})
		uvHL := texture.Mul3x1(mgl32.Vec3{0, 0.25, 1})
		uvHR := texture.Mul3x1(mgl32.Vec3{0, 0.75, 1})
		uvVR := texture.Mul3x1(mgl32.Vec3{0, 1, 1})

		vVL := gen(pVL, pVL.Normalize(), uvVL)
		vHL := gen(pHL, pHL.Normalize(), uvHL)
		vHR := gen(pHR, pHR.Normalize(), uvHR)
		vVR := gen(pVR, pVR.Normalize(), uvVR)

		*v = append(*v, vHR, v2R, vVR)
		*v = append(*v, vHR, v1R, v2R)
		*v = append(*v, vHR, v00, v1R)

		*v = append(*v, vHL, vVL, v2L)
		*v = append(*v, vHL, v2L, v1L)
		*v = append(*v, vHL, v1L, v00)

		_ = vVL
		_ = vHL
		_ = vHR
		_ = vVR
	})
}

func CubeSideHexagonalStar(k float32) func(mgl32.Mat3, *[]blockVertex) {
	// k = 0.5 is rhombic dodecahedron
	if k < 0.01 {
		k = 0.01
	} else if k > 1 {
		k = 1 // cube
	}
	return func(side mgl32.Mat3, v *[]blockVertex) {
		// should be: height < 0.5
		// should be: base < 0.5 (for star use base = height / 4)

		const height = 0.5
		base := k * height

		n := mgl32.Vec3{-base, 0, height - base}.Normalize()

		p0 := mgl32.Vec3{-base, -base, base}
		p1 := mgl32.Vec3{-base, base, base}
		pc := mgl32.Vec3{0, 0, height}

		uv0 := mgl32.Vec3{0, 1, 1}
		uv1 := mgl32.Vec3{0, 0, 1}
		uvc := mgl32.Vec3{0.5, 0.5, 1}

		quarterSide(side, func(model, texture mgl32.Mat3) {
			norm := model.Mul3x1(n)
			*v = append(*v,
				gen(model.Mul3x1(p0), norm, texture.Mul3x1(uv0)),
				gen(model.Mul3x1(p1), norm, texture.Mul3x1(uv1)),
				gen(model.Mul3x1(pc), norm, texture.Mul3x1(uvc)),
			)
		})
	}
}

func CubeSideOctagonalStar(k float32) func(mat3 mgl32.Mat3, v *[]blockVertex) {
	if k < 0.01 {
		k = 0.01
	} else if k > 1 {
		k = 1 // stellated octahedron
	}
	return func(side mgl32.Mat3, v *[]blockVertex) {
		const a = 0.5 // height = a * math.Sqrt2
		b := k * a    // base = b * math.Sqrt2

		p0 := mgl32.Vec3{-a, -a, a}
		p1 := mgl32.Vec3{-b, 0, 0}
		p2 := mgl32.Vec3{0, 0, b}

		uv0 := mgl32.Vec3{0, 1, 1}
		uv1 := mgl32.Vec3{(1 - k) * 0.5, 0.5, 1}
		uv2 := mgl32.Vec3{0.5, 0.5, 1}

		n := p2.Sub(p0).Cross(p1.Sub(p0)).Normalize()

		quarterSide(side, func(model, texture mgl32.Mat3) {
			norm := model.Mul3x1(n)
			*v = append(*v,
				gen(model.Mul3x1(p0), norm, texture.Mul3x1(uv0)),
				gen(model.Mul3x1(p1), norm, texture.Mul3x1(uv1)),
				gen(model.Mul3x1(p2), norm, texture.Mul3x1(uv2)),
			)
		})
	}
}

func halfSide(model mgl32.Mat3, halfSizePart func(model, texture mgl32.Mat3)) {
	for i := 0; i < 2; i++ {
		angle := float32(i) * math.Pi
		texture := mgl32.Ident3().
			Mul3(mgl32.Translate2D(0.5, 0.5)).
			Mul3(mgl32.HomogRotate2D(angle)).
			Mul3(mgl32.Translate2D(-0.5, -0.5))
		modelSidePart := model.Mul3(mgl32.Rotate3DZ(angle))

		//  |\ 1: p=-0.5,0.5,0.5 uv=0,0
		//  | \  center: p=0,0,0.5 uv=0.5,0.5
		//  |  \   up=0,0,1
		//  |---\ 2: p=0.5,-0.5,0.5 uv=1,1
		//  0: p=-0.5,-0.5,0 uv=0,1

		halfSizePart(modelSidePart, texture)
	}
}

func quarterSide(model mgl32.Mat3, quarterSizePart func(model, texture mgl32.Mat3)) {
	for i := 0; i < 4; i++ {
		angle := float32(i) * math.Pi / 2.0
		texture := mgl32.Ident3().
			Mul3(mgl32.Translate2D(0.5, 0.5)).
			Mul3(mgl32.HomogRotate2D(angle)).
			Mul3(mgl32.Translate2D(-0.5, -0.5))
		modelSidePart := model.Mul3(mgl32.Rotate3DZ(-angle))

		//  |\ 1: p=-0.5,0.5,0.5 uv=0,0
		//  | \  center: p=0,0,0.5 uv=0.5,0.5
		//  | /  up=0,0,1
		//  |/ 0: p=-0.5,-0.5,0.5 uv=0,1

		quarterSizePart(modelSidePart, texture)
	}
}

func makeCubeModel(makeSide func(side mgl32.Mat3, v *[]blockVertex)) []blockVertex {
	v := make([]blockVertex, 0, 64)

	center := mgl32.Ident3()
	for i := 0; i < 6; i++ {
		var side mgl32.Mat3

		switch i {
		case 0:
			// front
			side = center
		case 1:
			// back
			side = center.
				Mul3(mgl32.Rotate3DX(math.Pi))
		case 2:
			// top
			side = center.
				Mul3(mgl32.Rotate3DZ(math.Pi / 2)).
				Mul3(mgl32.Rotate3DY(math.Pi / 2))
		case 3:
			// bottom
			side = center.
				Mul3(mgl32.Rotate3DZ(math.Pi / 2)).
				Mul3(mgl32.Rotate3DY(-math.Pi / 2))
		case 4:
			// left
			side = center.
				Mul3(mgl32.Rotate3DZ(-math.Pi / 2)).
				Mul3(mgl32.Rotate3DX(math.Pi / 2))
		case 5:
			// right
			side = center.
				Mul3(mgl32.Rotate3DZ(-math.Pi / 2)).
				Mul3(mgl32.Rotate3DX(-math.Pi / 2))
		}

		makeSide(side, &v)
	}

	return v
}
