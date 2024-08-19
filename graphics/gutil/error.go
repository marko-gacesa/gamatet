// Copyright (c) 2024 by Marko Gaćeša

package gutil

import (
	"errors"
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
)

func GetError() error {
	errCode := gl.GetError()
	switch errCode {
	case gl.NO_ERROR:
		return nil
	case gl.INVALID_ENUM:
		return errors.New("OpenGL Error: INVALID_ENUM")
	case gl.INVALID_VALUE:
		return errors.New("OpenGL Error: INVALID_VALUE")
	case gl.INVALID_OPERATION:
		return errors.New("OpenGL Error: INVALID_OPERATION")
	case gl.OUT_OF_MEMORY:
		return errors.New("OpenGL Error: OUT_OF_MEMORY")
	case gl.INVALID_FRAMEBUFFER_OPERATION:
		return errors.New("OpenGL Error: INVALID_FRAMEBUFFER_OPERATION")
	default:
		return fmt.Errorf("OpenGL Error: %x", errCode)
	}
}
