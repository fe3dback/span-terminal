package terminal

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// how many logs can be buffered
const bufferCapacity = 8

type (
	Terminal struct {
		controlCtx context.Context
		target     *os.File

		buffer  actionChan
		sleepTo time.Time
		tea     *tea.Program

		sync.RWMutex
	}

	actionChan = chan action
)

func NewTerminal(controlCtx context.Context, target *os.File) *Terminal {
	t := &Terminal{
		controlCtx: controlCtx,
		target:     target,

		buffer:  make(actionChan, bufferCapacity),
		sleepTo: time.Now(),
	}

	renderer := tea.NewProgram(t)
	t.tea = renderer

	t.watch()
	return t
}

func (t *Terminal) Span(ctx context.Context, title string) (context.Context, *Span) {
	parent := spanFromContext(ctx)
	newSpan := newSpan(t, parent, title)

	return contextWithSpan(ctx, newSpan), newSpan
}

func (t *Terminal) watch() {
	finished := false

	err := t.tea.Start()
	if err != nil {
		fmt.Printf("failed init terminal: %v\n", err)
		return
	}

	for {
		if finished {
			break
		}

		select {
		case <-t.controlCtx.Done():
			finished = true
			break
		case act := <-t.buffer:
			t.rebuildView(act)
		}
	}

	t.tea.Quit()
}

func (t *Terminal) rebuildView(act action) {
	if act.Type() != actionTypeSpanEnd && time.Now().Before(t.sleepTo) {
		// do not update so often
		return
	}
}

func (t *Terminal) Init() tea.Cmd {
	return nil
}

func (t *Terminal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return t, nil
}

func (t *Terminal) View() string {
	// todo
	return fmt.Sprintf("time since start: %s", time.Since(t.sleepTo))
}
