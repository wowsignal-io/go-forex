// Package forex provides an easy to use, performant API to find historical
// currency conversion rates.
//
// Historical exchange rates for about 50 currencies are sourced from central
// banks and cached locally after the first request. Custom sources can be
// ingested via Exchange.AddSource().
//
// Two preconfigured exchanges are provided: LiveExchange() refreshes data from
// online sources, while OfflineExchange() only uses a smaller historical
// database of rates embeedded in the package, and is suitable for use in
// environments not connected to the internet.
//
// Central banks don't provide all exchange rates directly, and some must be
// computed using a third (and sometimes a fourth) currency as an intermediate
// step. The algorithm in this package always discovers the shortest path
// available - it doesn't attempt to find the best exchange rate.
//
// The runtime cost of queries grows logarithmically with the length of
// historical data and linearly with the number of currencies. A query on the
// full dataset takes between 5000 and 10,000 ns on modern hardware.
//
// The computed exchange rates are for informational purposes only - they are
// unlikely to be the same as the rates actually offered, but the difference
// should be tolerable for home finance applications.
package forex

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/wowsignal-io/go-forex/forex/boc"
	"github.com/wowsignal-io/go-forex/forex/ecb"
	"github.com/wowsignal-io/go-forex/forex/exchange"
	"github.com/wowsignal-io/go-forex/forex/internal"
	"github.com/wowsignal-io/go-forex/forex/offline"
	"github.com/wowsignal-io/go-forex/forex/rba"
)

var (
	defaultOnce     sync.Once
	defaultExchange *Exchange
)

// LiveExchange sources exchange rates from multiple online sources, refreshing
// about twice per day.
func LiveExchange() *Exchange {
	defaultOnce.Do(func() {
		defaultExchange = &Exchange{
			CacheLife: DefaultCacheLife,
			CacheDir:  DefaultCacheDir(),
		}
		defaultExchange.AddSource("ECB", ecb.DefaultECBSource, ecb.Get)
		defaultExchange.AddSource("RBA", rba.DefaultRBASource, rba.Get)
		defaultExchange.AddSource("BOC", boc.DefaultBOCSource, boc.Get)
	})

	return defaultExchange
}

var (
	offlineOnce     sync.Once
	offlineExchange *Exchange
)

// OfflineExchange sources rates from a single source embedded in this package.
// It does not have live rates, only historical ones.
func OfflineExchange() *Exchange {
	offlineOnce.Do(func() {
		offlineExchange = &Exchange{
			CacheLife: DefaultCacheLife,
			CacheDir:  DefaultCacheDir(),
		}

		offlineExchange.AddSource("ECB (offline)",
			"data:text/csv;base64,"+base64.StdEncoding.EncodeToString([]byte(offline.HistoricalECBRates)),
			ecb.Get)
	})

	return offlineExchange
}

const DefaultCacheLife = 12 * time.Hour

func DefaultCacheDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, ".forex")
}

type Exchange struct {
	CacheLife time.Duration
	CacheDir  string

	mu           sync.RWMutex
	graph        exchange.Graph
	sources      []rateSource
	lastDownload time.Time
}

func (e *Exchange) String() string {
	sources := make([]string, len(e.sources))
	for i, s := range e.sources {
		sources[i] = s.name
	}
	return fmt.Sprintf("Exchange(%s, %d currencies)", strings.Join(sources, ", "), len(e.graph))
}

// Returns the mtime of the oldest on-disk cache file. Must be called with at
// least a readlock, and will result in calls to stat(). Should only be used
// once - subsequently lastDownload will be non-zero, and this codepath can be
// avoided.
func (e *Exchange) oldestCache() (time.Time, error) {
	var oldest time.Time
	for _, s := range e.sources {
		t, err := s.lastReload()
		if err != nil {
			return time.Time{}, err
		}
		if t.Before(oldest) || oldest.IsZero() {
			oldest = t
		}
	}

	return oldest, nil
}

func (e *Exchange) lockedRead() (exchange.Graph, error) {
	e.mu.RLock()
	g := e.graph
	lastDownload := e.lastDownload
	e.mu.RUnlock()

	// The graph is never modified, only replaced. If we have a pointer to it,
	// it's safe to read without holding the lock. Do a check on our copied data
	// to see if we need a refresh.
	now := time.Now()
	if g == nil || lastDownload.Before(now.Add(e.CacheLife)) {
		var err error
		if g, err = e.maybeRefresh(now); err != nil {
			return nil, err
		}
	}
	return g, nil
}

// Convert computes the exchange rate between the from and to currencies on a
// given day.
func (e *Exchange) Convert(from, to string, t time.Time, opts ...exchange.Option) (exchange.Result, error) {
	g, err := e.lockedRead()
	if err != nil {
		return exchange.Result{}, err
	}

	return exchange.Convert(g, from, to, t, opts...)
}

// Currencies returns the available currencies as a map of strings (a set).
//
// Note that technically nothing guarantees all of these currencies are mutually
// interconvertible, but in practice, they always are, because all data sources
// are related to one of the major currencies, and all of them are
// interconvertible.
func (e *Exchange) Currencies() (map[string]bool, error) {
	g, err := e.lockedRead()
	if err != nil {
		return nil, err
	}

	res := make(map[string]bool, len(e.graph))
	for c := range g {
		res[c] = true
	}

	return res, nil
}

type Freshness int16

const (
	FromMemory Freshness = iota
	FromLocalCache
	FromRemoteSource
)

func (e *Exchange) maybeRefresh(t time.Time) (exchange.Graph, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Repeat the check that brought us here, this time holding the write lock.
	// This ensures contention doesn't cause multiple reloads in quick sequence,
	// and also lets us figure out the level of refresh required.

	lvl := FromMemory

	if e.graph == nil {
		// This is the first operation on a new Exchange. We need to check the
		// age of on-disk caches.
		t, err := e.oldestCache()
		if err != nil {
			return nil, err
		}
		e.lastDownload = t
		lvl = FromLocalCache
	}

	// Reload from source if the oldest cache is stale. (If we don't have
	// on-disk cache, then lastDownload will be set to year 0, which will count
	// as stale.)
	if e.lastDownload.Before(time.Now().Add(-e.CacheLife)) {
		lvl = FromRemoteSource
	}

	if err := e.forceRefresh(lvl); err != nil {
		return nil, err
	}

	// Return the graph while we're still holding the lock - this saves the call
	// site having to reacquire the read lock.
	return e.graph, nil
}

// ForceRefresh will rebuild the exchange graph from the upstream source, which
// may be online or otherwise remote to this machine.
func (e *Exchange) ForceRefresh() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.forceRefresh(FromRemoteSource)
}

func (e *Exchange) forceRefresh(lvl Freshness) error {
	if lvl == FromMemory {
		return nil
	}

	now := time.Now()
	var rates []exchange.Rate

	// Loading the sources can start downloads over the network, so it makes
	// sense to do it in parallel. (This appears to speed up a lot with multiple
	// sources.)
	var err error
	ch := make(chan []exchange.Rate)
	errCh := make(chan error)
	var wg sync.WaitGroup

	for _, s := range e.sources {
		wg.Add(1)
		s := s
		go func() {
			defer wg.Done()

			r, err := s.reload(now, lvl == FromRemoteSource, e.CacheLife)
			if err != nil {
				errCh <- err
			}
			ch <- r
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

Loop:
	for {
		select {
		case r, ok := <-ch:
			if !ok {
				break Loop
			}
			rates = append(rates, r...)
		case err2 := <-errCh:
			if err == nil {
				err = err2
			}
		}
	}

	// An error could have come from one of the load routines.
	if err != nil {
		return err
	}

	g, err := exchange.Compile(rates)
	if err != nil {
		return err
	}
	e.graph = g

	if lvl == FromRemoteSource {
		e.lastDownload = now
	}

	return nil
}

// GetFunc knows how to load and parse exchange rates from a URL. It can be used
// with AddSource to register a new source of exchange rates.
type GetFunc func(url string) ([]exchange.Rate, error)

var pathFriendlyChars = regexp.MustCompile(`[^a-zA-Z0-9]`)

// AddSource adds a new source of exchange rates. The caller must call
// ForceReload if the Exchange has been recently used and has a local cache.
func (e *Exchange) AddSource(name string, url string, getter GetFunc) {
	e.mu.Lock()
	defer e.mu.Unlock()

	cachePath := filepath.Join(e.CacheDir, "forex_"+pathFriendlyChars.ReplaceAllString(name, "_")+"_cache")
	e.sources = append(e.sources, rateSource{
		name:      name,
		cachePath: cachePath,
		sourceURL: url,
		f:         getter,
	})
}

type rateSource struct {
	name       string
	cachePath  string
	sourceURL  string
	f          GetFunc
	reloadTime time.Time
}

func (s *rateSource) lastReload() (time.Time, error) {
	if s.reloadTime.IsZero() {
		st, err := os.Stat(s.cachePath)
		if os.IsNotExist(err) {
			return time.Time{}, nil
		}

		if err != nil {
			return time.Time{}, err
		}

		s.reloadTime = st.ModTime()
	}

	return s.reloadTime, nil
}

func (s *rateSource) reload(now time.Time, download bool, ttl time.Duration) (rates []exchange.Rate, err error) {
	if download {
		// Ignore the error here - whether or not this worked, the thing that
		// matters is the os.Create call.
		os.MkdirAll(filepath.Dir(s.cachePath), 0740)

		f, err := os.Create(s.cachePath)
		if err != nil {
			return nil, err
		}

		// Best effort - lock the file on systems that support it. (This is
		// cooperative, but the only code that should be touching this file is
		// this code.) Systems that don't support flock (e.g. Windows) typically
		// coordinate file access more strongly than UNIX, so things should even
		// out.
		syscall.Flock(int(f.Fd()), syscall.LOCK_EX)
		defer func() {
			err2 := f.Sync()
			if err == nil {
				err = err2
			}
			syscall.Flock(int(f.Fd()), syscall.F_UNLCK)
			// After Sync, Close has no reason to return error, but strange
			// things do happen.
			err2 = f.Close()
			if err == nil {
				err = err2
			}
		}()
		data, err := internal.Fetch(s.sourceURL)
		if err != nil {
			return nil, err
		}
		if _, err := f.Write(data); err != nil {
			return nil, err
		}
	}

	return s.f(s.cachePath)
}