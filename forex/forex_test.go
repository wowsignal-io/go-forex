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

func Example() {
	// Get the exchange rate for January 4, 2022, between the Papuan Kina
	// and the Indian Rupee. Of the date, only the day matters - the rest is
	// rounded off. (Exchange rates are published daily.)
	rate, err := LiveExchange().Convert("TWD", "CZK", time.Date(2023, time.February, 10, 0, 0, 0, 0, time.UTC), exchange.RateOnly)
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
	rate, _ = LiveExchange().Convert("TWD", "CZK", time.Date(2022, time.January, 4, 0, 0, 0, 0, time.UTC), exchange.FullTrace)

	// Passing exchange.FullTrace as the last argument means rate.Trace is now
	// populated.
	for i, step := range rate.Trace {
		fmt.Printf("Conversion step %d/%d: 1 %s = %f %s (source: %s)\n", i+1, len(rate.Trace), step.From, step.Rate, step.To, step.Info)
	}

	// The output shows that the exchange rate was generated from Bank of
	// Canada's data by first converting to the Canadian Dollar and then to the
	// Rupee.

	// Output:
	// The rate is 0.730520
	// Conversion step 1/3: 1 TWD = 0.046150 CAD (source: BOC)
	// Conversion step 2/3: 1 CAD = 0.697010 EUR (source: BOC (inverse))
	// Conversion step 3/3: 1 EUR = 24.745000 CZK (source: ECB)
}

// Simple example of how to convert between two currencies.
func ExampleExchange_Convert() {
	rate, err := LiveExchange().Convert("USD", "EUR", time.Date(2022, time.January, 4, 0, 0, 0, 0, time.UTC))
	if err != nil {
		// Handle error.
	}
	fmt.Printf("The conversion rate from USD to EUR on January 4, 2022 was %f\n", rate.Rate)
	// Output: The conversion rate from USD to EUR on January 4, 2022 was 0.886603
}

// Shows how to properly handle ErrNotFound.
func ExampleExchange_Convert_errNotFound() {
	rate, err := LiveExchange().Convert("USD", "EUR", time.Date(2022, time.January, 2, 0, 0, 0, 0, time.UTC))
	if errors.Is(err, exchange.ErrNotFound) {
		fmt.Printf("No data for a Sunday\n")
	} else if err != nil {
		fmt.Printf("Got unknown error %v\n", err)
	} else {
		fmt.Printf("The conversion rate from USD to EUR on January 2, 2022 was %f\n", rate.Rate)
	}
	// Output: No data for a Sunday
}

// Shows how to accept an older exchange rate, when no data is available.
func ExampleExchange_Convert_acceptOlderRate() {
	rate, err := LiveExchange().Convert("USD", "EUR", time.Date(2022, time.January, 2, 0, 0, 0, 0, time.UTC), exchange.AcceptOlderRate(7))
	if errors.Is(err, exchange.ErrNotFound) {
		fmt.Printf("No data within a week of January 2, 2022\n")
	}
	if err != nil {
		fmt.Printf("Got unknown error %v\n", err)
	}
	fmt.Printf("The conversion rate from USD to EUR in the week before January 2, 2022 was %f\n", rate.Rate)
	// Output: The conversion rate from USD to EUR in the week before January 2, 2022 was 0.882924
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
					{From: "USD", To: "EUR", Day: time.Date(2012, time.July, 19, 0, 0, 0, 0, time.UTC), Rate: 1.0 / 1.22},
					{From: "EUR", To: "CZK", Day: time.Date(2012, time.July, 19, 0, 0, 0, 0, time.UTC), Rate: 25.3},
				},
			},
		},
		{
			comment:  "four currencies full trace",
			from:     "TWD",
			to:       "CZK",
			day:      time.Date(2023, time.February, 10, 0, 0, 0, 0, time.UTC),
			exchange: LiveExchange(),
			opts:     []exchange.Option{exchange.FullTrace},

			want: exchange.Result{
				Rate: 0.73,
				Trace: []exchange.Rate{
					{From: "TWD", To: "CAD", Day: time.Date(2023, time.February, 10, 0, 0, 0, 0, time.UTC), Rate: 0.04429},
					{From: "CAD", To: "EUR", Day: time.Date(2023, time.February, 10, 0, 0, 0, 0, time.UTC), Rate: 0.696},
					{From: "EUR", To: "CZK", Day: time.Date(2023, time.February, 10, 0, 0, 0, 0, time.UTC), Rate: 23.69},
				},
			},
		},
		{
			comment:  "four currencies no trace",
			from:     "TWD",
			to:       "CZK",
			day:      time.Date(2022, time.February, 10, 0, 0, 0, 0, time.UTC),
			exchange: LiveExchange(),

			want: exchange.Result{
				Rate: 0.73,
			},
		},
		{
			comment:  "unknown currency",
			from:     "XXX",
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
					{From: "USD", To: "EUR", Day: time.Date(2012, time.July, 19, 0, 0, 0, 0, time.UTC), Rate: 1.0 / 1.22},
					{From: "EUR", To: "CZK", Day: time.Date(2012, time.July, 19, 0, 0, 0, 0, time.UTC), Rate: 25.3},
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
					{From: "USD", To: "EUR", Day: time.Date(2022, time.February, 11, 0, 0, 0, 0, time.UTC), Rate: 1.0 / 1.14},
					{From: "EUR", To: "CZK", Day: time.Date(2022, time.February, 11, 0, 0, 0, 0, time.UTC), Rate: 24.36},
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
