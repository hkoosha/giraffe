package named

type Named struct {
	Value any
	Name  string
}

func Of(name string, v any) Named {
	return Named{
		Name:  name,
		Value: v,
	}
}
