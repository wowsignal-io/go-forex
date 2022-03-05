package exchange

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestConvert(t *testing.T) {
	for _, tc := range []struct {
		comment    string
		data       []Rate
		want       Result
		wantErr    error
		from, to   string
		day        time.Time
		resultType ResultType
	}{
		{
			comment: "empty",
			wantErr: ErrNotFound,
		},
		{
			comment: "direct",
			data: []Rate{
				{
					From: "USD",
					To:   "EUR",
					Day:  time.Date(2022, time.January, 2, 0, 0, 0, 0, time.UTC),
					Rate: 0.9,
				},
			},
			from: "USD",
			to:   "EUR",
			day:  time.Date(2022, time.January, 2, 0, 0, 0, 0, time.UTC),
			want: Result{Rate: 0.9},
		},
		{
			comment: "inverse",
			data: []Rate{
				{
					From: "USD",
					To:   "EUR",
					Day:  time.Date(2022, time.January, 2, 0, 0, 0, 0, time.UTC),
					Rate: 0.9,
				},
			},
			from: "EUR",
			to:   "USD",
			day:  time.Date(2022, time.January, 2, 0, 0, 0, 0, time.UTC),
			want: Result{Rate: 1 / 0.9},
		},
		{
			comment: "wrong day (early)",
			data: []Rate{
				{
					From: "USD",
					To:   "EUR",
					Day:  time.Date(2022, time.January, 2, 0, 0, 0, 0, time.UTC),
					Rate: 0.9,
				},
			},
			from:    "USD",
			to:      "EUR",
			day:     time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
			wantErr: ErrNotFound,
		},
		{
			comment: "wrong day (late)",
			data: []Rate{
				{
					From: "USD",
					To:   "EUR",
					Day:  time.Date(2022, time.January, 2, 0, 0, 0, 0, time.UTC),
					Rate: 0.9,
				},
			},
			from:    "USD",
			to:      "EUR",
			day:     time.Date(2022, time.January, 3, 0, 0, 0, 0, time.UTC),
			wantErr: ErrNotFound,
		},
		{
			comment: "shortest path",
			data: []Rate{
				{
					From: "EUR",
					To:   "USD",
					Day:  time.Date(2022, time.January, 2, 0, 0, 0, 0, time.UTC),
					Rate: 1.2,
				},
				{
					From: "EUR",
					To:   "CZK",
					Day:  time.Date(2022, time.January, 2, 0, 0, 0, 0, time.UTC),
					Rate: 25,
				},
				{
					From: "EUR",
					To:   "CHF",
					Day:  time.Date(2022, time.January, 2, 0, 0, 0, 0, time.UTC),
					Rate: 1.1,
				},
				{
					From: "CZK",
					To:   "CHF",
					Day:  time.Date(2022, time.January, 2, 0, 0, 0, 0, time.UTC),
					Rate: 23,
				},
			},
			from:       "USD",
			to:         "CHF",
			day:        time.Date(2022, time.January, 2, 0, 0, 0, 0, time.UTC),
			resultType: FullTrace,
			// The longer path is USD -> EUR -> CZK -> CHF yielding (1/1.2) * 25 * (1/23)
			want: Result{
				Rate: (1 / 1.2) * 1.1,
				Trace: []Rate{
					{From: "USD", To: "EUR", Rate: 1 / 1.2},
					{From: "EUR", To: "CHF", Rate: 1.1},
				},
			},
		},
	} {
		t.Run(tc.comment, func(t *testing.T) {
			g, err := Compile(tc.data)
			if err != nil {
				t.Fatalf("Compile(%#v) -> nil, %v", tc.data, err)
			}
			t.Logf("Compile(%#v)", tc.data)

			result, err := Convert(g, tc.from, tc.to, tc.day, tc.resultType)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Convert(%#v, %q, %q, %v, %v) -> err=%v wanted (err=%v)", g, tc.from, tc.to, tc.day, tc.resultType, err, tc.wantErr)
			}

			if diff := cmp.Diff(tc.want, result, cmpopts.EquateApprox(0, 0.0001), cmpopts.IgnoreFields(Rate{}, "Day", "Info")); diff != "" {
				t.Errorf("Convert(%#v, %q, %q, %v, %v) -> (-) wanted vs. (+) got:\n%s ", g, tc.from, tc.to, tc.day, tc.resultType, diff)
			}
		})
	}
}
