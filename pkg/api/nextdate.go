package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const formatDate = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat cannot be empty")
	}

	startTime, err := time.Parse(formatDate, dstart)
	if err != nil {
		return "", fmt.Errorf("incorrect start date: %w", err)
	}

	if repeat == "y" {
		for {
			startTime = startTime.AddDate(1, 0, 0)
			if startTime.After(now) {
				return startTime.Format(formatDate), nil
			}
		}
	} else if after, ok := strings.CutPrefix(repeat, "d "); ok {
		daysString := after
		days, err := strconv.Atoi(daysString)
		if err != nil {
			return "", fmt.Errorf("wrong format: %w", err)
		}
		if days <= 0 || days > 400 {
			return "", fmt.Errorf("daily interval should be from 1 to 400")
		}

		for {
			startTime = startTime.AddDate(0, 0, days)
			if startTime.After(now) {
				return startTime.Format(formatDate), nil
			}

		}
	} else if repeat == "w" || repeat == "m" {
		return "", fmt.Errorf("wrong format")
	} else {
		return "", fmt.Errorf("unknown format")
	}
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowString := r.URL.Query().Get("now")
	dateString := r.URL.Query().Get("date")
	repeatString := r.URL.Query().Get("repeat")

	var now time.Time
	var err error

	if nowString == "" {
		now = time.Now()
	} else {
		now, err = time.Parse("20060202", nowString)
		if err != nil {
			http.Error(w, fmt.Sprintf("wrong format: %v", err), http.StatusBadRequest)
			return
		}
	}

	result, err := NextDate(now, dateString, repeatString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, result)
}
