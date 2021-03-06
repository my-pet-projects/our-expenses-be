package domain

import (
	"time"

	"github.com/pkg/errors"
)

// DateRange holds data range.
type DateRange struct {
	from time.Time
	to   time.Time
}

// NewDateRange returns date range struct.
func NewDateRange(from time.Time, to time.Time) (*DateRange, error) {
	if from.After(to) {
		return nil, errors.New("From date is after To date")
	}
	dr := &DateRange{
		from: from,
		to:   to,
	}

	return dr, nil
}

// From returns from date.
func (dr DateRange) From() time.Time {
	return dr.from
}

// To returns to date.
func (dr DateRange) To() time.Time {
	return dr.to
}

// DatesInBetween returns all dates between from and to dates.
func (dr DateRange) DatesInBetween() []time.Time {
	var dates []time.Time
	from := dr.from
	for {
		dates = append(dates, from)
		fromYear, fromMonth, fromDay := from.Date()
		toYear, toMonth, toDay := dr.to.Date()
		if fromYear == toYear && fromMonth == toMonth && fromDay == toDay {
			break
		}
		from = from.Add(1 * 24 * time.Hour)
	}

	return dates
}
