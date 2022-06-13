package v2

import (
	"sync"
	"time"
)

var globalSpanID spanID = 0
var globalSpanMux sync.Mutex

type (
	spanID int64
	depth  int64

	Span struct {
		id      spanID  // unique spanID
		parent  *Span   // ref to parent, nil on root spans
		child   []*Span // refs to all child
		depth   depth   // 0 = root, +1 for child
		logical bool    // span will not store logs, and propagate it next to non-logical parent

		title     string    // span title to display
		container container // logs container, layout depend on terminal spawner
		progress  int       // progress in %, 0 .. 100

		startAt  time.Time
		endAt    time.Time
		finished bool

		sync.RWMutex
	}
)

// todo: notify channel
func newSpan(parent *Span, container container, logical bool, title string) *Span {
	globalSpanMux.Lock()
	defer globalSpanMux.Unlock()

	globalSpanID++

	span := &Span{
		id:      globalSpanID,
		parent:  parent,
		child:   make([]*Span, 0),
		logical: logical,

		title:     title,
		container: container,
		progress:  0,

		startAt:  time.Now(),
		endAt:    time.Time{},
		finished: false,
	}

	if parent != nil {
		span.depth = parent.depth + 1
		parent.child = append(parent.child, span)
	}

	return span
}

// Append log to this span
func (s *Span) Write(src string) {
	if s == nil {
		return
	}

	s.Lock()
	defer s.Unlock()

	if s.logical {
		// propagate next to physical parent
		s.parent.Write(src)
		return
	}

	s.container.write(src)
	// todo: notify here
}

// UpdateProgress get any value between 0 and 1
// where 1 = 100% and output this progress in terminal
func (s *Span) UpdateProgress(progress float64) {
	if s == nil {
		return
	}

	s.Lock()
	defer s.Unlock()

	if s.finished {
		return
	}

	if progress > 1 {
		progress = 1
	}

	if progress < 0 {
		progress = 0
	}

	s.progress = int(progress * 100)
}

// End will close this span
// It will ignore all other method calls to this span
// also time took will be calculated after span ending
// this method will automatically close all child spans (if it`s not closed yet)
func (s *Span) End() {
	if s == nil {
		return
	}

	s.Lock()
	defer s.Unlock()

	if s.finished {
		return
	}

	for _, subSpan := range s.child {
		subSpan.End()
	}

	s.progress = 100
	s.finished = true
	s.endAt = time.Now()
	s.container = nil
}
