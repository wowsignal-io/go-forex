// Package rba provides foreign exchange rates from the Reserve Bank of
// Australia.
//
// The data go back to January 2018. Rates are available from AUD to 21 other
// currencies, including 3 that don't appear on the ECB list and 2 that don't
// appear on the COB or ECB lists. (Consult currencies.txt for the full list.)
package rba

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/wowsignal-io/go-forex/forex/exchange"
	"github.com/wowsignal-io/go-forex/forex/internal"
)

const DefaultRBASource = "https://www.rba.gov.au/statistics/tables/csv/f11.1-data.csv"

// The line where currency units are named
const currenciesLine = 5

// The first line where exchange rate data are
const firstDataLine = 11

func parseHeader(cr *csv.Reader) (map[int]string, error) {

	err := internal.SkipLinesCSV(cr, currenciesLine)
	if err != nil {
		return nil, err
	}
	record, err := cr.Read()
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	if err != nil {
		return nil, err
	}

	result := make(map[int]string)
	for i, value := range record {
		if len(value) == 3 {
			// Currency symbols are three letters. Some other units are
			// indices that are of a different length and we ignore them.
			result[i] = value
		}
	}
	return result, nil
}

func parse(r io.Reader) ([]exchange.Rate, error) {
	cr := csv.NewReader(r)
	cr.ReuseRecord = true

	header, err := parseHeader(cr)
	if err != nil {
		return nil, err
	}

	internal.SkipLinesCSV(cr, firstDataLine-currenciesLine)
	result := []exchange.Rate{}
	for {
		record, err := cr.Read()
		if err == io.EOF {
			return result, nil
		}
		if err != nil {
			return nil, err
		}

		t, err := time.Parse("02-Jan-2006", record[0])
		if err != nil {
			return nil, parseError(err, 0, cr)
		}
		t = t.UTC().Truncate(24 * time.Hour)

	ColumnLoop:
		for field := 1; field < len(record); field++ {
			currency := header[field]
			if currency == "" {
				// This field is not in a currency column.
				continue ColumnLoop
			}

			if record[field] == "" {
				// No data on this day.
				continue ColumnLoop
			}

			x, err := strconv.ParseFloat(record[field], 64)
			if err != nil {
				return nil, parseError(err, field, cr)
			}

			result = append(result, exchange.Rate{
				From: "AUD",
				To:   currency,
				Day:  t,
				Rate: x,
				Info: "RBA",
			})
		}
	}
}

func Get(uri string) ([]exchange.Rate, error) {
	raw, err := internal.Fetch(uri)
	if err != nil {
		return nil, err
	}
	return parse(bytes.NewReader(raw))
}

func parseError(err error, field int, cr *csv.Reader) error {
	line, column := cr.FieldPos(field)
	return fmt.Errorf("%w on line %d, column %d", err, line, column)
}
