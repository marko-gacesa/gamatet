// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package setup

import "github.com/marko-gacesa/bitdata"

type writer interface {
	Write(w *bitdata.Writer)
}

type reader interface {
	Read(r *bitdata.ReaderError)
}

func Pack[T writer](t T) []byte {
	w := bitdata.NewWriter()
	t.Write(w)
	return w.BitData()
}

func Unpack[T reader](t T, data []byte) error {
	r := bitdata.NewReaderError(data)
	t.Read(r)
	return r.Error()
}
