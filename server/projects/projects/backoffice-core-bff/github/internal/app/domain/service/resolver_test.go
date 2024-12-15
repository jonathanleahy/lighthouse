package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func TestSetResolver(t *testing.T) {
	resolver := NewResolver(nil)

	assert.NotNil(t, resolver)
}

func TestCustomFormatter(t *testing.T) {
	t.Run("error is not a gqlerror", func(t *testing.T) {
		err := CustomFormatter(context.Background(), context.Canceled)
		assert.NotNil(t, err)
		assert.Equal(t, "internal server error", err.Message)
	})

	t.Run("error is gqlerror without extensions", func(t *testing.T) {
		gqlerror := &gqlerror.Error{}
		err := CustomFormatter(context.Background(), gqlerror)
		assert.NotNil(t, err)
	})

	t.Run("error is gqlerror with extensions and with code", func(t *testing.T) {
		gqlerror := &gqlerror.Error{
			Extensions: map[string]interface{}{
				"code": "GRAPHQL_VALIDATION_FAILED",
			},
		}
		err := CustomFormatter(context.Background(), gqlerror)
		assert.NotNil(t, err)
		assert.Equal(t, "invalid request data", err.Message)
	})

}

