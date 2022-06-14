package terminal

import (
	"context"
	"os"
	"sync"
	"time"

	tm "github.com/buger/goterm"
)

// logs will be updated at least once per interval
const forceUpdateInterval = time.Millisecond * 50

// pause between normal updates
const pauseBetweenUpdates = time.Millisecond * 5

// todo: capture stdout to container
// todo: close todos

type Terminal struct {
	writer      *os.File
	optMaxLines int

	rootSpans   []*Span
	active      bool
	watchCtx    context.Context
	watchCancel func()

	mux sync.RWMutex
}

func NewTerminal(initializers ...TerminalInitializer) *Terminal {
	t := &Terminal{
		writer:      os.Stdout,
		optMaxLines: 4,

		rootSpans: make([]*Span, 0),
		active:    false,
	}

	for _, initializer := range initializers {
		initializer(t)
	}

	return t
}

func (t *Terminal) capture() {
	t.mux.Lock()
	defer t.mux.Unlock()

	if t.active {
		return
	}

	// create watch context
	ctx, cancel := context.WithCancel(context.Background())

	t.watchCtx = ctx
	t.watchCancel = cancel

	t.active = true
	go t.watch()
}

func (t *Terminal) release() {
	t.mux.Lock()
	defer t.mux.Unlock()

	if !t.active {
		return
	}

	t.active = false
	t.watchCancel()

	// wait for closing watch and current loop processing
	time.Sleep(time.Millisecond * 500)
}

func (t *Terminal) span(ctx context.Context, opts ...StartOpt) (context.Context, *Span) {
	if !t.active {
		return ctx, nil
	}

	parent := spanFromContext(ctx)

	currentDepth := depth(0)
	if parent != nil {
		currentDepth = parent.depth + 1
	}

	newSpan := newSpan(
		parent,
		newContainer(currentDepth, t.optMaxLines),
		!currentDepth.isRoot(),
	)

	for _, enrich := range opts {
		enrich(newSpan)
	}

	if currentDepth == depth(0) {
		t.rootSpans = append(t.rootSpans, newSpan)
	}

	return contextWithSpan(ctx, newSpan), newSpan
}

func (t *Terminal) watch() {
	watching := true

	ticker := time.NewTicker(forceUpdateInterval)
	spanUpdated := make(chan struct{}, 1)
	nextUpdateAt := time.Now()

	go func() {
		lastUpdatedAt := time.Now()

		for watching {
			latestSpanChangeAt := t.latestSpanChangeAt()
			if latestSpanChangeAt.After(lastUpdatedAt) {
				lastUpdatedAt = latestSpanChangeAt
				spanUpdated <- struct{}{}
			}

			time.Sleep(time.Millisecond)
		}
	}()

	for watching {
		select {
		case <-t.watchCtx.Done():
			watching = false
			break
		case <-ticker.C:
			t.update()
			break
		case <-spanUpdated:
			if time.Now().Before(nextUpdateAt) {
				break
			}

			nextUpdateAt = time.Now().Add(pauseBetweenUpdates)
			t.update()
			break
		}
	}
}

func (t *Terminal) latestSpanChangeAt() time.Time {
	latest := time.Time{}

	for _, span := range t.rootSpans {
		if span.changedAt.After(latest) {
			latest = span.changedAt
		}
	}

	return latest
}

func (t *Terminal) update() {
	t.mux.Lock()
	defer t.mux.Unlock()

	tm.Clear()
	tm.MoveCursor(1, 1)

	for _, rootSpan := range t.rootSpans {
		spanContent := renderSpanWithOptions(rootSpan) // todo: opts
		_, _ = tm.Print(spanContent)
	}

	tm.Flush()
}
