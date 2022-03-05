// forex-convert
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/wowsignal-io/go-forex/forex"
	"github.com/wowsignal-io/go-forex/forex/exchange"
)

var (
	from      = flag.String("from", "USD", "the currency to convert from")
	to        = flag.String("to", "EUR", "the currency to convert to")
	tolerance = flag.Int("tolerance", 0, "how many days before the specified date to search for the forex rate")
	verbose   = flag.Bool("v", false, "print more info, mainly the conversion trace")
	offline   = flag.Bool("offline", false, "")
	date      = flag.String("date", "", "effective date as YYYY-MM-DD - if empty, then use today and use a 3-day tolerance")
)

func getDate() (time.Time, error) {
	if *date == "" {
		return time.Now(), nil
	}

	return time.Parse("2006-01-02", *date)
}

func flagProvided(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func getOpts() []exchange.Option {
	tolerance := *tolerance
	if *date == "" && !flagProvided("tolerance") {
		tolerance = 3
	}

	opts := []exchange.Option{
		exchange.AcceptOlderRate(tolerance),
	}

	if *verbose {
		opts = append(opts, exchange.FullTrace)
	}

	return opts
}

func getExchange() *forex.Exchange {
	if *offline {
		return forex.OfflineExchange()
	}

	return forex.LiveExchange()
}

func main() {
	flag.Parse()

	t, err := getDate()
	if err != nil {
		log.Fatalf("Invalid date: %v", err)
	}
	e := getExchange()

	rate, err := e.Convert(*from, *to, t, getOpts()...)
	if err != nil {
		log.Fatalf("Convert: %v", err)
	}

	if *verbose {
		for i, step := range rate.Trace {
			fmt.Printf("Conversion step %d/%d: 1 %s = %f %s (source: %s)\n", i+1, len(rate.Trace), step.From, step.Rate, step.To, step.Info)
		}
		fmt.Print("Computed rate: ")
	}
	fmt.Printf("%f\n", rate.Rate)
}
