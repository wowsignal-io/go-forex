package forex

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/wowsignal-io/go-forex/forex/exchange"
	"github.com/wowsignal-io/go-forex/forex/internal"
)

func ExampleLiveExchange() {
	// Get the exchange rate for January 4, 2022, between the Papuan Kina
	// and the Indian Rupee. Of the date, only the day matters - the rest is
	// rounded off. (Exchange rates are published daily.)
	rate, err := LiveExchange().Convert("PGK", "INR", time.Date(2022, time.January, 4, 0, 0, 0, 0, time.UTC), exchange.RateOnly)
	if err != nil {
		fmt.Printf("Convert failed (%v): are you connected to the internet?", err)
	}
	fmt.Printf("The rate is %f\n", rate.Rate)

	// We can also request more details about how the exchange rate was
	// obtained.
	//
	// Note that repeated calls to Convert don't download exchange rate data
	// from the internet more than once every 12 hours, even if the program exits and
	// restarts.
	rate, _ = LiveExchange().Convert("PGK", "INR", time.Date(2022, time.January, 4, 0, 0, 0, 0, time.UTC), exchange.FullTrace)

	// Passing exchange.FullTrace as the last argument means rate.Trace is now
	// populated.
	for i, step := range rate.Trace {
		fmt.Printf("Conversion step %d/%d: 1 %s = %f %s (source: %s)\n", i+1, len(rate.Trace), step.Src, step.Rate, step.Dst, step.Info)
	}

	// The output shows that the exchange rate was generated from Bank of
	// Canada's data by first converting to the Canadian Dollar and then to the
	// Rupee.

	// Output:
	// The rate is 21.255237
	// Conversion step 1/2: 1 PGK = 0.395226 AUD (source: RBA (inverse))
	// Conversion step 2/2: 1 AUD = 53.780000 INR (source: RBA)
}

func TestLiveExchange(t *testing.T) {
	want, err := internal.Uniq("currencies.txt")
	if err != nil {
		t.Fatal(err)
	}

	got, err := LiveExchange().Currencies()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Supported currencies (live exchange) -> (-) wanted vs. (+) got:\n%s", diff)
	}
}

func TestConvert(t *testing.T) {
	for _, tc := range []struct {
		comment  string
		from, to string
		day      time.Time
		exchange *Exchange
		opts     []exchange.Option

		want    exchange.Result
		wantErr error
	}{
		{
			comment:  "live full trace",
			from:     "USD",
			to:       "CZK",
			day:      time.Date(2012, time.July, 19, 0, 0, 0, 0, time.UTC),
			exchange: LiveExchange(),
			opts:     []exchange.Option{exchange.FullTrace},

			want: exchange.Result{
				Rate: 20.5895,
				Trace: []exchange.Rate{
					{Src: "USD", Dst: "EUR", Day: time.Date(2012, time.July, 19, 0, 0, 0, 0, time.UTC), Rate: 1.0 / 1.22},
					{Src: "EUR", Dst: "CZK", Day: time.Date(2012, time.July, 19, 0, 0, 0, 0, time.UTC), Rate: 25.3},
				},
			},
		},
		{
			comment:  "four currencies full trace",
			from:     "PGK",
			to:       "CZK",
			day:      time.Date(2022, time.February, 10, 0, 0, 0, 0, time.UTC),
			exchange: LiveExchange(),
			opts:     []exchange.Option{exchange.FullTrace},

			want: exchange.Result{
				Rate: 6.08,
				Trace: []exchange.Rate{
					{Src: "PGK", Dst: "AUD", Day: time.Date(2022, time.February, 10, 0, 0, 0, 0, time.UTC), Rate: 0.397},
					{Src: "AUD", Dst: "EUR", Day: time.Date(2022, time.February, 10, 0, 0, 0, 0, time.UTC), Rate: 0.63},
					{Src: "EUR", Dst: "CZK", Day: time.Date(2022, time.February, 10, 0, 0, 0, 0, time.UTC), Rate: 24.35},
				},
			},
		},
		{
			comment:  "four currencies no trace",
			from:     "PGK",
			to:       "CZK",
			day:      time.Date(2022, time.February, 10, 0, 0, 0, 0, time.UTC),
			exchange: LiveExchange(),

			want: exchange.Result{
				Rate: 6.08,
			},
		},
		{
			comment:  "unknown currency",
			from:     "TWD",
			to:       "CZK",
			day:      time.Date(2022, time.February, 10, 0, 0, 0, 0, time.UTC),
			exchange: OfflineExchange(),

			wantErr: exchange.ErrNotFound,
		},
		{
			comment:  "offline full trace",
			from:     "USD",
			to:       "CZK",
			day:      time.Date(2012, time.July, 19, 0, 0, 0, 0, time.UTC),
			exchange: OfflineExchange(),
			opts:     []exchange.Option{exchange.FullTrace},

			want: exchange.Result{
				Rate: 20.5895,
				Trace: []exchange.Rate{
					{Src: "USD", Dst: "EUR", Day: time.Date(2012, time.July, 19, 0, 0, 0, 0, time.UTC), Rate: 1.0 / 1.22},
					{Src: "EUR", Dst: "CZK", Day: time.Date(2012, time.July, 19, 0, 0, 0, 0, time.UTC), Rate: 25.3},
				},
			},
		},
		{
			comment:  "no rate on sunday",
			from:     "USD",
			to:       "CZK",
			day:      time.Date(2022, time.February, 13, 0, 0, 0, 0, time.UTC),
			exchange: LiveExchange(),
			opts:     []exchange.Option{exchange.FullTrace},

			wantErr: exchange.ErrNotFound,
		},
		{
			comment:  "no rate on sunday, use most recent previous",
			from:     "USD",
			to:       "CZK",
			day:      time.Date(2022, time.February, 13, 0, 0, 0, 0, time.UTC),
			exchange: LiveExchange(),
			opts:     []exchange.Option{exchange.FullTrace, exchange.AcceptOlderRate(5)},

			want: exchange.Result{
				Rate: 21.33,
				Trace: []exchange.Rate{
					{Src: "USD", Dst: "EUR", Day: time.Date(2022, time.February, 11, 0, 0, 0, 0, time.UTC), Rate: 1.0 / 1.14},
					{Src: "EUR", Dst: "CZK", Day: time.Date(2022, time.February, 11, 0, 0, 0, 0, time.UTC), Rate: 24.36},
				},
			},
		},
		{
			comment:  "no rate on sunday, use previous, tolerance too short",
			from:     "USD",
			to:       "CZK",
			day:      time.Date(2022, time.February, 13, 0, 0, 0, 0, time.UTC),
			exchange: LiveExchange(),
			opts:     []exchange.Option{exchange.FullTrace, exchange.AcceptOlderRate(1)},

			wantErr: exchange.ErrNotFound,
		},
	} {
		t.Run(tc.comment, func(t *testing.T) {
			result, err := tc.exchange.Convert(tc.from, tc.to, tc.day, tc.opts...)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("%v.Convert(%q, %q, %v, %v) -> error %v (wanted error %v)", tc.exchange, tc.from, tc.to, tc.day, tc.opts, err, tc.wantErr)
			}

			if diff := cmp.Diff(tc.want, result, cmpopts.EquateApprox(0, 0.05), cmpopts.IgnoreFields(exchange.Rate{}, "Info")); diff != "" {
				t.Errorf("%v.Convert(%q, %q, %v, %v) -> (-) wanted vs. (+) got:\n%s", tc.exchange, tc.from, tc.to, tc.day, tc.opts, diff)
			}
		})
	}

}

func BenchmarkConvertRateOnly(b *testing.B) {
	// Warm up the caches.
	e := LiveExchange()
	_, err := e.Convert("USD", "CZK", time.Date(2012, time.July, 19, 0, 0, 0, 0, time.UTC), exchange.RateOnly)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Convert("USD", "CZK", time.Date(2012, time.July, 19, 0, 0, 0, 0, time.UTC), exchange.RateOnly)
	}
}

func BenchmarkConvertFullTrace(b *testing.B) {
	// Warm up the caches.
	e := LiveExchange()
	_, err := e.Convert("USD", "CZK", time.Date(2012, time.July, 19, 0, 0, 0, 0, time.UTC), exchange.FullTrace)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Convert("USD", "CZK", time.Date(2012, time.July, 19, 0, 0, 0, 0, time.UTC), exchange.FullTrace)
	}
}
