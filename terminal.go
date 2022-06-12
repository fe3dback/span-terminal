package terminal

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// how many logs can be buffered
const bufferCapacity = 8

// how many lines for each span
const maxLines = 4

// long lines will be trimmed
const maxLineLength = 80

// max updates rate
const updateInterval = time.Millisecond * 5

type (
	Terminal struct {
		controlCtx context.Context

		topSpans    spans
		buffer      actionChan
		sleepTo     time.Time
		tea         *tea.Program
		lastMessage []byte
		finished    bool

		bufferedText map[spanID]*container

		sync.RWMutex
	}

	actionChan = chan action
)

func NewTerminal(controlCtx context.Context, target *os.File) *Terminal {
	t := &Terminal{
		controlCtx: controlCtx,

		topSpans:     make(spans),
		buffer:       make(actionChan, bufferCapacity),
		bufferedText: make(map[spanID]*container),
		sleepTo:      time.Now(),
		finished:     false,
	}

	lipgloss.SetColorProfile(termenv.ANSI)
	renderer := tea.NewProgram(t, tea.WithOutput(target))
	t.tea = renderer

	go t.watch()
	return t
}

func (t *Terminal) shutdown() {
	if t.finished {
		return
	}

	t.finished = true

	t.rebuildView()
	time.Sleep(time.Millisecond * 100)
	t.tea.Kill()
	time.Sleep(time.Millisecond * 100)
}

func (t *Terminal) span(ctx context.Context, title string) (context.Context, *Span) {
	if t.finished {
		return ctx, nil
	}

	parent := spanFromContext(ctx)
	newSpan := newSpan(t, parent, title)

	t.RWMutex.Lock()
	if newSpan.isRoot() {
		t.topSpans[newSpan.id] = newSpan
	}
	t.bufferedText[newSpan.id] = newContainer(maxLines)
	t.RWMutex.Unlock()

	return contextWithSpan(ctx, newSpan), newSpan
}

func (t *Terminal) watch() {
	teaQuit := make(chan struct{})
	finished := false

	go func() {
		err := t.tea.Start()
		if err != nil {
			finished = true
			fmt.Printf("failed init terminal: %v\n", err)
			return
		}

		<-teaQuit
		t.rebuildView()
		t.tea.Kill()
	}()

	for {
		if t.finished || finished {
			break
		}

		select {
		case <-t.controlCtx.Done():
			finished = true
			break
		case act := <-t.buffer:
			if t.updateState(act) {
				t.rebuildView()
			}
		}
	}

	t.finished = true
	close(teaQuit)
}

func (t *Terminal) updateState(act action) bool {
	if act.Type() != actionTypeSpanEnd && time.Now().Before(t.sleepTo) {
		// do not update so often
		return false
	}

	if act.Type() == actionTypeSpanEnd {
		// update status
		return true
	}

	t.sleepTo = time.Now().Add(updateInterval)

	// -- update span containers

	t.RWMutex.RLock()
	cont, exist := t.bufferedText[act.SpanID()]
	t.RWMutex.RUnlock()

	if exist {
		cont.append(act.String())
	}

	return true
}

func (t *Terminal) rebuildView() {
	var buf bytes.Buffer
	t.renderSpans(&buf, t.topSpans)

	t.lastMessage = buf.Bytes()
}

func (t *Terminal) renderSpans(buf *bytes.Buffer, list spans) {
	sortedList := make([]*Span, 0, len(list))
	for _, span := range list {
		sortedList = append(sortedList, span)
	}

	sort.Slice(sortedList, func(i, j int) bool {
		return sortedList[i].id < sortedList[j].id
	})

	for _, span := range sortedList {
		if span == nil {
			continue
		}

		if span.isFinished() {
			buf.Write(t.outputDone(span))
			continue
		}

		buf.Write(t.outputProcess(span, !span.hasActiveChild()))

		if span.hasActiveChild() {
			t.renderSpans(buf, span.child)
			continue
		}

		t.RWMutex.RLock()
		cont, exist := t.bufferedText[span.id]
		t.RWMutex.RUnlock()

		if !exist {
			continue
		}

		for _, s := range cont.content() {
			buf.Write(t.outputLine(s))
		}
	}
}

func (t *Terminal) outputDone(s *Span) []byte {
	return []byte(
		styleStatusDone.Render(fmt.Sprintf("[ %5s ] %s",
			t.formatDuration(s),
			s.title,
		)) + "\n",
	)
}

func (t *Terminal) formatDuration(s *Span) string {
	took := s.endAt.Sub(s.startAt)

	if took.Seconds() > 1 {
		return fmt.Sprintf("%.0fs", took.Seconds())
	}

	return fmt.Sprintf("%dms", took.Milliseconds())
}

func (t *Terminal) outputProcess(s *Span, active bool) []byte {
	text := s.title
	pb := fmt.Sprintf("[ %4d%% ] ", s.percent)

	percents := ""
	if s.percent > 0 {
		percents = styleProgressActive.Render(pb)
	} else {
		percents = styleProgressWait.Render(pb)
	}

	if active {
		return []byte(percents + styleStatusActive.Render(text) + "\n")
	}

	return []byte(percents + styleStatusWait.Render(text) + "\n")
}

func (t *Terminal) outputLine(line string) []byte {
	if len(line) > maxLineLength {
		line = string(line[:maxLineLength-2]) + ".."
	}

	return []byte(styleStatusActive.Render("  | ") + styleLogs.Render(line) + "\n")
}

func (t *Terminal) Init() tea.Cmd {
	return nil
}

func (t *Terminal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			t.finished = true
			return t, tea.Quit
		}
	}

	return t, func() tea.Msg {
		return true
	}
}

func (t *Terminal) View() string {
	if t == nil {
		return ""
	}

	return string(t.lastMessage)
}
