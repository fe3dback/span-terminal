package terminal

type (
	depth int64
)

func (d depth) isRoot() bool {
	return d == 0
}

func (d depth) isOperation2nd() bool {
	return d == 1
}

func (d depth) isDetails3nd() bool {
	return d == 2
}

func (d depth) isDetailsOrDeeper3nd() bool {
	return d >= 2
}
