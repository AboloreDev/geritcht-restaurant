package utils

import (
	"fmt"
	"time"

	"gorm.io/datatypes"
)

var timeSlots = []string{
	"18:00:00",
	"20:00:00",
	"21:00:00",
	"22:00:00",
	"23:00:00",
	"02:00:00",
}

func IsValidTimeSlots(slot string) bool {
	for _, validSlot := range timeSlots {
		if validSlot == slot {
			return true
		}
	}
	return false
}

func ParseToDataTypesTime(slot string) (datatypes.Time, error) {
	t, err := time.Parse("15:04:00", slot)
	if err != nil {
		return 0, err
	}
	nanos :=
		int64(t.Hour())*3600*1e9 +
			int64(t.Minute())*60*1e9 +
			int64(t.Second())*1e9

	return datatypes.Time(nanos), nil
}

func FormatDataTypesTime(t datatypes.Time) string {
	totalSeconds := int64(t) / 1e9

	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

// TimeToDataTypesTime converts a time.Time directly into datatypes.Time,
// using only its wall-clock hour/minute/second. No string parsing involved,
// so it can't fall victim to layout-mismatch bugs.
func TimeToDataTypesTime(t time.Time) datatypes.Time {
	nanos := int64(t.Hour())*3600*1e9 +
		int64(t.Minute())*60*1e9 +
		int64(t.Second())*1e9
	return datatypes.Time(nanos)
}
