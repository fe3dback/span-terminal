package terminal

type (
	container interface {
		write(string)
		content() []string
	}

	emptyContainer struct{}

	multiLineContainer struct {
		maxLines int
		lines    []string
	}
)

func newContainer(depth depth, maxLines int) container {
	if depth.isRoot() {
		return newMultiLineContainer(maxLines)
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

func (e *emptyContainer) content() []string {
	return nil
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

func (e *multiLineContainer) content() []string {
	return e.lines
}
