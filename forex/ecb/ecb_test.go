package ecb

import (
	"testing"

	"github.com/wowsignal-io/go-forex/forex/internal"
)

func TestGet(t *testing.T) {
	rates, err := Get("testdata/eurofxref-hist.zip")
	if err != nil {
		t.Fatal(err)
	}

	const expectRateCount = 41 /* columns */ *5916 /* rows */ - 57188 /* number of N/A cells */
	if len(rates) != expectRateCount {
		t.Errorf("Found %d rates (expected %d)", len(rates), expectRateCount)
	}

	wantCurrencies, err := internal.Uniq("currencies.txt")
	if err != nil {
		t.Fatal(err)
	}

	notFound := internal.ValidateAll(rates, wantCurrencies, func(i int, warnings []string) {
		for _, warning := range warnings {
			t.Errorf("Rate %d/%d invalid: %s", i+1, len(rates), warning)
		}
	})

	for currency := range notFound {
		t.Errorf("Currency %s declared in currencies.txt, but not found in the output rates", currency)
	}
}
