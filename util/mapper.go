package util

import (
	"fmt"
)

func Float64ToString(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func PercentToFloat64(percent int) float64 {
	return float64(percent) / 100.0
}

func PercentToFloat64Pointer(percent *int) *float64 {
	f := PercentToFloat64(*percent)
	return &f
}

func GetIconResourceNameByAccuWeatherIcon(iconId *uint8) *string {
	return nil
}

func GetIconResourceNameByWeatherIcon(weatherIcon uint8) *string {
	return GetIconResourceNameByAccuWeatherIcon(&weatherIcon)
}
