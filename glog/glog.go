package glog

type GLog interface {
	Named(string) GLog

	Debug(msg string, fields ...any)

	Info(msg string, fields ...any)

	Warn(msg string, fields ...any)

	Error(msg string, fields ...any)
}
