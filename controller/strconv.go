package controller

import (
	"strconv"
	"strings"
	"time"
)

func IsEmpty(s string) bool {
	if s != "" {
		return strings.TrimLeft(s, " ") == ""
	}
	return true
}

func Trim(s string) string {
	s1 := strings.TrimLeft(s, " ")
	return strings.TrimRight(s1, " ")
}

func ParseState(s string) (names []string, values []string) {
	if s == "" {
		return
	}
	fields := strings.Split(Trim(s), ",")
	for _, fld := range fields {
		pair := strings.Split(Trim(fld), ":")
		names = append(names, Trim(pair[0]))
		if len(pair) == 1 {
			values = append(values, "")
		} else {
			values = append(values, Trim(pair[1]))
		}
	}
	return
}

func ExtractState(state, name string) string {
	if state == "" {
		return ""
	}
	fields := strings.Split(Trim(state), ",")
	for _, fld := range fields {
		pair := strings.Split(Trim(fld), ":")
		if Trim(pair[0]) == name {
			if len(pair) > 1 {
				return Trim(pair[1])
			} else {
				return ""
			}

		}
	}
	return ""
}

func ConvertDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, nil
	}
	tokens := strings.Split(s, "ms")
	if len(tokens) == 2 {
		val, err := strconv.Atoi(tokens[0])
		if err != nil {
			return 0, err
		}
		return time.Duration(val) * time.Millisecond, nil
	}
	tokens = strings.Split(s, "µs")
	if len(tokens) == 2 {
		val, err := strconv.Atoi(tokens[0])
		if err != nil {
			return 0, err
		}
		return time.Duration(val) * time.Microsecond, nil
	}
	tokens = strings.Split(s, "m")
	if len(tokens) == 2 {
		val, err := strconv.Atoi(tokens[0])
		if err != nil {
			return 0, err
		}
		return time.Duration(val) * time.Minute, nil
	}
	// Assume seconds
	tokens = strings.Split(s, "s")
	if len(tokens) == 2 {
		s = tokens[0]
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return time.Duration(val) * time.Second, nil
}
