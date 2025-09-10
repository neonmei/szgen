package generator

import (
	"fmt"
	"strconv"
	"strings"
)

func parseValue[T int64 | float64](v string) (T, error) {
	value := strings.TrimSpace(v)
	if value == "" {
		return T(0), fmt.Errorf("empty value")
	}

	switch any(T(0)).(type) {
	case int64:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return T(0), fmt.Errorf("invalid int64 value: %s", value)
		}
		return T(v), nil
	case float64:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return T(0), fmt.Errorf("invalid float64 value: %s", value)
		}
		return T(v), nil
	default:
		return T(0), fmt.Errorf("unsupported type")
	}
}

func parseRange[T int64 | float64](v string) ([]T, error) {
	value := strings.TrimSpace(v)
	if value == "" {
		return nil, fmt.Errorf("empty range")
	}

	parts := strings.Split(value, ",")
	result := make([]T, 0, len(parts))

	for _, v := range parts {
		parsedValue, err := parseValue[T](v)
		if err != nil {
			return nil, err
		}

		result = append(result, parsedValue)
	}

	return result, nil
}
