package ctxstats

import (
	"context"
	"fmt"
	"github.com/glassechidna/ctxtags"
	"io"
	"os"
	"sort"
	"strings"
	"time"
)

type Lambda struct {
	Writer io.Writer
	nower nower
}

type nower interface {
	Now() time.Time
}

func (l *Lambda) Count(ctx context.Context, name string, value int64, tags map[string]string, rate float64) error {
	tags = ctxtags.Tags(ctxtags.WithTags(ctx, tags))
	return l.write(name, fmt.Sprintf("%d", value), "count", tags)
}

func (l *Lambda) Gauge(ctx context.Context, name string, value float64, tags map[string]string, rate float64) error {
	tags = ctxtags.Tags(ctxtags.WithTags(ctx, tags))
	return l.write(name, fmt.Sprintf("%f", value), "gauge", tags)
}

func (l *Lambda) Histogram(ctx context.Context, name string, value float64, tags map[string]string, rate float64) error {
	tags = ctxtags.Tags(ctxtags.WithTags(ctx, tags))
	return l.write(name, fmt.Sprintf("%f", value), "histogram", tags)
}

func (l *Lambda) write(name, value, typ string, tags map[string]string) error {
	now := time.Now().Unix()
	if l != nil && l.nower != nil {
		now = l.nower.Now().Unix()
	}

	line := fmt.Sprintf("MONITORING|%d|%s|%s|%s|%s\n", now, value, typ, name, statsdTags(tags))

	var w io.Writer = os.Stdout
	if l != nil && l.Writer != nil {
		w = l.Writer
	}

	_, err := io.WriteString(w, line)
	return err
}

func statsdTags(m map[string]string) string {
	var keys []string
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sb := &strings.Builder{}
	for _, k:= range keys {
		sb.WriteString(k)
		sb.WriteString(":")
		sb.WriteString(m[k])
		sb.WriteString(",")
	}

	res := sb.String()
	if len(res) > 0 {
		res = "#" + strings.TrimSuffix(res, ",")
	}

	return res
}
