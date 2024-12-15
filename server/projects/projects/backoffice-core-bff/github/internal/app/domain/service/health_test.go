package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewHealthService(t *testing.T) {
	s := NewHealthService()
	health, e := s.GetMessage()

	assert.Nil(t, e)
	assert.Equal(t, "it works", health.Message)
}

