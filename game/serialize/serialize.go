// Copyright (c) 2020 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package serialize

import (
	"encoding/binary"
	"fmt"
	"io"
)

func Write8(w io.Writer, i byte) error {
	_, err := w.Write([]byte{i})
	return err
}

func Read8(r io.Reader) (byte, error) {
	var buffer [1]byte

	if _, err := io.ReadFull(r, buffer[:]); err != nil {
		return 0, fmt.Errorf("failed to read int8: %w", err)
	}

	return buffer[0], nil
}

func ReadInt8(r io.Reader) (int8, error) {
	n, err := Read8(r)
	return int8(n), err
}

func Write16(w io.Writer, i uint16) error {
	var buffer [2]byte
	binary.LittleEndian.PutUint16(buffer[:], i)
	_, err := w.Write(buffer[:])
	return err
}

func Read16(r io.Reader) (uint16, error) {
	var buffer [2]byte

	_, err := io.ReadFull(r, buffer[:])
	if err != nil {
		return 0, fmt.Errorf("failed to read int16: %w", err)
	}

	return binary.LittleEndian.Uint16(buffer[:]), nil
}

func ReadInt16(r io.Reader) (int16, error) {
	i, err := Read16(r)
	return int16(i), err
}

func Write32(w io.Writer, i uint32) error {
	var buffer [4]byte
	binary.LittleEndian.PutUint32(buffer[:], i)
	_, err := w.Write(buffer[:])
	return err
}

func Read32(r io.Reader) (uint32, error) {
	var buffer [4]byte

	_, err := io.ReadFull(r, buffer[:])
	if err != nil {
		return 0, fmt.Errorf("failed to read int32: %w", err)
	}

	return binary.LittleEndian.Uint32(buffer[:]), nil
}

func ReadInt(r io.Reader) (int, error) {
	i, err := Read32(r)
	return int(i), err
}

func Write64(w io.Writer, i uint64) error {
	var buffer [8]byte
	binary.LittleEndian.PutUint64(buffer[:], i)
	_, err := w.Write(buffer[:])
	return err
}

func Read64(r io.Reader) (uint64, error) {
	var buffer [8]byte

	_, err := io.ReadFull(r, buffer[:])
	if err != nil {
		return 0, fmt.Errorf("failed to read int64")
	}

	return binary.LittleEndian.Uint64(buffer[:]), nil
}

func ReadInt64(r io.Reader) (int64, error) {
	i, err := Read64(r)
	return int64(i), err
}
