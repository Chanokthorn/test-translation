package traverser

import "context"

type Traverser interface {
	Traverse(ctx context.Context, data any) error
}

type traverser struct{}

func (t *traverser) Traverse(ctx context.Context, data any) error {
	return nil
}
