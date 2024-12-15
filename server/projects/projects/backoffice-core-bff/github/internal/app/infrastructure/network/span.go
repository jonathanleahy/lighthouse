package network

import (
	"errors"
	"go.opentelemetry.io/otel/attribute"
)

type Span struct {
	Name   string
	Values []attribute.KeyValue
	Ignore bool
}

func (s *Span) validate() error {
	if !s.Ignore && len(s.Name) == 0 {
		return errors.New(missingSpanName)
	}

	return nil
}

