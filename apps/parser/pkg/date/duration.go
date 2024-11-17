package date

import (
	"fmt"
	"strings"
	"time"
)

func dateDiff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = y2 - y1
	month = int(M2 - M1)
	day = d2 - d1
	hour = h2 - h1
	min = m2 - m1
	sec = s2 - s1

	// Normalize negative values.
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// Days in month.
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

type DurationOptsHide struct {
	Years   bool
	Months  bool
	Days    bool
	Hours   bool
	Minutes bool
	Seconds bool
}

type DurationOpts struct {
	UseUtc bool
	Hide   DurationOptsHide
}

func Duration(from time.Time, opts *DurationOpts) string {
	if opts == nil {
		opts = &DurationOpts{}
	}

	var date strings.Builder

	var y, m, d, h, mi, s int

	if opts.UseUtc == true {
		y, m, d, h, mi, s = dateDiff(from, time.Now().UTC())
	} else {
		y, m, d, h, mi, s = dateDiff(from, time.Now())
	}

	if y > 0 && !opts.Hide.Years {
		fmt.Fprintf(&date, "%dy ", y)
	}

	if m > 0 && !opts.Hide.Months {
		fmt.Fprintf(&date, "%dmo ", m)
	}

	if d > 0 && !opts.Hide.Days {
		fmt.Fprintf(&date, "%dd ", d)
	}

	if h > 0 && !opts.Hide.Hours {
		fmt.Fprintf(&date, "%dh ", h)
	}

	if mi > 0 && !opts.Hide.Minutes {
		fmt.Fprintf(&date, "%dm ", mi)
	}

	if s > 0 && !opts.Hide.Seconds {
		fmt.Fprintf(&date, "%ds", s)
	}

	return date.String()
}
