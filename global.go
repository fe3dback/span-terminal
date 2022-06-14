package terminal

import "context"

var globalTerminal *Terminal
var globalTerminalInitialized = false

func init() {
	globalTerminal = NewTerminal()
}

// StartSpan will extract parent Span from context (if exist) and create new Span from it
// new Span will be written to new context
// In case when terminal is not active, Span can be nil, but it is totally safe to use
func StartSpan(ctx context.Context, title string, opts ...StartOpt) (context.Context, *Span) {
	if globalTerminal == nil {
		return ctx, nil
	}

	return globalTerminal.span(ctx, append([]StartOpt{WithTitle(title)}, opts...)...)
}

// SetGlobalTerminal allow to customize terminal
// and set it as default for span creation and output
// method can be called only once, all other calls
// will be ignored
func SetGlobalTerminal(t *Terminal) {
	if globalTerminalInitialized {
		return
	}

	globalTerminal = t
	globalTerminalInitialized = true
}

// CaptureOutput will capture control on output to stdout/stderr
// and display custom logs from spans
// all other print/logs will be redirected and printed in special
// region alongside span logs
func CaptureOutput() {
	if globalTerminal == nil {
		return
	}

	globalTerminal.capture()
}

// ReleaseOutput will stop terminal for outputting span logs
// release stdout to default control
// all span updates in released mode, may be ignored
func ReleaseOutput() {
	if globalTerminal == nil {
		return
	}

	globalTerminal.release()
}
