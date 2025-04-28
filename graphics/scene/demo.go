// Copyright (c) 2024, 2025 by Marko Gaćeša

package scene

const DemoBlocks DemoName = "demo-blocks"
const DemoFields DemoName = "demo-fields"

type DemoName string

type DemoScreenConfig struct {
	Name DemoName
}

func Demo(name DemoName) DemoScreenConfig {
	return DemoScreenConfig{Name: name}
}
