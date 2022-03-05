package internal

import (
	"fmt"

	"github.com/wowsignal-io/go-forex/forex/exchange"
)

func ValidateAll(rates []exchange.Rate, allowed map[string]bool, callback func(i int, warnings []string)) (notFound map[string]struct{}) {
	seen := map[string]bool{}
	for k := range allowed {
		seen[k] = false
	}

	for i, rate := range rates {
		if valid, warnings := Validate(rate, allowed); !valid {
			callback(i, warnings)
		}
		seen[rate.From] = true
		seen[rate.To] = true
	}

	notFound = map[string]struct{}{}
	for k, v := range seen {
		if !v {
			notFound[k] = struct{}{}
		}
	}

	return notFound
}

func Validate(r exchange.Rate, allowed map[string]bool) (bool, []string) {
	var warnings []string

	if r.Day.IsZero() {
		warnings = append(warnings, "zero Day value")
	}

	if r.Rate == 0 {
		warnings = append(warnings, "zero Rate value")
	}

	if r.From == r.To {
		warnings = append(warnings, "source and target currency are the same")
	}

	if r.From == "" {
		warnings = append(warnings, "missing source currency")
	}

	if r.To == "" {
		warnings = append(warnings, "missing target currency in rate")
	}

	if !allowed[r.From] {
		warnings = append(warnings, fmt.Sprintf("source currency %q not allowed here", r.From))
	}

	if !allowed[r.To] {
		warnings = append(warnings, fmt.Sprintf("target currency %q not allowed here", r.To))
	}

	return len(warnings) == 0, warnings
}
