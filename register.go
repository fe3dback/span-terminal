package terminal

import (
	"context"
)

var globalTerminal *Terminal

func RegisterTerminal(t *Terminal) {
	globalTerminal = t
}

func StartSpan(ctx context.Context, title string) (context.Context, *Span) {
	if globalTerminal == nil {
		return ctx, nil
	}

	return globalTerminal.span(ctx, title)
}

func Shutdown() {
	if globalTerminal == nil {
		return
	}

	globalTerminal.shutdown()
}
