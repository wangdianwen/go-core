package utils

import (
	"errors"
	"time"
)

func ParseDateTime(str string) (*time.Time, error) {
	var parsed time.Time
	parsed, err := time.Parse("2006-01-02 15:04:05", str)
	if err == nil {
		return &parsed, nil
	}
	parsed, err = time.Parse("2006-01-02", str)
	if err == nil {
		return &parsed, nil
	}
	return nil, errors.New("datetime format are '2006-01-02 15:04:05' or '2006-01-02")
}
