// Copyright (c) 2020 by Marko Gaćeša

package material

import (
	"fmt"
	"gamatet/graphics/camera"
	"gamatet/graphics/geometry"
	"gamatet/graphics/gtypes"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"reflect"
	"time"
	"unsafe"
)

type program struct {
	program  uint32
	uniforms map[string]int32

	uniView  int32
	uniModel int32
	uniNorm  int32
	uniTime  int32
	uniTex   int32
	uniColor int32

	attribVert  uint32
	attribNorm  uint32
	attribTexUV uint32

	meshPrimitiveType  uint32
	meshPrimitiveCount int32
}

const z = "\000"

func newProgram(vertexShaderSource, fragmentShaderSource string) (*program, error) {
	programID, err := compile(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		return nil, err
	}

	uniforms := make(map[string]int32)

	uniView := gl.GetUniformLocation(programID, gl.Str("viewMatrix"+z))
	uniModel := gl.GetUniformLocation(programID, gl.Str("modelMatrix"+z))
	uniNorm := gl.GetUniformLocation(programID, gl.Str("normalMatrix"+z))
	uniTime := gl.GetUniformLocation(programID, gl.Str("time"+z))
	uniTex := gl.GetUniformLocation(programID, gl.Str("textureSampler"+z))
	uniColor := gl.GetUniformLocation(programID, gl.Str("objectColor"+z))

	// fragment program output
	gl.BindFragDataLocation(programID, 0, gl.Str("outputColor"+z))

	attribVert := uint32(gl.GetAttribLocation(programID, gl.Str("geometryPosition"+z)))
	attribNorm := uint32(gl.GetAttribLocation(programID, gl.Str("geometryNormal"+z)))
	attribTexUV := uint32(gl.GetAttribLocation(programID, gl.Str("geometryTexture"+z)))

	p := program{
		program: programID,

		uniforms: uniforms,
		uniView:  uniView,
		uniModel: uniModel,
		uniNorm:  uniNorm,
		uniTex:   uniTex,
		uniTime:  uniTime,
		uniColor: uniColor,

		attribVert:  attribVert,
		attribNorm:  attribNorm,
		attribTexUV: attribTexUV,

		meshPrimitiveType:  0,
		meshPrimitiveCount: 0,
	}

	return &p, nil
}

func (p *program) Use() {
	gl.UseProgram(p.program)
}

func (p *program) registerUniform(name string) {
	uniLoc := gl.GetUniformLocation(p.program, gl.Str(name+z))
	p.uniforms[name] = uniLoc
}

func (p *program) Render() {
	gl.DrawArrays(p.meshPrimitiveType, 0, p.meshPrimitiveCount)
}

func (p *program) Camera(cam *camera.Camera) {
	gl.UniformMatrix4fv(p.uniView, 1, false, &cam.GetView()[0])
}

func (p *program) Model(model *mgl32.Mat4) {
	gl.UniformMatrix4fv(p.uniModel, 1, false, &model[0])

	// Transforming normals
	// https://www.scratchapixel.com/lessons/mathematics-physics-for-computer-graphics/geometry/transforming-normals
	normal := model.Mat3().Inv().Transpose()
	gl.UniformMatrix3fv(p.uniNorm, 1, false, &normal[0])
}

func (p *program) Refresh() {
	gl.Uniform1f(p.uniTime, float32(time.Now().Sub(gtypes.Time).Seconds()))
}

func (p *program) Color(color mgl32.Vec4) {
	gl.Uniform4fv(p.uniColor, 1, &color[0])
}

func (p *program) Texture(texture uint32) {
	gl.Uniform1i(p.uniTex, int32(texture)-int32(gl.TEXTURE0))
	gl.ActiveTexture(texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
}

func (p *program) uniform1i(name string, i int) {
	gl.Uniform1i(p.uniforms[name], int32(i))
}

func (p *program) uniform1f(name string, v float32) {
	gl.Uniform1f(p.uniforms[name], v)
}

func (p *program) uniformVec3(name string, v mgl32.Vec3) {
	gl.Uniform3fv(p.uniforms[name], 1, &v[0])
}

func (p *program) uniformVec4(name string, v mgl32.Vec4) {
	gl.Uniform4fv(p.uniforms[name], 1, &v[0])
}

func (p *program) uniformMat3(name string, v *mgl32.Mat3) {
	gl.UniformMatrix3fv(p.uniforms[name], 1, false, &v[0])
}

func (p *program) uniformMat4(name string, v *mgl32.Mat4) {
	gl.UniformMatrix4fv(p.uniforms[name], 1, false, &v[0])
}

func (p *program) uniformTexture(name string, texture uint32) {
	gl.Uniform1i(p.uniforms[name], int32(texture)-int32(gl.TEXTURE0))
	gl.ActiveTexture(texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
}

func (p *program) Geometry(geometry geometry.Geometry) {
	geometry.GLBind()

	vertexSize := int32(geometry.VertexSize())

	gl.EnableVertexAttribArray(p.attribVert)
	gl.VertexAttribPointer(p.attribVert, 3, gl.FLOAT, false, vertexSize,
		gl.PtrOffset(geometry.DataOffsetVertex()))

	gl.EnableVertexAttribArray(p.attribNorm)
	gl.VertexAttribPointer(p.attribNorm, 3, gl.FLOAT, false, vertexSize,
		gl.PtrOffset(geometry.DataOffsetNormal()))

	gl.EnableVertexAttribArray(p.attribTexUV)
	gl.VertexAttribPointer(p.attribTexUV, 2, gl.FLOAT, false, vertexSize,
		gl.PtrOffset(geometry.DataOffsetTextureUV()))

	p.meshPrimitiveCount = int32(geometry.VertexCount())
	p.meshPrimitiveType = geometry.PrimitiveType()
}

func compile(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	prog := gl.CreateProgram()

	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)

	var status int32
	gl.GetProgramiv(prog, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(prog, gl.INFO_LOG_LENGTH, &logLength)

		log, pLog := getLogBuffer(logLength)
		gl.GetProgramInfoLog(prog, logLength, nil, pLog)

		return 0, fmt.Errorf("failed to link shader: %s", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return prog, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()

	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log, pLog := getLogBuffer(logLength)
		gl.GetShaderInfoLog(shader, logLength, nil, pLog)

		return 0, fmt.Errorf("failed to compile:\n%s", log)
	}

	return shader, nil
}

func getLogBuffer(length int32) (b []byte, ptr *uint8) {
	b = make([]byte, int(length))
	header := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	ptr = (*uint8)(unsafe.Pointer(header.Data))
	return
}
