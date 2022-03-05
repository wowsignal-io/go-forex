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
	tolerance = flag.Int("tolerance", 0, "")
	verbose   = flag.Bool("v", false, "print more info, mainly the conversion trace")
	offline   = flag.Bool("offline", false, "")
	date      = flag.String("date", "2022-01-04", "effective date as YYYY-MM-DD")
)

func main() {
	flag.Parse()

	opts := []exchange.Option{
		exchange.Tolerance(*tolerance),
	}

	if *verbose {
		opts = append(opts, exchange.FullTrace)
	}

	t, err := time.Parse("2006-01-02", *date)
	if err != nil {
		log.Fatalf("Invalid date: %v", err)
	}

	e := forex.LiveExchange()
	if *offline {
		e = forex.OfflineExchange()
	}
	rate, err := e.Convert(*from, *to, t, opts...)
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
