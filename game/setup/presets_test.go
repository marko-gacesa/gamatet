// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package setup

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPresets(t *testing.T) {
	for i := range SinglePlayerPresetCount {
		setup := SinglePlayerPreset(i)
		setupClone := setup
		if setupClone.SanitizeSingle() {
			t.Errorf("single player preset %d required sanitation\n%s\n",
				i, cmp.Diff(setupClone, setup))
		}
	}

	for i := range MultiPlayerPresetCount {
		setup := MultiPlayerPreset(i)
		setupClone := setup
		if setupClone.SanitizeMulti() {
			t.Errorf("multi player preset %d required sanitation\n%s\n",
				i, cmp.Diff(setupClone, setup))
		}
	}
}
