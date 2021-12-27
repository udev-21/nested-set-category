package repos

import "context"

type Tree interface {
	GetByID(ctx context.Context, userFilter int64) error
	GetChildrenWithDepth(ctx context.Context, id int64, depth int) error

	CreateRoot(ctx context.Context, name string) error
	AppendInto(ctx context.Context) error
	AppendAfter(ctx context.Context, name string) error
	AppendBefore(ctx context.Context, name string) error

	MoveAfter(ctx context.Context, name string) error
	MoveBefore(ctx context.Context, name string) error

	Delete(ctx context.Context, userId uint64) error
	DeleteWithChildren(ctx context.Context, userId uint64) error

	RecalculateDepth(ctx context.Context) error
}
