package network

type PathParameter struct {
	Value     interface{}
	Sensitive bool
}

func (p *PathParameter) gdprValue() interface{} {
	if p.Sensitive {
		return sensitivePlaceholderValue
	}

	return p.Value
}

