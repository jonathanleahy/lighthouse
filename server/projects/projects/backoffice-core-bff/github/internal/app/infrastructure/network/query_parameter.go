package network

type QueryParameter struct {
	Name      string
	Value     interface{}
	Sensitive bool
}

func (q *QueryParameter) gdprValue() interface{} {
	if q.Sensitive {
		return sensitivePlaceholderValue
	}

	return q.Value
}

