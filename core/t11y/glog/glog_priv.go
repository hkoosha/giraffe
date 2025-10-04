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
	//nolint:forbidigo
	fmt.Printf(
		`{"message": "%s", "details": "%s"}%s`,
		poorJSON(message),
		poorJSON(details),
		"\n",
	)
}

type poorManGLog struct{}

func (p poorManGLog) Log(level Level, msg string, fields ...any) {
	switch level {
	case Debug:
		p.Debug(msg, fields...)
	case Info:
		p.Info(msg, fields...)
	case Warn:
		p.Warn(msg, fields...)
	case Error:
		p.Error(msg, fields...)
	case Disabled:
		// Nothing.
	default:
		p.Error(msg, fields...)
	}
}

func (p poorManGLog) IsDebug() bool {
	return true
}

func (p poorManGLog) IsInfo() bool {
	return true
}

func (p poorManGLog) IsWarn() bool {
	return true
}

func (p poorManGLog) IsError() bool {
	return true
}

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

func (p poorManGLog) Err(msg string, err error, fields ...any) {
	poorLog(msg, append(fields, err))
}

var global Lg = &poorManGLog{}
