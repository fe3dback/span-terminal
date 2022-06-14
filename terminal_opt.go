package terminal

import "os"

type (
	TerminalInitializer = func(*Terminal)
)

func WithWriter(writer *os.File) TerminalInitializer {
	return func(terminal *Terminal) {
		terminal.writer = writer
	}
}

func WithContainerMaxLines(maxLines int) TerminalInitializer {
	return func(terminal *Terminal) {
		terminal.optMaxLines = maxLines
	}
}
