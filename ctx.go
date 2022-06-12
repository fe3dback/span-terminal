package terminal

import "context"

type ctxSpan = struct{}

func spanFromContext(ctx context.Context) *Span {
	if span, ok := ctx.Value(ctxSpan{}).(*Span); ok {
		return span
	}

	return nil
}

func contextWithSpan(ctx context.Context, span *Span) context.Context {
	return context.WithValue(ctx, ctxSpan{}, span)
}
