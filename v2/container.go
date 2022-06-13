package v2

import "strings"

type (
	container interface {
		write(string)
		content() []byte
	}

	emptyContainer struct{}

	singleLineContainer struct {
		line string
	}

	multiLineContainer struct {
		maxLines int
		lines    []string
	}
)

func newContainer(depth depth, maxLines int) container {
	if depth == 0 {
		return newMultiLineContainer(maxLines)
	}

	if depth == 1 {
		return newMultiLineContainer(maxLines)
	}

	if depth == 3 {
		return newSingleLineContainer()
	}

	return newEmptyContainer()
}

// ------------------

func newEmptyContainer() *emptyContainer {
	return &emptyContainer{}
}

func (e *emptyContainer) write(_ string) {
	return
}

func (e *emptyContainer) content() []byte {
	return nil
}

// ------------------

func newSingleLineContainer() *singleLineContainer {
	return &singleLineContainer{}
}

func (e *singleLineContainer) write(s string) {
	e.line = s
}

func (e *singleLineContainer) content() []byte {
	return []byte(e.line)
}

// ------------------

func newMultiLineContainer(maxLines int) *multiLineContainer {
	return &multiLineContainer{
		maxLines: maxLines,
		lines:    make([]string, 0, maxLines),
	}
}

func (e *multiLineContainer) write(s string) {
	if len(e.lines) < e.maxLines {
		e.lines = append(e.lines, s)
		return
	}

	e.lines = e.lines[1:]
	e.lines = append(e.lines, s)
	return
}

func (e *multiLineContainer) content() []byte {
	return []byte(strings.Join(e.lines, "\n"))
}
