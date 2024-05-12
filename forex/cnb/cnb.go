package cnb

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/wowsignal-io/go-forex/forex/exchange"
	"github.com/wowsignal-io/go-forex/forex/internal"
)

func SourceURLForDate(date time.Time) string {
	switch date.Weekday() {
	case time.Saturday:
		date = date.AddDate(0, 0, -1)
	case time.Sunday:
		date = date.AddDate(0, 0, -2)
	}
	return fmt.Sprintf("https://www.cnb.cz/cs/financni-trhy/devizovy-trh/kurzy-devizoveho-trhu/kurzy-devizoveho-trhu/denni_kurz.txt?date=%s", date.Format("02.01.2006"))
}

func Get(uri string) ([]exchange.Rate, error) {
	raw, err := internal.Fetch(uri)
	if err != nil {
		return nil, err
	}
	return parse(raw)
}

func parseCzechDecimal(s string) (float64, error) {
	s = strings.ReplaceAll(s, ",", ".")
	return strconv.ParseFloat(s, 64)
}

func parse(raw []byte) ([]exchange.Rate, error) {
	t, err := parseDate(raw)
	if err != nil {
		return nil, err
	}

	rates := []exchange.Rate{}
	raw, err = internal.SkipLinesBytes(raw, 2)
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	if err != nil {
		return nil, err
	}
	cr := csv.NewReader(bytes.NewReader(raw))
	cr.Comma = '|'
	cr.LazyQuotes = true
	cr.ReuseRecord = true
	for {
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(record) < 5 {
			continue
		}

		amount, err := parseCzechDecimal(record[2])
		if err != nil {
			return nil, fmt.Errorf("parse amount: %w", err)
		}
		rate, err := parseCzechDecimal(record[4])
		if err != nil {
			return nil, fmt.Errorf("parse rate: %w", err)
		}

		rates = append(rates, exchange.Rate{
			To:   "CZK",
			From: record[3],
			Day:  t,
			Rate: rate / amount,
		})
	}

	return rates, nil
}

func parseDate(raw []byte) (time.Time, error) {
	format := "02.01.2006"
	raw = raw[:len(format)]
	loc, err := time.LoadLocation("Europe/Prague")
	if err != nil {
		panic(err)
	}
	return time.ParseInLocation("02.01.2006", string(raw), loc)
}
