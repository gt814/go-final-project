package tasks

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, strDate string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat rule not specified")
	}

	date, err := time.Parse("20060102", strDate)
	if err != nil {
		return "", errors.New("invalid date format")
	}

	repeatFields := strings.Fields(repeat)
	if len(repeatFields) < 1 {
		return "", errors.New("invalid repeat format")
	}

	rule := repeatFields[0]
	value := ""
	if len(repeatFields) > 1 {
		value = repeatFields[1]
	}

	var nextDate time.Time
	switch rule {
	case "d":
		if value == "" {
			return "", errors.New("invalid repeat format")
		}
		days, err := strconv.Atoi(value)
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("invalid repeat format")
		}
		nextDate = date.AddDate(0, 0, days)
		for now.After(nextDate) {
			nextDate = nextDate.AddDate(0, 0, days)
		}

	case "y":
		nextDate = date.AddDate(1, 0, 0)
		for now.After(nextDate) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}

	default:
		return "", errors.New("invalid repeat rule")
	}

	return nextDate.Format("20060102"), nil
}
