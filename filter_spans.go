package terminal

import "sort"

const (
	spanBestA spanBest = iota
	spanBestB
	spanEqual
)

type spanBest = int
type spanFilter = func(a, b *Span) spanBest
type spanFilters = []spanFilter

var spanRelevantFilters = spanFilters{
	filterByFinished(),
	filterById(),
}

func mostRelevantSpans(spans []*Span, limit int) []*Span {
	if len(spans) <= limit {
		return spans
	}

	// most relevant
	best := spans
	sort.Slice(best, func(i, j int) bool {
		for _, filter := range spanRelevantFilters {
			switch filter(best[i], best[j]) {
			case spanEqual:
				continue
			case spanBestA:
				return true
			case spanBestB:
				return false
			}
		}

		return false
	})

	// get most relevant
	best = best[:limit]

	// sort by id
	sort.Slice(best, func(i, j int) bool {
		return best[i].id < best[j].id
	})

	return best
}

// Finished is less relevant for display
func filterByFinished() spanFilter {
	return func(a, b *Span) spanBest {
		if a.finished && !b.finished {
			return spanBestB
		}

		if !a.finished && b.finished {
			return spanBestA
		}

		return spanEqual
	}
}

// Older is less relevant
func filterById() spanFilter {
	return func(a, b *Span) spanBest {
		if a.id > b.id {
			return spanBestA
		}

		if b.id > a.id {
			return spanBestB
		}

		return spanEqual
	}
}
