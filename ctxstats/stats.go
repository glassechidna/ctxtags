package ctxstats

import "context"

type Stats interface {
	Count(ctx context.Context, name string, value int64, tagMap map[string]string, rate float64) error
	Gauge(ctx context.Context, name string, value float64, tagMap map[string]string, rate float64) error
	Histogram(ctx context.Context, name string, value float64, tagMap map[string]string, rate float64) error
}
