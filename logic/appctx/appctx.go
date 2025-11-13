// Copyright (c) 2022 by Marko Gaćeša

package appctx

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Context is application level context. It is canceled when the application receives a signal from the host operating system to terminate.
var Context = func() context.Context {
	ctx, cancelFn := context.WithCancel(context.Background())
	go func() {
		signalStop := make(chan os.Signal, 1)
		signal.Notify(signalStop, syscall.SIGINT, syscall.SIGTERM)

		defer func() {
			signal.Stop(signalStop)
			cancelFn()
		}()

		<-signalStop
	}()

	return ctx
}()
