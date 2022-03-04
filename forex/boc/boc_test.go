package boc

import (
	"testing"

	"github.com/wowsignal-io/go-forex/forex/internal"
)

func TestGet(t *testing.T) {
	rates, err := Get("testdata/FX_RATES_DAILY-sd-2017-01-03-2.csv")
	if err != nil {
		t.Fatal(err)
	}

	// number of valid floats in the CSV file determined by
	// `grep -o -P '\d+\.\d+' FX_RATES_DAILY-sd-2017-01-03-2.csv | wc -l``
	const expectRateCount = 31549
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
