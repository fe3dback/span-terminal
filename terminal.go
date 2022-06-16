package terminal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// logs will be updated at least once per interval
const forceUpdateInterval = time.Millisecond * 50

// pause between normal updates
const pauseBetweenUpdates = time.Millisecond * 5

type Terminal struct {
	opts terminalOpts

	isANSITerminal bool
	rootSpans      []*Span
	active         bool
	watchCtx       context.Context
	watchCancel    func()

	stdoutBuffer  *bytes.Buffer
	realStdout    *os.File
	termOs        *termOS
	logsContainer container
	watchFinished chan struct{}

	mux sync.RWMutex
}

func NewTerminal(initializers ...OptsInitializer) *Terminal {
	opt := &terminalOpts{
		containerMaxLines: OptDefaultContainerMaxLines,
		stdoutMaxLines:    OptDefaultStdoutMaxLines,
		renderOpts:        defaultRenderOpts,
	}
	for _, initializer := range initializers {
		initializer(opt)
	}

	return &Terminal{
		opts: *opt,

		isANSITerminal: termenv.ColorProfile() != termenv.Ascii,
		rootSpans:      make([]*Span, 0),
		active:         false,

		stdoutBuffer:  bytes.NewBuffer(nil),
		realStdout:    os.Stdout,
		termOs:        newTermOs(os.Stdout),
		logsContainer: newMultiLineContainer(opt.stdoutMaxLines),
	}
}

func (t *Terminal) capture() {
	t.mux.Lock()
	defer t.mux.Unlock()

	if t.active {
		return
	}

	if !t.isANSITerminal {
		return
	}

	// create watch context
	ctx, cancel := context.WithCancel(context.Background())

	t.watchCtx = ctx
	t.watchCancel = cancel

	lipgloss.SetColorProfile(termenv.ANSI) // force set to simple ANSI
	t.stdoutBuffer.Reset()
	t.watchFinished = make(chan struct{}) // this channel will be closed, after watch is completed

	t.active = true
	go t.redirectAllStdoutToContainer()
	go t.watch()
}

func (t *Terminal) redirectAllStdoutToContainer() {
	bufioStdout(t.watchCtx, func(message bufioMessage) {
		if message.err != nil {
			if !errors.Is(message.err, io.EOF) {
				_, _ = fmt.Fprint(os.Stderr, fmt.Sprintf("failed buffer stdout: %v", message.err))
			}

			return
		}

		t.stdoutBuffer.Write(message.data)
		t.stdoutBuffer.WriteString("\n")
		t.logsContainer.write(string(message.data))
	})
}

func (t *Terminal) release() {
	t.mux.Lock()
	defer t.mux.Unlock()

	if !t.active {
		return
	}

	t.active = false

	time.Sleep(time.Millisecond * 500) // wait for all io term events done
	t.watchCancel()

	// wait for watch is finished gracefully
	<-t.watchFinished
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
		newContainer(currentDepth, t.opts.containerMaxLines),
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
			t.update()                         // last update
			time.Sleep(time.Millisecond * 500) // wait for last term update
			t.dumpBufferedStdout()             // restore buffered logs to stdout
			time.Sleep(time.Millisecond * 500) // wait for buffer writing
			close(t.watchFinished)             // signal that we can finish restoring terminal
			break
		case <-ticker.C:
			t.update() // force update
			break
		case <-spanUpdated:
			if time.Now().Before(nextUpdateAt) {
				break
			}

			nextUpdateAt = time.Now().Add(pauseBetweenUpdates)
			t.update() // something changed
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
	// clear
	t.termOs.clear()
	t.termOs.moveCursor(1, 1)

	// render main logs
	if t.active {
		// don`t show normal stdout, because we
		// dump in normal mode right after release
		t.termOs.print(renderMainContainer(t.logsContainer) + "\n")
	}

	// render top spans
	for _, rootSpan := range mostRelevantSpans(t.rootSpans, t.opts.renderOpts.spansMaxRoots) {
		t.termOs.print(renderSpanWithOptions(rootSpan, t.opts.renderOpts))
	}

	// output to term
	t.termOs.flush()
}

func (t *Terminal) dumpBufferedStdout() {
	// print all captured and hidden messages and logs
	// back to stdout
	_, _ = t.realStdout.WriteString("\n")
	_, _ = t.realStdout.Write(t.stdoutBuffer.Bytes())
	_, _ = t.realStdout.WriteString("\n")
	t.stdoutBuffer.Reset()
}
