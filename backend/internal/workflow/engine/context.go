package engine

import "context"

type mutationContextKey struct{}

// MutationMeta carries reentrancy metadata on context for record hooks.
type MutationMeta struct {
	Source            string
	Depth             int
	ExcludeWorkflowID string
}

func WithMutationMeta(ctx context.Context, meta MutationMeta) context.Context {
	return context.WithValue(ctx, mutationContextKey{}, meta)
}

func MutationMetaFrom(ctx context.Context) (MutationMeta, bool) {
	v, ok := ctx.Value(mutationContextKey{}).(MutationMeta)
	return v, ok
}
