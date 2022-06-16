package terminal

const OptDefaultContainerMaxLines = 4
const OptDefaultStdoutMaxLines = 8

type (
	terminalOpts = struct {
		containerMaxLines int
		stdoutMaxLines    int
		renderOpts        renderOpts
	}

	OptsInitializer = func(*terminalOpts)
)

// WithContainerMaxLines set max log lines for each root span
// default = OptDefaultContainerMaxLines
func WithContainerMaxLines(maxLines int) OptsInitializer {
	return func(opt *terminalOpts) {
		opt.containerMaxLines = maxLines
	}
}

// WithStdoutMaxLines set max lines for captured stdout
// all fmt.* and log.* functions will print to this area
// default = OptDefaultStdoutMaxLines
func WithStdoutMaxLines(maxLines int) OptsInitializer {
	return func(opt *terminalOpts) {
		opt.stdoutMaxLines = maxLines
	}
}

// WithRenderOpts allow to customize spans printing
func WithRenderOpts(initializers ...RenderOptInitializer) OptsInitializer {
	return func(opts *terminalOpts) {
		renderOpts := defaultRenderOpts

		for _, initializer := range initializers {
			initializer(&renderOpts)
		}

		opts.renderOpts = renderOpts
	}
}
