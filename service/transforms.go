package service

import (
	"fmt"
	"strconv"
)

// formatIntegerWithMetricPrefix formats an integer into a string with metric prefix
func formatIntegerWithMetricPrefix(n int) string {
	var metricPrefix string
	var result float64
	switch {
	case n > 10e8:
		metricPrefix = "G"
		result = float64(n) / 10e8
	case n > 10e5:
		metricPrefix = "M"
		result = float64(n) / 10e5
	case n > 10e2:
		metricPrefix = "k"
		result = float64(n) / 10e2
	default:
		metricPrefix = ""
		result = float64(n)
	}

	precision := 0
	if metricPrefix != "" {
		precision = 3 - len(strconv.Itoa(int(result)))
	}

	formatSpecifier := fmt.Sprintf("%%.%df%s", precision, metricPrefix)
	return fmt.Sprintf(formatSpecifier, result)
}
