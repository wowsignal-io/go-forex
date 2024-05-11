package offline

import (
	_ "embed"
)

//go:embed eurofxref-hist.zip
var HistoricalECBRates string

//go:embed boc_offline_rates.csv
var HistoricalBOCRates string
