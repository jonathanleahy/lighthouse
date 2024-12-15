package handler

import (
	"github.com/pismo/backoffice-core-bff/internal/app/domain/service"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetManager(t *testing.T) {
	expectedResolver := service.NewResolver(nil)
	SetResolver(expectedResolver)
	assert.Equal(t, expectedResolver, resolver)
}

