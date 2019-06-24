package ctxtags

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTags(t *testing.T) {
	ctx := context.Background()
	ctx = WithTags(ctx, map[string]string{
		"top1": "topval1",
		"top2": "topval2",
	})
	ctx = WithTags(ctx, map[string]string{
		"bottom1": "bottomval1",
		"top2":    "bottom2",
	})

	flat := Tags(ctx)
	assert.Equal(t, flat, map[string]string{
		"top1":    "topval1",
		"top2":    "bottom2",
		"bottom1": "bottomval1",
	})
}