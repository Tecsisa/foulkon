package grpc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Tecsisa/foulkon/api"
)

// Note: This doesn't test the switch statement in statusFromError, only
//       interesting properties.
func TestStatusFromError(t *testing.T) {
	t.Run("returns error with underlying grpc status", func(t *testing.T) {
		err := statusFromError(&api.Error{Code: api.USER_ALREADY_EXIST})
		s, ok := status.FromError(err)
		if assert.True(t, ok) {
			assert.NotNil(t, s)
		}
	})

	t.Run("returns 'Internal' code for unknown API error codes", func(t *testing.T) {
		err := statusFromError(&api.Error{Code: "thisisunknowninvalid"})
		s, ok := status.FromError(err)
		if assert.True(t, ok) && assert.NotNil(t, s) {
			assert.Equal(t, s.Code(), codes.Internal)
		}
	})

	t.Run("returns error as-is if it's not an API error", func(t *testing.T) {
		oops := fmt.Errorf("oops")
		err := statusFromError(oops)
		s, ok := status.FromError(err)
		assert.False(t, ok)
		assert.Nil(t, s)
		assert.Equal(t, oops, err)
	})

	t.Run("returns status which includes original message", func(t *testing.T) {
		msg := "oh fish"
		err := statusFromError(&api.Error{Code: api.USER_ALREADY_EXIST, Message: msg})
		s, ok := status.FromError(err)
		if assert.True(t, ok) && assert.NotNil(t, s) {
			assert.Equal(t, s.Message(), msg)
		}
	})
}
