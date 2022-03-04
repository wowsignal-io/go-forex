// Package ecb provides foreign exchange rates from the European Central Bank.
//
// The data goes back to January 1999, when the Euro was introduced. Rates are
// available from EUR to 41 other currencies. (See currencies.txt for the full
// list.)
//
//
package ecb

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"time"

	"github.com/wowsignal-io/go-forex/forex/exchange"
	"github.com/wowsignal-io/go-forex/forex/internal"
)

const DefaultECBSource = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist.zip"

func Get(uri string) ([]exchange.Rate, error) {
	raw, err := internal.Fetch(uri)
	if err != nil {
		return nil, err
	}

	rc, err := decompress(raw)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	return parse(rc)
}

func parseError(err error, field int, cr *csv.Reader) error {
	line, column := cr.FieldPos(field)
	return fmt.Errorf("%w on line %d, column %d", err, line, column)
}

func parse(r io.Reader) ([]exchange.Rate, error) {
	result := []exchange.Rate{}
	cr := csv.NewReader(r)
	cr.ReuseRecord = true

	// The CSV file starts with a header, which names the Dst currency in each
	// column. The rest of the lines in the file have the date in column 0, with
	// the remaining columns containing forex rates. The Src currency is always
	// EUR.

	record, err := cr.Read()
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	if err != nil {
		return nil, err
	}
	// csv.Reader will reuse `record`, so we have to make a copy.
	header := make([]string, len(record))
	copy(header, record)

	for {
		record, err := cr.Read()
		if err == io.EOF {
			return result, nil
		}

		if err != nil {
			return nil, err
		}

		t, err := time.Parse("2006-01-02", record[0])
		if err != nil {
			return nil, parseError(err, 0, cr)
		}
		t = t.UTC().Truncate(24 * time.Hour)

	ColumnLoop:
		for field := 1; field < len(record); field++ {
			value := record[field]
			if value == "" || value == "N/A" {
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
				Src:  "EUR",
				Dst:  currency,
				Day:  t,
				Rate: rate,
				Info: "ECB",
			})
		}
	}
}

func decompress(p []byte) (io.ReadCloser, error) {
	r, err := zip.NewReader(bytes.NewReader(p), int64(len(p)))
	if err != nil {
		return nil, err
	}
	var f *zip.File
	for _, f = range r.File {
		if filepath.Ext(f.Name) == "csv" {
			break
		}
	}

	if f == nil {
		return nil, errors.New("no csv file found")
	}

	return f.Open()
}
