package terminal

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	tsize "github.com/kopoli/go-terminal-size"
)

// clear screen
const osClear = "\033[2J"

type termOS struct {
	terminal *os.File
	writer   *bufio.Writer
	screen   *bytes.Buffer
}

func newTermOs(terminal *os.File) *termOS {
	return &termOS{
		terminal: terminal,
		writer:   bufio.NewWriter(terminal),
		screen:   new(bytes.Buffer),
	}
}

func (t *termOS) clear() {
	_, _ = t.writer.WriteString(osClear)
}

func (t *termOS) moveCursor(x, y int) {
	_, _ = fmt.Fprintf(t.screen, "\033[%d;%dH", y, x)
}

func (t *termOS) print(src string) {
	_, _ = fmt.Fprint(t.screen, src)
}

func (t *termOS) flush() {
	size, err := tsize.FgetSize(t.terminal)
	if err != nil {
		return
	}

	for idx, str := range strings.SplitAfter(t.screen.String(), "\n") {
		if idx > size.Height {
			return
		}

		_, _ = t.writer.WriteString(str)
	}

	_ = t.writer.Flush()
	t.screen.Reset()
}
