package model

import (
	"strconv"
)

type Record = map[string]interface{}

func ParseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

func CleanValue(s string) interface{} {
	if s == "--" || s == "-" || s == "" {
		return nil
	}
	// Try integer
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}
	// Try float
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return s
}

func CleanFloat(s string) interface{} {
	if s == "--" || s == "-" || s == "" {
		return nil
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return s
}
