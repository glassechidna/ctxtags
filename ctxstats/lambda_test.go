package ctxstats

import (
	"context"
	"github.com/glassechidna/ctxtags"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

type mockNower struct{}

func (m *mockNower) Now() time.Time {
	return time.Unix(100, 0)
}

func TestLambda_Count(t *testing.T) {
	t.Run("nil tags", func(t *testing.T) {
		w := &strings.Builder{}
		l := &Lambda{Writer: w, nower: &mockNower{}}
		err := l.Count(context.Background(), "users", 7, nil, 1)
		assert.NoError(t, err)
		assert.Equal(t, "MONITORING|100|7|count|users|\n", w.String())
	})

	t.Run("empty tags", func(t *testing.T) {
		w := &strings.Builder{}
		l := &Lambda{Writer: w, nower: &mockNower{}}
		err := l.Count(context.Background(), "users", 7, map[string]string{}, 1)
		assert.NoError(t, err)
		assert.Equal(t, "MONITORING|100|7|count|users|\n", w.String())
	})
}

func TestLambda_Tags(t *testing.T) {
	t.Run("one tag", func(t *testing.T) {
		w := &strings.Builder{}
		l := &Lambda{Writer: w, nower: &mockNower{}}
		err := l.Count(context.Background(), "users", 7, map[string]string{"key": "value"}, 1)
		assert.NoError(t, err)
		assert.Equal(t, "MONITORING|100|7|count|users|#key:value\n", w.String())
	})

	t.Run("two tags", func(t *testing.T) {
		w := &strings.Builder{}
		l := &Lambda{Writer: w, nower: &mockNower{}}
		err := l.Count(context.Background(), "users", 7, map[string]string{"key": "value", "second": "tag"}, 1)
		assert.NoError(t, err)
		assert.Equal(t, "MONITORING|100|7|count|users|#key:value,second:tag\n", w.String())
	})

	t.Run("non-overlapping context", func(t *testing.T) {
		ctx := context.Background()
		ctx = ctxtags.WithTags(ctx, map[string]string{"ctxkey": "ctxval"})

		w := &strings.Builder{}
		l := &Lambda{Writer: w, nower: &mockNower{}}
		err := l.Count(ctx, "users", 7, map[string]string{"key": "value"}, 1)
		assert.NoError(t, err)
		assert.Equal(t, "MONITORING|100|7|count|users|#ctxkey:ctxval,key:value\n", w.String())
	})

	t.Run("overlapping context", func(t *testing.T) {
		ctx := context.Background()
		ctx = ctxtags.WithTags(ctx, map[string]string{"ctxkey": "ctxval"})

		w := &strings.Builder{}
		l := &Lambda{Writer: w, nower: &mockNower{}}
		err := l.Count(ctx, "users", 7, map[string]string{"ctxkey": "topval"}, 1)
		assert.NoError(t, err)
		assert.Equal(t, "MONITORING|100|7|count|users|#ctxkey:topval\n", w.String())
	})
}

func TestLambda_Gauge(t *testing.T) {
	w := &strings.Builder{}
	l := &Lambda{Writer: w, nower: &mockNower{}}
	err := l.Gauge(context.Background(), "memory", 7, nil, 1)
	assert.NoError(t, err)
	assert.Equal(t, "MONITORING|100|7.000000|gauge|memory|\n", w.String())
}

func TestLambda_Histogram(t *testing.T) {
	w := &strings.Builder{}
	l := &Lambda{Writer: w, nower: &mockNower{}}
	err := l.Histogram(context.Background(), "memory", 7, nil, 1)
	assert.NoError(t, err)
	assert.Equal(t, "MONITORING|100|7.000000|histogram|memory|\n", w.String())
}
