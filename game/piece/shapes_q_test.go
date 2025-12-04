// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package piece

import (
	"testing"
)

func TestShapesQ(t *testing.T) {
	q := _initShapeRect(5, 6, []bool{
		__, XX, XX, XX, __,
		XX, __, __, __, XX,
		__, __, XX, XX, __,
		__, __, XX, __, __,
		__, __, __, __, __,
		__, __, XX, __, __,
	})
	if want, got := q, shapeQ; want != got {
		t.Errorf("want = %s got %s", want.def(), got.def())
	}
}

func TestShapesO(t *testing.T) {
	o := _initShapeRect(4, 4, []bool{
		__, XX, XX, __,
		XX, __, __, XX,
		XX, __, __, XX,
		__, XX, XX, __,
	})
	if want, got := o, shapeO; want != got {
		t.Errorf("want = %s got %s", want.def(), got.def())
	}
}
