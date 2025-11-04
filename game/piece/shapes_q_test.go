// Copyright (c) 2025 by Marko Gaćeša

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
		t.Errorf("want = %s, got %s", want.def(), got.def())
	}
}
