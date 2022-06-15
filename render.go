package terminal

import (
	"fmt"
	"strings"
	"time"
)

type renderOpts struct {
	maxRootSpans      int    // todo
	maxChild          int    // will filter only ended todo
	progressZeroLabel string // should be 3 chars long
	logsMaxLength     int
	logsPrefix        string
}

type renderOptInitializer func(*renderOpts)

func renderSpanWithOptions(span *Span, optsInitializers ...renderOptInitializer) string {
	if span == nil {
		return ""
	}

	opts := &renderOpts{
		maxRootSpans:      4,
		maxChild:          8,
		progressZeroLabel: "...",
		logsMaxLength:     80,
		logsPrefix:        "| ",
	}

	for _, initializer := range optsInitializers {
		initializer(opts)
	}

	return renderSpan(span, opts)
}

func renderSpan(span *Span, opt *renderOpts) string {
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

func renderSpanRoot(span *Span, opt *renderOpts) string {
	childContent := ""
	for _, subSpan := range span.child {
		childContent += "" +
			renderSpanPadding(span) +
			renderSpan(subSpan, opt)
	}

	logs := ""
	if !span.finished {
		logs = renderContainer(span.container, opt) + "\n"
	}

	return "" +
		styleHeader.Render("[+] "+span.title) + "\n" +
		logs +
		childContent + "\n"
}

func renderSpanSecond(span *Span, opt *renderOpts) string {
	if span.finished {
		return renderSpanStatusLine(span, opt) + "\n"
	}

	childContent := ""
	for _, subSpan := range span.child {
		childContent += "" +
			renderSpanPadding(span) +
			renderSpan(subSpan, opt)
	}

	return "" +
		renderSpanStatusLine(span, opt) + "\n" +
		childContent + "\n"
}

func renderSpanThird(span *Span, opt *renderOpts) string {
	return renderSpanStatusLine(span, opt) + "\n"
}

func renderSpanStatusLine(span *Span, opt *renderOpts) string {
	prefix := " "
	delimiter := " "

	if span.depth.isDetailsOrDeeper3nd() {
		delimiter = " | "
	}

	if span.depth.isOperation2nd() {
		prefix = ">"
	}

	content := prefix + renderSpanProgress(span, opt) + delimiter + span.title

	if span.finished {
		return styleStatusDone.Render(content)
	}

	return styleStatusActive.Render(content)
}

func renderSpanProgress(span *Span, opt *renderOpts) string {
	if span.finished {
		return fmt.Sprintf("%5s", renderDuration(span.startAt, span.endAt))
	}

	if span.progress == 0 {
		return fmt.Sprintf("  %s", opt.progressZeroLabel)
	}

	return fmt.Sprintf("  %2d%%", span.progress)
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
	return strings.Repeat(" ", int(span.depth)) + " "
}

func renderContainer(c container, opt *renderOpts) string {
	logs := ""

	for _, line := range c.content() {
		if len(line) > opt.logsMaxLength {
			line = line[:opt.logsMaxLength]
		}

		logs += opt.logsPrefix + line + "\n"
	}

	return styleLogs.Render(logs)
}

func renderMainContainer(c container) string {
	return renderContainer(c, &renderOpts{
		logsMaxLength: 80,
		logsPrefix:    "",
	})
}
