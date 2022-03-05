// Package exchange provides currency conversion based on historical data.
//
// Most applications should use package forex, which provides precompiled
// exchange rate data and convenient caching.
//
// Available rates must be passed to Compile to bake an Exchange graph. Convert
// is used to compute an exchange rate between any two currencies based on the
// available data.
//
// The algorithm is a simple breadth-first search to discover the shortest path
// to convert between two currencies. (E.g. the shortest path from CZK to AED
// might be  be CZK -> EUR -> AUD -> TWD). Query complexity grows linearly with
// the number of currencies and logarithmically with the number of days of
// historical data.
//
// The computed conversion rates are for informational purposes only - they are
// unlikely to be the same as the rates actually offered, but the difference
// should be tolerable for home finance applications.
package exchange

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

// ErrNotFound is returned when no exchange data is available to satisfy a
// query.
var ErrNotFound = errors.New("no forex data")

// Rate represents the conversion rate between two currencies on a given day.
type Rate struct {
	// The rate applies when converting from the From currency to the To
	// currency.
	From, To string
	// The effective rate.
	Rate float64
	// Valid on this day, in UTC.
	Day time.Time
	// Additional information about how the rate was sourced. Usually the name
	// of the central bank whose data was used.
	Info string
}

// Graph is a compiled graph of currencies connected by their conversion rates.
// Query using Convert, not directly.
type Graph map[string]*currency

type currency struct {
	symbol string
	// Must remain sorted from the most recent Rate.Day.
	rates []edge
}

type edge struct {
	src, dst *currency
	rate     float64
	day      time.Time
	info     string
	inverse  bool
}

// Compile produces a graph used for currency conversion.
func Compile(rates []Rate) (Graph, error) {
	m := map[string]*currency{}
	for _, rate := range rates {
		day := rate.Day.Truncate(24 * time.Hour)
		src, ok := m[rate.From]
		if !ok {
			src = &currency{symbol: rate.From}
			m[rate.From] = src
		}

		dst, ok := m[rate.To]
		if !ok {
			dst = &currency{symbol: rate.To}
			m[rate.To] = dst
		}

		src.rates = append(src.rates, edge{
			src:  src,
			dst:  dst,
			rate: rate.Rate,
			day:  day,
			info: rate.Info})

		dst.rates = append(dst.rates, edge{
			src:     dst,
			dst:     src,
			rate:    1 / rate.Rate,
			day:     day,
			info:    rate.Info + " (inverse)",
			inverse: true,
		})
	}

	for _, c := range m {
		sort.Slice(c.rates, func(i, j int) bool {
			return c.rates[i].day.After(c.rates[j].day)
		})
	}

	return m, nil
}

func filterEdges(edges []edge, day time.Time, tolerance time.Duration) []edge {
	i := sort.Search(len(edges), func(i int) bool { return !edges[i].day.After(day) })
	if i == len(edges) {
		return nil
	}

	edges = edges[i:]
	day = day.Add(-24 * time.Hour).Add(-tolerance)
	i = sort.Search(len(edges), func(i int) bool { return !edges[i].day.After(day) })

	if i == len(edges) {
		return edges
	}

	// The call site expects to modify the contents of this slice, so we must
	// return a copy.
	res := make([]edge, i)
	copy(res, edges[:i])

	return res
}

// Result is a computed currency conversion rate obtained from Convert.
type Result struct {
	// The computed rate. If a central bank published conversion rates for the
	// currency pair, this is that rate. For currency pairs with no published
	// rate, this is a computed rate, obtained by converting to an intermediate
	// currency first.
	Rate float64
	// The conversion trace. Trace[0] is the from currency, Trace[len(Trace)-1]
	// is the to currency. If Rate was computed through an intermediate currency
	// (or two), then len(Trace) will be 3 or 4.
	//
	// Only populated if Convert was called with the FullTrace option.
	Trace []Rate
}

// ResultType is an option for Convert. It specifies which fields of Result
// should be populated.
//
// The default value is RateOnly.
type ResultType int16

func (rt ResultType) apply(opts *options) {
	opts.resultType = rt
}

func (rt ResultType) String() string {
	switch rt {
	case RateOnly:
		return "RateOnly"
	case FullTrace:
		return "FullTrace"
	default:
		return "<invalid ResultType>"
	}
}

const (
	// Only populate Result.Rate.
	RateOnly ResultType = iota
	// Populate Result.Rate and Result.Trace.
	FullTrace
)

// Tolerance is an option for Convert. When exchange data is not available on
// the desired day, Tolerance specifies how many earlier days may be checked.
//
// The default value is 0 (exact match only).
type Tolerance time.Duration

func (t Tolerance) apply(opts *options) {
	opts.tolerance = time.Duration(t)
}

func (t Tolerance) String() string {
	return fmt.Sprintf("Tolerance(%d days)", t/Tolerance(time.Hour)/24)
}

func AcceptOlderRate(maxAgeDays int) Tolerance {
	return Tolerance(maxAgeDays) * 24 * Tolerance(time.Hour)
}

// Option for the Convert function. Specifies optional arguments, like whether
// to accept stale exchange rates. See the list of types that implement this
// interface for a list of options.
type Option interface {
	apply(*options)
}

type options struct {
	resultType ResultType
	tolerance  time.Duration
}

// Convert from the from currency to the to currency using the provided exchange
// graph. Only rates from the specified date will be used (but see Tolerance).
//
// Most users should use forex.Convert instead. The only reason to use this
// function is if the application wants finer control over exchange data and
// caching.
func Convert(exchange Graph, from, to string, t time.Time, opts ...Option) (Result, error) {
	// The exchange rate is a graph with possible cycles. Each edge is only
	// valid on a specific day, and the edges in each vertex are stored in
	// ascending order of day, enabling binary search.
	//
	// The search algorithm is a BFS using a slice as a FIFO queue.* Edges are
	// filtered by time: binary search determines the lowest offset for valid
	// edges in each vertex.
	//
	// As an added complication, we may want to keep track of the edges that
	// contributed to generating the resulting exchange rate. This is only done
	// if the ResultType parameter is set to FullTrace, because it requires
	// about 4 times more storage. If full trace is requested, the `trace` map
	// is used to keep track of which edge was used to visit each unique
	// currency. (This works, because the `seen` set prevents revisiting
	// currencies.)
	//
	// *: It's customary to use a linked list, but benchmarks in Go consistently
	// show slices performing better.

	var o options
	for _, opt := range opts {
		opt.apply(&o)
	}

	t = t.UTC().Truncate(24 * time.Hour)
	c := exchange[from]
	if c == nil {
		return Result{}, ErrNotFound
	}

	q := filterEdges(c.rates, t, o.tolerance)
	// What currencies have been visited in the QueueLoop
	seen := make(map[string]bool, len(exchange))
	// What edge was last seen per target currency by the RateLoop.
	//
	// TODO: this map seems to confuse heap escape analysis. (It adds about 10
	// allocs per lookup.)
	seenEdges := make(map[*currency]*edge, len(exchange))
	seen[from] = true
	var trace map[*currency]edge
	if o.resultType == FullTrace {
		trace = make(map[*currency]edge, len(exchange))
	}

QueueLoop:
	for len(q) > 0 {
		candidate := q[0]
		q = q[1:]

		if seen[candidate.dst.symbol] {
			continue QueueLoop
		}

		if candidate.dst.symbol == to {
			return finalize(candidate.rate, candidate, trace), nil
		}

		// Binary search over the available rates (egdes). The rates are sorted
		// by Day starting from the most recent. This finds the first possibly
		// valid rate.
		pred := func(i int) bool { return !candidate.dst.rates[i].day.After(t) }

	RateLoop:
		for i := sort.Search(len(candidate.dst.rates), pred); i < len(candidate.dst.rates); i++ {
			e := candidate.dst.rates[i]
			if t.Sub(e.day) > o.tolerance {
				// No rates found on the day, or within tolerance. Move on to
				// the next candidate in the BFS queue.
				break RateLoop
			}
			// Only process the most recent edge - don't check multiple days of
			// edges leading to the same currency.
			if seen[e.dst.symbol] || seenEdges[e.dst] == &candidate {
				continue RateLoop
			}
			seenEdges[e.dst] = &candidate

			// The edge is valid on this day - push it onto the queue. If we're
			// not keeping track of the full trace, then we also need to
			// calculate the rate as we go. Here the partial product gets stored
			// in the queue copy of each edge.
			if trace == nil {
				e.rate *= candidate.rate
			}
			q = append(q, e)
		}

		seen[candidate.dst.symbol] = true
		if trace != nil {
			trace[candidate.dst] = candidate
		}
	}

	return Result{}, ErrNotFound
}

func finalize(rate float64, e edge, trace map[*currency]edge) Result {
	if trace == nil {
		return Result{Rate: rate}
	}

	path := []Rate{}
	for {
		path = append(path, Rate{From: e.src.symbol, To: e.dst.symbol, Rate: e.rate, Day: e.day, Info: e.info})
		prev, ok := trace[e.src]
		if !ok {
			break
		}
		e = prev
		rate *= e.rate
	}

	// The trace is in the wrong order (going back to the start).
	for i := 0; i < len(path)/2; i++ {
		j := len(path) - i - 1
		path[i], path[j] = path[j], path[i]
	}

	return Result{Trace: path, Rate: rate}
}
