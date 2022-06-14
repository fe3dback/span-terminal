package terminal

type (
	StartOpt func(*Span)
)

func WithTitle(title string) StartOpt {
	return func(span *Span) {
		span.title = title
	}
}

func WithInitialProgress(initialProgress float64) StartOpt {
	return func(span *Span) {
		if initialProgress < 0 {
			initialProgress = 0
		}

		if initialProgress > 1 {
			initialProgress = 1
		}

		span.progress = int(initialProgress * 100)
	}
}
