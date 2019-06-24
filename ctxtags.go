package ctxtags

import "context"

type keyT int
const key keyT = 0

func WithTags(ctx context.Context, m map[string]string) context.Context {
	stack, _ := ctx.Value(key).(*mapStack)
	stack = stack.push(m)
	return context.WithValue(ctx, key, stack)
}

func Tags(ctx context.Context) map[string]string {
	if ctx == nil {
		return map[string]string{}
	}

	stack, _ := ctx.Value(key).(*mapStack)
	return stack.flattened()
}

type mapStack struct {
	m    map[string]string
	next *mapStack
}

func (ms *mapStack) push(m map[string]string) *mapStack {
	return &mapStack{m: m, next: ms}
}

func (ms *mapStack) flattened() map[string]string {
	flat := map[string]string{}

	for ms != nil {
		for k, v := range ms.m {
			if _, found := flat[k]; !found {
				flat[k] = v
			}
		}

		ms = ms.next
	}

	return flat
}
