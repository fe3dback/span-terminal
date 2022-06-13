package v2

import (
	"fmt"
	"strings"
	"time"
)

type renderOpts struct {
	maxRootSpans      int    // todo
	maxChild          int    // will filter only ended todo
	progressZeroLabel string // should be 3 chars long
}

type renderOptInitializer func(*renderOpts)

func renderSpanWithOptions(span *Span, optsInitializers ...renderOptInitializer) []byte {
	if span == nil {
		return nil
	}

	opts := renderOpts{
		maxRootSpans:      4,
		maxChild:          8,
		progressZeroLabel: "Zzz",
	}

	for _, initializer := range optsInitializers {
		initializer(&opts)
	}

	return []byte(renderSpan(span, opts))
}

func renderSpan(span *Span, opt renderOpts) string {
	switch span.depth {
	case 0:
		return renderSpanRoot(span, opt)
	case 1:
		return renderSpanSecond(span, opt)
	case 2:
		return renderSpanThird(span, opt)
	}

	return ""
}

func renderSpanRoot(span *Span, opt renderOpts) string {
	childContent := make([]byte, 0)
	for _, subSpan := range span.child {
		childContent = append(childContent, []byte(renderSpan(subSpan, opt)+"\n")...)
	}

	return "" +
		renderSpanStatusLine(span, opt) + "\n" +
		string(span.container.content()) + "\n" +
		string(childContent) + "\n"
}

func renderSpanSecond(span *Span, opt renderOpts) string {
	childContent := make([]byte, 0)
	for _, subSpan := range span.child {
		childContent = append(childContent, []byte(""+
			renderSpanPadding(span)+
			renderSpan(subSpan, opt)+"\n")...)
	}

	return "" +
		renderSpanStatusLine(span, opt) + "\n" +
		string(childContent) + "\n" +
		string(span.container.content()) + "\n"
}

func renderSpanThird(span *Span, _ renderOpts) string {
	return span.title + " | " + string(span.container.content()) + "\n"
}

func renderSpanStatusLine(span *Span, opt renderOpts) string {
	content := "" + "[ " + renderSpanProgress(span, opt) + " ] " + span.title

	if span.finished {
		return styleStatusDone.Render(content)
	}

	return styleStatusActive.Render(content)
}

func renderSpanProgress(span *Span, opt renderOpts) string {
	if span.finished {
		return fmt.Sprintf("%5s", renderDuration(span.startAt, span.endAt))
	}

	if span.progress == 0 {
		return fmt.Sprintf(" %s ", opt.progressZeroLabel)
	}

	return fmt.Sprintf(" %2d%% ", span.progress)
}

func renderDuration(from, to time.Time) string {
	took := to.Sub(from)

	if took.Hours() > 1 {
		return fmt.Sprintf("%.0fh", took.Hours())
	}

	if took.Minutes() > 5 {
		return fmt.Sprintf("%.0fm", took.Minutes())
	}

	if took.Seconds() > 1 {
		return fmt.Sprintf("%.0fs", took.Seconds())
	}

	return fmt.Sprintf("%dms", took.Milliseconds())
}

func renderSpanPadding(span *Span) string {
	//  |- Third span of A 1 | -~ -~ -~ -~ -~ -~ -~ -~ -~
	return strings.Repeat(" ", int(span.depth)) + "| "
}
