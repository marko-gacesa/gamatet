// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package setup

import "testing"

func TestPresets(t *testing.T) {
	for i := range SinglePlayerPresetCount {
		setup := SinglePlayerPreset(i)
		if setup.SanitizeSingle() {
			t.Errorf("single player preset %d required sanitation", i)
		}
	}

	for i := range MultiPlayerPresetCount {
		setup := MultiPlayerPreset(i)
		if setup.SanitizeMulti() {
			t.Errorf("multi player preset %d required sanitation", i)
		}
	}
}
