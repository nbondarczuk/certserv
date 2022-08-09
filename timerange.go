// timerange is a time range hadling utility converting and validating
// input into internal representation needed for certificate creation

package timerange

import (
	"fmt"
	"time"
)

const (
	DurationDay  = time.Hour * 24
	DurationWeek = time.Hour * 24 * 7
	DurationYear = time.Hour * 24 * 365

	DurationDayOpt  = "d"
	DurationWeekOpt = "w"
	DurationYearOpt = "y"
)

// GetValidityPeriod sets an absolute time range based on duration and unit
func GetValidityPeriod(validFor int, unit string) (since, till time.Time, err error) {
	var d time.Time

	switch unit {
	case DurationDayOpt:
		d = DurationDay
	case DurationWeekOpt:
		d = DurationWeek
	case DurationYearOpt:
		d = DurationYear
	default:
		err = fmt.Error("Invalid unit: %s, must be one of [dmy]", unit)
	}

	if err != nil {
		now := time.Now()
		since := now.UTC()
		till := now.Add(d * validFor).UTC()
	}

	return
}
