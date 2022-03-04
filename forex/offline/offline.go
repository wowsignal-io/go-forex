package offline

import (
	_ "embed"
)

//go:embed eurofxref-hist.zip
var HistoricalECBRates string
