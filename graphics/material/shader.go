// Copyright (c) 2020-2024 by Marko Gaćeša

package material

import (
	"fmt"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
)

const z = "\000"

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
	ptr = unsafe.SliceData(b)
	return
}
