// Copyright (c) 2024 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package app

import (
	"reflect"
	"testing"
)

func TestRoutes(t *testing.T) {
	tests := []struct {
		name    string
		ids     []route
		exp     []route
		expRmvd []route
	}{
		{
			name:    "push-one",
			ids:     []route{"one"},
			exp:     []route{"one"},
			expRmvd: nil,
		},
		{
			name:    "pop-one",
			ids:     []route{""},
			exp:     []route{""},
			expRmvd: []route{""},
		},
		{
			name:    "push-2-pop-1",
			ids:     []route{"1", "2", ""},
			exp:     []route{"1", "2", "1"},
			expRmvd: []route{"2"},
		},
		{
			name:    "push-4-pop-2",
			ids:     []route{"1", "2", "3", "", "", "4", ""},
			exp:     []route{"1", "2", "3", "2", "1", "4", "1"},
			expRmvd: []route{"3", "2", "4"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var r routes
			var rmvd []route
			for i, id := range test.ids {
				if id == "" {
					rmvd = append(rmvd, r.pop())
				} else {
					r.push(id)
				}

				if want, got := test.exp[i], r.curr(); want != got {
					t.Errorf("idx=%d, want=%s got=%s", i, want, got)
				}
			}

			if want, got := test.expRmvd, rmvd; !reflect.DeepEqual(want, got) {
				t.Errorf("rmvd: want=%s got=%s", want, got)
			}
		})
	}
}
