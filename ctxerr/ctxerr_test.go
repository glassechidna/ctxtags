package ctxerr

import (
	"context"
	"github.com/glassechidna/ctxtags"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWrap(t *testing.T) {
	t.Run("empty ctx", func(t *testing.T) {
		err := errors.New("base error")
		wrapped := Wrap(context.Background(), err)
		assert.Equal(t, "ctx: {}: base error", wrapped.Error())
	})

	t.Run("single tag", func(t *testing.T) {
		ctx := context.Background()
		ctx = ctxtags.WithTags(ctx, map[string]string{"key": "val"})
		err := errors.New("base error")
		wrapped := Wrap(ctx, err)
		assert.Equal(t, `ctx: {"key":"val"}: base error`, wrapped.Error())
	})

	t.Run("returns nil", func(t *testing.T) {
		wrapped := Wrap(context.Background(), nil)
		assert.Nil(t, wrapped)
	})

	t.Run("nil ctx is fine", func(t *testing.T) {
		err := errors.New("base error")
		wrapped := Wrap(nil, err)
		assert.Equal(t, "ctx: {}: base error", wrapped.Error())
	})
}

