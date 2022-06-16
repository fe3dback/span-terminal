package terminal

const RenderOptDefaultSpanMaxRoots = 4
const RenderOptDefaultSpanMaxChild = 6
const RenderOptDefaultSpanMaxDetails = 12
const RenderOptDefaultProgressZeroLabel = "..."
const RenderOptDefaultLogsMaxLength = 80
const RenderOptDefaultLogsPrefix = "| "

type (
	renderOpts struct {
		spansMaxRoots     int
		spansMaxChild     int
		spansMaxDetails   int
		progressZeroLabel string
		logsMaxLength     int
		logsPrefix        string
	}

	RenderOptInitializer func(*renderOpts)
)

var defaultRenderOpts = renderOpts{
	spansMaxRoots:     RenderOptDefaultSpanMaxRoots,
	spansMaxChild:     RenderOptDefaultSpanMaxChild,
	spansMaxDetails:   RenderOptDefaultSpanMaxDetails,
	progressZeroLabel: RenderOptDefaultProgressZeroLabel,
	logsMaxLength:     RenderOptDefaultLogsMaxLength,
	logsPrefix:        RenderOptDefaultLogsPrefix,
}

// WithRenderOptSpanMaxRoots set maximum span tasks to display (1 level)
// default = RenderOptDefaultSpanMaxRoots
func WithRenderOptSpanMaxRoots(max int) RenderOptInitializer {
	return func(opts *renderOpts) {
		opts.spansMaxRoots = max
	}
}

// WithRenderOptSpanMaxChild set maximum span subtasks to display (2 level)
// default = RenderOptDefaultSpanMaxChild
func WithRenderOptSpanMaxChild(max int) RenderOptInitializer {
	return func(opts *renderOpts) {
		opts.spansMaxChild = max
	}
}

// WithRenderOptSpanMaxChild set maximum span subtasks details to display (3 level)
// default = RenderOptDefaultSpanMaxDetails
func WithRenderOptSpanMaxDetails(max int) RenderOptInitializer {
	return func(opts *renderOpts) {
		opts.spansMaxDetails = max
	}
}

// WithRenderOptProgressZeroLabel set default label for zero progress (0%)
// should be 3 chars long, or output will be bad
// default = RenderOptDefaultProgressZeroLabel
func WithRenderOptProgressZeroLabel(label string) RenderOptInitializer {
	return func(opts *renderOpts) {
		opts.progressZeroLabel = label
	}
}

// WithRenderOptLogsMaxLength set line length limit
// all rows with `length > max` - will be trimmed (from center)
// default = RenderOptDefaultLogsMaxLength
func WithRenderOptLogsMaxLength(max int) RenderOptInitializer {
	return func(opts *renderOpts) {
		opts.logsMaxLength = max
	}
}

// WithRenderOptLogsPrefix set prefix for log lines
// default = RenderOptDefaultLogsPrefix
func WithRenderOptLogsPrefix(prefix string) RenderOptInitializer {
	return func(opts *renderOpts) {
		opts.logsPrefix = prefix
	}
}
