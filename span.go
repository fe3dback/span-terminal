package terminal

import (
	"sync"
	"time"
)

var globalId spanID
var globalSpanMux sync.RWMutex

type (
	spanID int64
	spans  map[spanID]*Span

	Span struct {
		term     *Terminal
		parent   *Span
		child    spans
		depth    int
		id       spanID
		title    string
		lastLine string
		percent  int

		startAt time.Time
		endAt   time.Time
		isEnd   bool
	}
)

func newSpan(term *Terminal, prev *Span, title string) *Span {
	globalSpanMux.Lock()
	defer globalSpanMux.Unlock()

	globalId++
	currentTime := time.Now()

	newSpan := &Span{
		term:    term,
		parent:  prev,
		depth:   0,
		child:   make(spans),
		id:      globalId,
		title:   title,
		startAt: currentTime,
		endAt:   currentTime,
		isEnd:   false,
	}

	if prev != nil {
		prev.child[newSpan.id] = newSpan
		newSpan.depth = prev.depth + 1
	}

	return newSpan
}

func (s *Span) WriteMessage(result string) {
	if s == nil {
		return
	}
	if s.isEnd {
		return
	}

	s.updateLastLine(result)
	s.write(newAction(s, actionTypeMessage, result))
}

func (s *Span) UpdateProgress(progress float64) {
	if s == nil {
		return
	}
	if s.isEnd {
		return
	}

	if progress < 0 {
		progress = 0
	}

	if progress > 1 {
		progress = 1
	}

	s.percent = int(progress * 100)
}

func (s *Span) End() {
	if s == nil {
		return
	}
	if s.isEnd {
		return
	}

	s.write(newAction(s, actionTypeSpanEnd, ""))

	s.endAt = time.Now()
	s.percent = 100
	s.isEnd = true
}

func (s *Span) updateLastLine(result string) {
	s.lastLine = result

	if s.parent != nil {
		s.parent.updateLastLine(result)
	}
}

func (s *Span) write(act action) {
	if s.term.finished {
		return
	}

	s.term.buffer <- act
}

func (s *Span) isRoot() bool {
	return s.parent == nil
}

func (s *Span) isFinished() bool {
	return s.isEnd
}

func (s *Span) getParent() *Span {
	return s.parent
}

func (s *Span) getChild() map[spanID]*Span {
	return s.child
}

func (s *Span) hasActiveChild() bool {
	globalSpanMux.RLock()
	defer globalSpanMux.RUnlock()

	for _, span := range s.child {
		if !span.isFinished() {
			return true
		}
	}

	return false
}
