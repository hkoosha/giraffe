package internal

import "time"

func Marshal(
	format string,
	t time.Time,
) ([]byte, error) {
	f := `"` + t.Format(format) + `"`

	return []byte(f), nil
}

func Unmarshal(
	format string,
	b []byte,
) (time.Time, error) {
	str := string(b)

	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	return time.Parse(format, str)
}
