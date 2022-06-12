package terminal

import "fmt"

type actionType int

const (
	actionTypeUnknown actionType = iota
	actionTypeMessage
	actionTypeSpanEnd
)

type (
	action interface {
		fmt.Stringer
		SpanID() spanID
		Span() *Span
		Type() actionType
	}
)

type (
	commonAction struct {
		span       *Span
		actionType actionType
		message    string
	}
)

func newAction(span *Span, actionType actionType, message string) *commonAction {
	return &commonAction{
		span:       span,
		actionType: actionType,
		message:    message,
	}
}

func (a *commonAction) String() string {
	return a.message
}

func (a *commonAction) Type() actionType {
	return a.actionType
}

func (a *commonAction) Span() *Span {
	return a.span
}

func (a *commonAction) SpanID() spanID {
	if a.span == nil {
		return 0
	}

	return a.span.id
}
