package repository

import (
	"context"

	"github.com/udev-21/nested-set-go/models"
)

type Category interface {
	FetchByID(ctx context.Context, id int64) (models.Category, error)                                              //done
	FetchAllChildren(ctx context.Context, parent models.Category) ([]models.Category, error)                       //done
	FetchAllChildrenWithDepth(ctx context.Context, parent models.Category, depth int64) ([]models.Category, error) //done
	FetchLeafs(ctx context.Context, node models.Category) ([]models.Category, error)                               //done
	AppendRoot(ctx context.Context, node models.Category) error                                                    //done
	AppendChildren(ctx context.Context, parentNode models.Category, children []models.Category) error              //done
	AppendAfter(ctx context.Context, afterNode models.Category, children []models.Category) error                  //done
	AppendBefore(ctx context.Context, beforeNode models.Category, children []models.Category) error                //done
	Delete(ctx context.Context, node models.Category) error                                                        //done
	DeleteWithChildren(ctx context.Context, node models.Category) error                                            //done
	RecalculateDepth(ctx context.Context) error                                                                    //done
	GetAllRootIds(ctx context.Context) ([]string, error)                                                           //done
	MoveAfter(ctx context.Context, afterNode, target models.Category) error                                        //done
	MoveBefore(ctx context.Context, beforeNode, target models.Category) error                                      //done
	MoveInto(ctx context.Context, parentNode, target models.Category) error                                        //done
}
