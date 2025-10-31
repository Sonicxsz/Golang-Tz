package dto

import (
	"fmt"
	"time"
)

// ParseMonthYear парсит строку "MM-YYYY" в time.Time
func ParseMonthYear(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, fmt.Errorf("date string is empty")
	}
	return time.Parse("01-2006", dateStr)
}
func FormatMonthYear(t time.Time) string {
	return t.Format("01-2006")
}
