package port

import "context"

type AccessCheck interface {
	Check(ctx context.Context, action string, ids ...string) (allowedResources []string, err error)
	GetEmail(ctx context.Context) string
}
