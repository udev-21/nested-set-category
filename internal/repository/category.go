package repository

import (
	"context"

	"github.com/udev-21/nested-set-go/pkg/models"
)

type Category interface {
	FetchByID(ctx context.Context, id int64) (models.Category, error)
	FetchAllChildren(ctx context.Context, parent models.Category) ([]models.Category, error)
	FetchAllChildrenWithDepth(ctx context.Context, parent models.Category, depth int64) ([]models.Category, error)

	FetchLeafs(ctx context.Context, node models.Category) ([]models.Category, error)

	AppendRoot(ctx context.Context, node models.Category) error
	AppendChildren(ctx context.Context, parentNode models.Category, children models.Category) error

	AppendAfter(ctx context.Context, afterNode models.Category, children models.Category) error
	AppendBefore(ctx context.Context, beforeNode models.Category, children models.Category) error

	MoveAfter(ctx context.Context, afterNode, target models.Category) error
	MoveBefore(ctx context.Context, beforeNode, target models.Category) error

	Delete(ctx context.Context, node models.Category) error
	DeleteWithChildren(ctx context.Context, node models.Category) error

	RecalculateDepth(ctx context.Context) error
	GetAllRootIds(ctx context.Context) ([]string, error)
}
