package gzapadapter

type Hook = func(msg string, err error, fields ...any)
