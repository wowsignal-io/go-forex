// Package cbuae provides foreign exchange rates from the UAE central bank.
//
// Historical rates are published monthly as excel spreadsheets. Daily rates are
// available as HTML from a fairly convenient URL.
//
// At the moment, we don't implement historical rates - instead, we just grab
// the last three days of daily rates.
package cbuae

import (
	"bytes"
	_ "embed"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/wowsignal-io/go-forex/forex/exchange"
	"github.com/wowsignal-io/go-forex/forex/internal"
)

//go:embed currencies.txt
var currenciesTxt string

//go:embed currency_names.txt
var currencyNamesTxt string

func DownloadOption(req *http.Request, client *http.Client) *http.Request {
	req2, err := http.NewRequest("POST", req.URL.String(), nil)
	if err != nil {
		panic(err)
	}

	req2.Host = "www.centralbank.ae"
	req2.Header.Set("Accept", "text/html, */*; q=0.01")
	req2.Header.Set("Sec-Fetch-Site", "same-origin")
	req2.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req2.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req2.Header.Set("Sec-Fetch-Mode", "cors")
	req2.Header.Set("Origin", "https://www.centralbank.ae")
	req2.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15")
	req2.Header.Set("Referer", "https://www.centralbank.ae/en/forex-eibor/exchange-rates/")
	req2.Header.Set("Content-Length", "0")
	req2.Header.Set("Connection", "keep-alive")
	req2.Header.Set("Sec-Fetch-Dest", "empty")
	req2.Header.Set("X-Requested-With", "XMLHttpRequest")
	return req2
}

func nameToISOMap() map[string]string {
	m := map[string]string{}
	names := bytes.Split([]byte(currencyNamesTxt), []byte("\n"))
	symbols := bytes.Split([]byte(currenciesTxt), []byte("\n"))
	if len(names) != len(symbols) {
		panic(fmt.Sprintf("mismatched currency names and symbols: %d vs %d", len(names), len(symbols)))
	}
	for i, name := range names {
		m[string(name)] = string(symbols[i])
	}
	return m
}

func SourceURLForDate(date time.Time) string {
	return fmt.Sprintf("https://www.centralbank.ae/umbraco/Surface/Exchange/GetExchangeRateAllCurrencyDate?dateTime=%s", date.Format("2006-01-02"))
}

func Get(uri string) ([]exchange.Rate, error) {
	raw, err := internal.Fetch(uri, DownloadOption)
	if err != nil {
		return nil, err
	}

	return parse(raw)
}

func parseDate(raw []byte) (time.Time, error) {
	const needle = "Last updated:"
	const endNeedle = "</p>"
	idx := bytes.Index(raw, []byte(needle))
	if idx < 0 {
		return time.Time{}, fmt.Errorf("date not found")
	}
	raw = raw[idx+len(needle):]
	idx = bytes.Index(raw, []byte(endNeedle))
	if idx < 0 {
		return time.Time{}, fmt.Errorf("date not found")
	}
	raw = raw[:idx]
	raw = bytes.TrimSpace(raw)
	loc, err := time.LoadLocation("Asia/Dubai")
	if err != nil {
		panic(err)
	}
	return time.ParseInLocation("Monday 02 January 2006 03:04:05 PM", string(raw), loc)
}

func parse(raw []byte) ([]exchange.Rate, error) {
	const needleStartCell = `<td class="font-r fs-small text-navy-custom">`
	const needleEndCell = `</td>`
	namesToISO := nameToISOMap()
	date, err := parseDate(raw)
	if err != nil {
		return nil, err
	}

	rates := []exchange.Rate{}
	var rate *exchange.Rate
	parseCell := func(cell []byte) error {
		// Three options: currency name, rate, or empty cell.

		// Skip empty cells.
		if len(cell) == 0 {
			return nil
		}

		// If it's a valid rate, store it and flush the cell.
		f, err := strconv.ParseFloat(string(cell), 64)
		if err == nil {
			if rate != nil {
				rate.Rate = f
				rates = append(rates, *rate)
				rate = nil
			}
			return nil
		}

		// We have a new currency. Start the rate struct and wait for the
		// number next.
		iso, ok := namesToISO[string(cell)]
		if !ok {
			// Skip, on purpose.
			return nil
		}
		rate = &exchange.Rate{
			Info: "CBUAE",
			Day:  date,
			From: "AED",
			To:   iso,
		}
		return nil
	}

	state := "next_cell"
StateLoop:
	for {
		switch state {
		case "next_cell":
			idx := bytes.Index(raw, []byte(needleStartCell))
			if idx < 0 {
				break StateLoop
			}
			state = "read_cell"
			raw = raw[idx+len(needleStartCell):]
		case "read_cell":
			idx := bytes.Index(raw, []byte(needleEndCell))
			if idx < 0 {
				break StateLoop
			}
			cell := raw[:idx]
			state = "next_cell"
			if err := parseCell(cell); err != nil {
				return nil, err
			}
			raw = raw[idx:]
		}
	}

	return rates, nil
}
