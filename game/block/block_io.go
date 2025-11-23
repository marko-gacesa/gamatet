// Copyright (c) 2020 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package block

import (
	"encoding/binary"
	"io"
)

func (b *Block) Write(w io.Writer) error {
	var buffer [6]byte

	buffer[0] = byte(b.Type)
	buffer[1] = b.Hardness
	binary.LittleEndian.PutUint32(buffer[2:6], b.Color)

	if _, err := w.Write(buffer[:]); err != nil {
		return err
	}

	return nil
}

func (b *Block) Read(r io.Reader) error {
	var buffer [6]byte

	if _, err := io.ReadFull(r, buffer[:]); err != nil {
		return err
	}

	b.Type = Type(buffer[0])
	b.Hardness = buffer[1]
	b.Color = binary.LittleEndian.Uint32(buffer[2:6])

	return nil
}
