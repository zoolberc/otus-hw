package utils

import (
	"time"

	"github.com/google/uuid"
	"github.com/zoolberc/otus-hw/hw12_13_14_15_calendar/internal/constants"
)

func ParseDate(date string) (time.Time, error) {
	d, err := time.Parse(constants.DateFormatLayout, date)
	if err != nil {
		return time.Time{}, constants.ErrIncorrectDateFormat
	}
	return d, nil
}

func GenerateUUID() string {
	return uuid.NewString()
}
