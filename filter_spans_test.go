package terminal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mostRelevantSpans(t *testing.T) {
	span1 := &Span{id: 1, finished: false}
	span2 := &Span{id: 2, finished: true}
	span3 := &Span{id: 3, finished: false}
	span4 := &Span{id: 4, finished: true}
	spans := []*Span{span1, span2, span3, span4}

	tests := []struct {
		name  string
		limit int
		want  []*Span
	}{
		{
			name:  "all feat",
			limit: 4,
			want:  spans,
		},
		{
			name:  "limited to 3",
			limit: 3,
			want:  []*Span{span3, span1, span4},
		},
		{
			name:  "limited to 1",
			limit: 1,
			want:  []*Span{span3},
		},
	}
	for _, tt := range tests {
		got := mostRelevantSpans(spans, tt.limit)
		assert.Equal(t, tt.want, got)
	}
}
