# go-forex

Simple and efficient Go library for getting daily foreign exchange rates.
Built-in support for ca. 80 currencies.

Also includes a simple [commandline tool](#commandline-interface).

Motivations, alternatives and trade-offs are discussed in the technical [design document](https://www.wowsignal.io/articles/go-forex).

## Examples

The following example will automatically download (and cache) exchange rates
from several central banks (see the list below) and then compute the exchange
rate.

```go
import "github.com/wowsignal-io/go-forex/forex"

rate, err := forex.LiveExchange().Convert("USD", "EUR", time.Date(2022, time.January, 4, 0, 0, 0, 0, time.UTC))
if err != nil { /* Handle errors. */ }
fmt.Printf("The conversion rate from USD to EUR on January 4, 2022 was %f\n", rate.Rate)
// Output: The conversion rate from USD to EUR on January 4, 2022 was 0.886603.
```

Currency combinations that are not directly published are supported by automatic
intermediate conversions. For example, we can convert Papuan Kina to the Indian
Rupee.

```go
rate, err := forex.LiveExchange().Convert("PKG", "INR", time.Date(2022, time.January, 4, 0, 0, 0, 0, time.UTC))
if err != nil { /* Handle errors. */ }
fmt.Printf("The conversion rate from PKG to INR on January 4, 2022 was %f\n", rate.Rate)
// Output: The conversion rate from PKG to INR on January 4, 2022 was 21.255237.
```

Enabling the full trace shows us how the rate was computed. Here, for example,
it was done by conversion through the Australian Dollar, using rates from the
Royal Bank of Australia (RBA).

```go
rate, err := forex.LiveExchange().Convert("PKG", "INR", time.Date(2022, time.January, 4, 0, 0, 0, 0, time.UTC), exchange.FullTrace)
if err != nil { /* Handle errors. */ }

// Passing exchange.FullTrace as the last argument means rate.Trace is now
// populated.
for i, step := range rate.Trace {
    fmt.Printf("Conversion step %d/%d: 1 %s = %f %s (source: %s)\n", i+1, len(rate.Trace), step.From, step.Rate, step.To, step.Info)
}

// Output:
// Conversion step 1/2: 1 PGK = 0.395226 AUD (source: RBA (inverse))
// Conversion step 2/2: 1 AUD = 53.780000 INR (source: RBA)
```

## Commandline interface

A command called `forex-convert` is provided exposing the above API over the
commandline.

To install:

```sh
go install github.com/wowsignal-io/go-forex/cmd/forex-convert@latest
# This places forex-convert in $GOBIN.
```

Convert some currencies:

```sh
forex-convert -from=PGK -to=INR -date=2021-03-01 -v
# Outputs:
# Conversion step 1/2: 1 PGK = 0.367985 AUD (source: RBA (inverse))
# Conversion step 2/2: 1 AUD = 56.850000 INR (source: RBA)
# Computed rate: 20.919963
```

Or, to get just the number:

```sh
forex-convert -from=PGK -to=INR -date=2021-03-01
# Outputs:
# 20.919963
```

Also supports other options, such as offline operation and search tolerances. Run `forex-convert --help`.

## Offline operation

All above examples use the `LiveExchange`, which downloads and caches exchange
data from the internet. For offline operation, `OfflineExchange` can be used as
a drop-in replacement. However, fewer currencies are supported and only
historical data is available. (However, it's easy to update the data used by the
`OfflineExchange` using a cronjob or similar.)

## Performance

The algorithm is breadth-first walk through the data while filtering edges
sorted by time. Runtime scales linearly with the number of currencies and
logarithmically with the length of historical data. On an M1 MacBook, computing
and indirect exchange rate takes about 4,000 ns and requires about 7,000 bytes
of storage.

## Supported currencies and sources

Data are sourced from the following banks:

* European Central Bank (ECB)
* Royal Bank of Australia (RBA)
* Bank of Canada (BOC)
* Central Bank of the U.A.E. (CBUAE)
* The Czech National Bank (CNB)

Data are refreshed every 12 hours (or manually) and cached locally in /tmp or
similar path. See [currencies.txt](forex/currencies.txt) for a full list of
supported currencies.

The computed exchange rates are for informational purposes only - they are
unlikely to be the same as the rates actually offered, but the difference should
be tolerable for home finance applications.
