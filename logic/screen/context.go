// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package screen

import "context"

type Context interface {
	context.Context
	Stop()
}

type ctxWrapper struct {
	context.Context
	stop context.CancelFunc
}

func (ctx ctxWrapper) Stop() {
	ctx.stop()
}

func NewContext(ctx context.Context) Context {
	newCtx, cancelNewCtx := context.WithCancel(ctx)
	return &ctxWrapper{Context: newCtx, stop: cancelNewCtx}
}
