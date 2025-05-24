package glog

type Level uint

const (
	Disabled Level = iota
	Error
	Warn
	Info
	Debug
)

func (l Level) String() string {
	switch l {
	case Disabled:
		return "NOP"

	case Error:
		return "ERR"

	case Warn:
		return "WRN"

	case Info:
		return "INF"

	case Debug:
		return "DBG"

	default:
		return "UNK"
	}
}
