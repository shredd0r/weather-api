package util

import (
	"fmt"
)

func Float64ToString(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func GetIconResourceNameByAccuWeatherIcon(iconId *uint8) *string {
	return nil
}

func GetIconResourceNameByWeatherIcon(weatherIcon uint8) *string {
	return GetIconResourceNameByAccuWeatherIcon(&weatherIcon)
}
