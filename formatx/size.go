package formatx

import (
	"fmt"
	"math"
)

var (
	defaultUnits = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
)

func FormatSizePrecise(bytes int64) string {
	return FormatSizePreciseWithUnits(bytes, 1024, nil) //nolint:mnd // .
}

func FormatSizePreciseWithUnits(bytes, unit int64, units []string) string {
	if len(units) == 0 {
		units = defaultUnits
	}

	if bytes < unit {
		return fmt.Sprintf("%d %s", bytes, units[0])
	}

	exp := int(math.Log(float64(bytes)) / math.Log(float64(unit)))

	if exp >= len(units) {
		exp = len(units) - 1
	}

	value := float64(bytes) / math.Pow(float64(unit), float64(exp))

	precision := 2

	if value >= 100 { //nolint:mnd // .
		precision = 0
	} else if value >= 10 { //nolint:mnd // .
		precision = 1
	}

	return fmt.Sprintf("%.*f %s", precision, value, units[exp])
}
