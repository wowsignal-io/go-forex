// Package boc provides foreign exchange rates from Bank of Canada.
//
// By default, the data go back to January 2017. Rates are available from CAD
// to 26 other currencies, including 4 that don't appear on the ECB list.
// (Consult currencies.txt for the full list.)
package boc

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/wowsignal-io/go-forex/forex/exchange"
	"github.com/wowsignal-io/go-forex/forex/internal"
)

const DefaultBOCSource = "https://www.bankofcanada.ca/valet/observations/group/FX_RATES_DAILY/csv?start_date=2017-01-03"

func Get(uri string) ([]exchange.Rate, error) {
	raw, err := internal.Fetch(uri)
	if err != nil {
		return nil, err
	}

	needle := []byte("\"OBSERVATIONS\"")
	i := bytes.Index(raw, needle)
	if i < 0 {
		return nil, errors.New("invalid BOC sheet")
	}
	return parse(bytes.NewReader(raw[i+len(needle):]))
}

func parse(r io.Reader) ([]exchange.Rate, error) {
	cr := csv.NewReader(r)
	cr.ReuseRecord = true
	header, err := parseHeader(cr)
	if err != nil {
		return nil, err
	}

	result := []exchange.Rate{}
	for {
		record, err := cr.Read()
		if err == io.EOF {
			return result, nil
		}

		t, err := time.Parse("2006-01-02", record[0])
		if err != nil {
			return nil, parseError(err, 0, cr)
		}
		t = t.UTC().Truncate(24 * time.Hour)

	ColumnLoop:
		for field := 1; field < len(record); field++ {
			value := record[field]
			if value == "" {
				// No rate for this day.
				continue ColumnLoop
			}

			currency := header[field]
			if currency == "" {
				return nil, parseError(errors.New("don't know the currency"), field, cr)
			}

			rate, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, parseError(err, field, cr)
			}

			result = append(result, exchange.Rate{
				From: currency,
				Day:  t,
				To:   "CAD",
				Rate: rate,
				Info: "BOC",
			})
		}
	}
}

func parseHeader(cr *csv.Reader) (currencies map[int]string, err error) {
	record, err := cr.Read()
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	if err != nil {
		return nil, err
	}

	currencies = make(map[int]string)
	for i := 1; i < len(record); i++ {
		c, err := parseCurrencySymbol(record[i])
		if err != nil {
			return nil, err
		}
		currencies[i] = c
	}
	return currencies, nil
}

func parseError(err error, field int, cr *csv.Reader) error {
	line, column := cr.FieldPos(field)
	return fmt.Errorf("%w on line %d, column %d", err, line, column)
}

func parseCurrencySymbol(s string) (string, error) {
	if len(s) != 8 {
		return "", fmt.Errorf("currency symbol %q is is not 8 bytes", s)
	}

	if !strings.HasPrefix(s, "FX") {
		return "", fmt.Errorf("currency symbol %q does not start with FX", s)
	}

	if !strings.HasSuffix(s, "CAD") {
		return "", fmt.Errorf("currency symbol %q does not end in CAD", s)
	}

	return s[2:5], nil
}
