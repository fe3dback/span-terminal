package terminal

import (
	"sync"
	"time"
)

var globalId spanID
var createMux sync.Mutex

type (
	spanID int64
	spans  map[spanID]*Span

	Span struct {
		term   *Terminal
		parent *Span
		child  spans
		id     spanID
		title  string

		startAt time.Time
		endAt   time.Time
		isEnd   bool
	}
)

func newSpan(term *Terminal, prev *Span, title string) *Span {
	createMux.Lock()
	defer createMux.Unlock()

	globalId++
	currentTime := time.Now()

	newSpan := &Span{
		term:    term,
		parent:  prev,
		child:   make(spans),
		id:      globalId,
		title:   title,
		startAt: currentTime,
		endAt:   currentTime,
		isEnd:   false,
	}

	if prev != nil {
		prev.child[newSpan.id] = newSpan
	}

	return newSpan
}

func (s *Span) WriteMessage(result string) {
	if s == nil {
		return
	}

	s.write(newAction(s, actionTypeMessage, result))
}

func (s *Span) End(result string) {
	if s == nil {
		return
	}

	s.WriteMessage(result)
	s.write(newAction(s, actionTypeSpanEnd, ""))

	s.endAt = time.Now()
	s.isEnd = true
}

func (s *Span) write(act action) {
	if s == nil {
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
