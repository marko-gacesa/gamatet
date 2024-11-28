// Copyright (c) 2024 by Marko Gaćeša

package scene

import "context"

const DemoBlocks = "demo-blocks"
const DemoFields = "demo-fields"

type DemoScreenConfig struct {
	Name string
	Stop context.CancelFunc
}

func Demo(ctx context.Context, name string) (DemoScreenConfig, context.Context) {
	ctx, cancelCtx := context.WithCancel(ctx)
	return DemoScreenConfig{
		Name: name,
		Stop: cancelCtx,
	}, ctx
}
