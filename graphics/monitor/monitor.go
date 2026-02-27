// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package monitor

import (
	"fmt"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func GetMonitorByName(name string) *glfw.Monitor {
	if name != "" {
		monitors := glfw.GetMonitors()
		for _, m := range monitors {
			if m == nil {
				continue
			}

			if n := m.GetName(); n != "" && name == n {
				return m
			}
		}
	}

	return nil
}

type Monitor struct {
	Name           string
	Width          int
	Height         int
	ScaleWidth     float32
	ScaleHeight    float32
	PhysicalWidth  int
	PhysicalHeight int
}

func (m *Monitor) String() string {
	if m.Width == 0 || m.Height == 0 {
		return m.Name
	}

	w := int(float32(m.Width) * m.ScaleWidth)
	h := int(float32(m.Height) * m.ScaleHeight)

	return fmt.Sprintf("%s (%dx%d)", m.Name, w, h)
}

func GetMonitors() []Monitor {
	monitors := glfw.GetMonitors()
	list := make([]Monitor, 0, len(monitors))
	for _, monitor := range monitors {
		if monitor == nil {
			continue
		}

		m := getMonitor(monitor)

		if m.Name == "" {
			continue
		}

		list = append(list, m)
	}

	return list
}

func getMonitor(monitor *glfw.Monitor) Monitor {
	m := Monitor{
		Name: monitor.GetName(),
	}

	if vm := monitor.GetVideoMode(); vm != nil {
		m.Width = vm.Width
		m.Height = vm.Height
	}

	m.ScaleWidth, m.ScaleHeight = monitor.GetContentScale()
	m.PhysicalWidth, m.PhysicalHeight = monitor.GetPhysicalSize()

	return m
}
