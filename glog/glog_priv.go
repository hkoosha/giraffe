package glog

import (
	"fmt"
	"strings"
)

func poorJSON(value any) string {
	if value == nil {
		return ""
	}

	var s string

	switch value := value.(type) {
	case string:
		s = value

	default:
		s = fmt.Sprintf("%v", value)
	}

	s = strings.ReplaceAll(s, `"`, "'")
	s = strings.ReplaceAll(s, "\n", " ")

	return s
}

func poorLog(
	message string,
	details any,
) {
	fmt.Printf(
		`{"message": "%s", "details": "%s"}%s`,
		poorJSON(message),
		poorJSON(details),
		"\n",
	)
}

type poorManGLog struct{}

func (p poorManGLog) Named(string) Lg {
	return p
}

func (p poorManGLog) Debug(msg string, fields ...any) {
	poorLog(msg, fields)
}

func (p poorManGLog) Info(msg string, fields ...any) {
	poorLog(msg, fields)
}

func (p poorManGLog) Warn(msg string, fields ...any) {
	poorLog(msg, fields)
}

func (p poorManGLog) Error(msg string, fields ...any) {
	poorLog(msg, fields)
}

func (p poorManGLog) Of(key string, value ...any) any {
	value = append(value, "")
	copy(value[1:], value[:len(value)-1])
	value[0] = key
	return value
}

var global Lg = &poorManGLog{}
