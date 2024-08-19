// Copyright (c) 2020-2024 by Marko Gaćeša

package material

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type program struct {
	programID          uint32
	meshPrimitiveType  uint32
	meshPrimitiveCount int32
}

const programOutputFragData = "outputColor"

func newProgram(vertexShaderSource, fragmentShaderSource string) (*program, error) {
	programID, err := compile(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		return nil, err
	}

	gl.BindFragDataLocation(programID, 0, gl.Str(programOutputFragData+z)) // fragment program output

	p := program{
		programID:          programID,
		meshPrimitiveType:  0,
		meshPrimitiveCount: 0,
	}

	return &p, nil
}

func (p *program) Use()    { gl.UseProgram(p.programID) }
func (p *program) Reset()  {}
func (p *program) Render() { gl.DrawArrays(p.meshPrimitiveType, 0, p.meshPrimitiveCount) }
func (p *program) Delete() { gl.DeleteProgram(p.programID) }

func (p *program) uniformLocation(name string) int32 {
	return gl.GetUniformLocation(p.programID, gl.Str(name+z))
}
func (p *program) attribLocation(name string) uint32 {
	return uint32(gl.GetAttribLocation(p.programID, gl.Str(name+z)))
}

func uniform1i(uni int32, i int)     { gl.Uniform1i(uni, int32(i)) }
func uniform1f(uni int32, v float32) { gl.Uniform1f(uni, v) }

func uniformVec2(uni int32, v mgl32.Vec2)   { gl.Uniform2fv(uni, 1, &v[0]) }
func uniformVec3(uni int32, v mgl32.Vec3)   { gl.Uniform3fv(uni, 1, &v[0]) }
func uniformVec4(uni int32, v mgl32.Vec4)   { gl.Uniform4fv(uni, 1, &v[0]) }
func uniformMat3(uni int32, v *mgl32.Mat3)  { gl.UniformMatrix3fv(uni, 1, false, &v[0]) }
func uniformMat4(uni int32, v *mgl32.Mat4)  { gl.UniformMatrix4fv(uni, 1, false, &v[0]) }
func uniformMat3T(uni int32, v *mgl32.Mat3) { gl.UniformMatrix3fv(uni, 1, true, &v[0]) }
func uniformMat4T(uni int32, v *mgl32.Mat4) { gl.UniformMatrix4fv(uni, 1, true, &v[0]) }

func uniformTexture(uni int32, texture uint32) {
	gl.Uniform1i(uni, int32(texture)-int32(gl.TEXTURE0))
	gl.ActiveTexture(texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
}

func uniformModel(uniModel, uniNormal int32, model *mgl32.Mat4) {
	uniformMat4(uniModel, model)

	// Transforming normals
	// https://www.scratchapixel.com/lessons/mathematics-physics-for-computer-graphics/geometry/transforming-normals
	normal := model.Mat3().Inv()
	uniformMat3T(uniNormal, &normal)
}
