// Copyright (c) 2020 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package anim

import (
	"testing"
	"time"
)

func TestList_Update(t *testing.T) {
	startTime := time.Date(2020, 11, 01, 12, 30, 45, 0, time.UTC)

	type expect struct {
		delay time.Duration
		count int
		dx    float32
		dy    float32
		dz    float32
	}

	tests := []struct {
		name     string
		anims    []Anim
		expected []expect
	}{
		{
			name:     "empty",
			anims:    []Anim{},
			expected: []expect{{}}, // one element: expect no changes
		},
		{
			name: "one anim",
			anims: []Anim{
				NewTransLin(startTime, 60*time.Second, 100, 0, 0),
			},
			expected: []expect{
				{
					delay: 0,
					count: 1,
					dx:    -100,
					dy:    0,
					dz:    0,
				},
				{
					delay: 30 * time.Second,
					count: 1,
					dx:    -50,
					dy:    0,
					dz:    0,
				},
				{
					delay: 60 * time.Second,
					count: 0,
					dx:    0,
					dy:    0,
					dz:    0,
				},
			},
		},
		{
			name: "two anims: the first is short",
			anims: []Anim{
				NewTransLin(startTime, 60*time.Second, 100, 0, 0),
				NewTransLin(startTime, 30*time.Second, 0, 100, 0),
			},
			expected: []expect{
				{
					delay: 0,
					count: 2,
					dx:    -100,
					dy:    -100,
					dz:    0,
				},
				{
					delay: 15 * time.Second,
					count: 2,
					dx:    -75,
					dy:    -50,
					dz:    0,
				},
				{
					delay: 30 * time.Second,
					count: 1,
					dx:    -50,
					dy:    0,
					dz:    0,
				},
				{
					delay: 60 * time.Second,
					count: 0,
					dx:    0,
					dy:    0,
					dz:    0,
				},
			},
		},
		{
			name: "two anims: the second is short",
			anims: []Anim{
				NewTransLin(startTime, 30*time.Second, 100, 0, 0),
				NewTransLin(startTime, 60*time.Second, 0, 100, 0),
			},
			expected: []expect{
				{
					delay: 15 * time.Second,
					count: 2,
					dx:    -50,
					dy:    -75,
					dz:    0,
				},
				{
					delay: 45 * time.Second,
					count: 1,
					dx:    0,
					dy:    -25,
					dz:    0,
				},
				{
					delay: 60 * time.Second,
					count: 0,
					dx:    0,
					dy:    0,
					dz:    0,
				},
			},
		},
		{
			name: "three anims: the middle one is short",
			anims: []Anim{
				NewTransLin(startTime, 120*time.Second, 0, 0, 100),
				NewTransLin(startTime, 30*time.Second, 100, 0, 0),
				NewTransLin(startTime, 60*time.Second, 0, 100, 0),
			},
			expected: []expect{
				{
					delay: 15 * time.Second,
					count: 3,
					dx:    -50,
					dy:    -75,
					dz:    -87.5,
				},
				{
					delay: 45 * time.Second,
					count: 2,
					dx:    0,
					dy:    -25,
					dz:    -62.5,
				},
				{
					delay: 90 * time.Second,
					count: 1,
					dx:    0,
					dy:    0,
					dz:    -25,
				},
				{
					delay: 200 * time.Second,
					count: 0,
					dx:    0,
					dy:    0,
					dz:    0,
				},
			},
		},
	}

	for _, test := range tests {
		list := &List{}

		for _, anim := range test.anims {
			list.Add(anim)
		}

		for _, e := range test.expected {
			now := startTime.Add(e.delay)
			r := list.Process(now)

			count := list.Count()

			if count != e.count {
				t.Errorf("test %q failed; after %q expected count %d but got %d", test.name, e.delay, e.count, count)
				continue
			}

			if r.DX != e.dx {
				t.Errorf("test %q failed; after %q expected dx=%f but got %f", test.name, e.delay, e.dx, r.DX)
				continue
			}
			if r.DY != e.dy {
				t.Errorf("test %q failed; after %q expected dy=%f but got %f", test.name, e.delay, e.dy, r.DY)
				continue
			}
			if r.DZ != e.dz {
				t.Errorf("test %q failed; after %q expected dz=%f but got %f", test.name, e.delay, e.dz, r.DZ)
				continue
			}
		}
	}
}

func TestList_Update_Features(t *testing.T) {
	startTime := time.Date(2020, 11, 01, 12, 30, 45, 0, time.UTC)

	type expect struct {
		delay time.Duration
		count int
		f     Feature
	}

	tests := []struct {
		name     string
		anims    []Anim
		expected []expect
	}{
		{
			name: "two anims: the first one is short",
			anims: []Anim{
				NewXLin(startTime, 30*time.Second, 100),
				NewZRotLin(startTime, 60*time.Second, 100),
			},
			expected: []expect{
				{
					delay: 15 * time.Second,
					count: 2,
					f:     Translate | Rotate,
				},
				{
					delay: 45 * time.Second,
					count: 1,
					f:     Rotate,
				},
				{
					delay: 60 * time.Second,
				},
			},
		},
		{
			name: "different anims: the middle one is short",
			anims: []Anim{
				NewPopIn(startTime, 120*time.Second),
				NewXLin(startTime, 30*time.Second, 100),
				NewZRotLin(startTime, 60*time.Second, 100),
			},
			expected: []expect{
				{
					delay: 15 * time.Second,
					count: 3,
					f:     Translate | Scale | Rotate,
				},
				{
					delay: 45 * time.Second,
					count: 2,
					f:     Scale | Rotate,
				},
				{
					delay: 90 * time.Second,
					count: 1,
					f:     Scale,
				},
				{
					delay: 200 * time.Second,
				},
			},
		},
	}

	list := List{}

	for _, test := range tests {
		//x.list.Clear()

		for _, anim := range test.anims {
			list.Add(anim)
		}

		for _, e := range test.expected {
			now := startTime.Add(e.delay)
			r := list.Process(now)

			count := list.Count()

			if count != e.count {
				t.Errorf("test %q failed; after %q expected count %d but got %d", test.name, e.delay, e.count, count)
				continue
			}

			if r.Feature != e.f {
				t.Errorf("test %q failed; after %q expected feature=%b but got %b", test.name, e.delay, e.f, r.Feature)
				continue
			}
		}
	}
}
