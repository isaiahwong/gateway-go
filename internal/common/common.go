package common

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// SignalContext which listens for a signal term which then
// cancels any operations with current context
func SignalContext(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		signal.Stop(sigs)
		close(sigs)

		// breaks line for trailing ^C
		fmt.Println()
		cancel()
	}()

	return ctx
}
