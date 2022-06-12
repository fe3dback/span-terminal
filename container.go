package terminal

type container struct {
	max   int
	lines []string
}

func newContainer(maxLines int) *container {
	return &container{
		max:   maxLines,
		lines: []string{},
	}
}

func (c *container) append(line string) {
	if len(c.lines) < c.max {
		c.lines = append(c.lines, line)
		return
	}

	c.lines = c.lines[1:]
	c.lines = append(c.lines, line)
}

func (c *container) content() []string {
	return c.lines
}
