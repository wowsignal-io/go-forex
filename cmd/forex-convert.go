// forex-convert
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/wowsignal-io/go-forex/forex"
	"github.com/wowsignal-io/go-forex/forex/exchange"
)

var (
	from      = flag.String("from", "", "the currency to convert from (3-letter symbol)")
	to        = flag.String("to", "", "the currency to convert to (3-letter symbol)")
	tolerance = flag.Int("tolerance", 0, "how many days before the specified date to search for the forex rate")
	verbose   = flag.Bool("v", false, "print more info, mainly the conversion trace")
	offline   = flag.Bool("offline", false, "don't connect to the internet, use only offline data")
	date      = flag.String("date", "today", "effective date as YYYY-MM-DD, or aliases 'today' and 'yesterday'")
	debug     = flag.Bool("debug", false, "print additional debugging information to stderr")
)

func getDate() (time.Time, error) {
	if *date == "today" {
		return time.Now(), nil
	}
	if *date == "yesterday" {
		return time.Now().Add(-24 * time.Hour), nil
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

func getTolerance() int {
	if (*date == "today" || *date == "yesterday") && !flagProvided("tolerance") {
		return 3
	}
	return *tolerance
}

func getOpts() []exchange.Option {
	opts := []exchange.Option{
		exchange.AcceptOlderRate(getTolerance()),
	}

	opts = append(opts, exchange.FullTrace)

	return opts
}

func getExchange() *forex.Exchange {
	if *offline {
		return forex.OfflineExchange()
	}

	return forex.LiveExchange()
}

func flagUsage(f *flag.Flag) {
	fmt.Fprintf(flag.CommandLine.Output(), "  -%s", f.Name)
	if f.DefValue != "" {
		fmt.Fprintf(flag.CommandLine.Output(), " (default=%s)", f.DefValue)
	}
	fmt.Fprintf(flag.CommandLine.Output(), "\n\t%s\n", f.Usage)
}

func printUsage() {
	fmt.Fprint(flag.CommandLine.Output(), "Usage: forex-convert -from FROM -to TO")
	fmt.Fprint(flag.CommandLine.Output(), " [-date YYYY-MM-DD] [-tolerance TOLERANCE] [-offline] [-v]\n")
	fmt.Fprint(flag.CommandLine.Output(), "Options:\n")
	flagUsage(flag.Lookup("from"))
	flagUsage(flag.Lookup("to"))
	flagUsage(flag.Lookup("date"))
	flagUsage(flag.Lookup("tolerance"))
	flagUsage(flag.Lookup("offline"))
	flagUsage(flag.Lookup("v"))
}

func getCurrency(s string) (string, error) {
	if s == "" {
		return "", errors.New("must specify currency")
	}

	if len(s) != 3 {
		return "", fmt.Errorf("%q is not a valid 3-letter currency symbol", s)
	}

	return strings.ToUpper(s), nil
}

func main() {
	flag.Usage = printUsage
	flag.Parse()

	if *from == "" || *to == "" {
		flag.Usage()
		os.Exit(1)
	}

	t, err := getDate()
	if err != nil {
		log.Fatalf("Invalid date: %v", err)
	}
	t = t.UTC().Truncate(24 * time.Hour)
	e := getExchange()

	src, err := getCurrency(*from)
	if err != nil {
		log.Fatalf("Invalid -from value: %v", err)
	}

	dst, err := getCurrency(*to)
	if err != nil {
		log.Fatalf("Invalid -to value: %v", err)
	}

	rate, err := e.Convert(src, dst, t, getOpts()...)
	if *debug {
		log.Printf("Cache dir=%s lifetime=%v", e.CacheDir, e.CacheLife)
		log.Printf("Using exchange %v", e)
	}
	if err != nil {
		log.Fatalf("Convert: %v", err)
	}

	for _, step := range rate.Trace {
		if !step.Day.Equal(t) {
			log.Printf("Warning: rate %s to %s is stale, dated %s (wanted %s, -tolerance=%d)",
				step.From, step.To, step.Day.Format("2006-01-02"), t.Format("2006-01-02"), getTolerance())
		}
	}

	if *verbose {
		for i, step := range rate.Trace {
			fmt.Printf("Conversion step %d/%d: 1 %s = %f %s (source: %s on %v)\n",
				i+1, len(rate.Trace), step.From, step.Rate, step.To, step.Info, step.Day.Format("2006-01-02"))
		}
		fmt.Print("Computed rate: ")
	}
	fmt.Printf("%f\n", rate.Rate)
}
