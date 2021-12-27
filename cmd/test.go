package main

import (
	"fmt"

	"github.com/udev-21/nested-set-go/pkg/models"
	"github.com/udev-21/nested-set-go/pkg/utils"
)

func main() {
	data := models.Tree{
		RawTree: models.RawTree{
			Name:  "root",
			Depth: 1,
		},
		Children: []models.Tree{
			{
				RawTree: models.RawTree{
					Name:  "child1",
					Depth: 2,
				},
				Children: []models.Tree{
					{
						RawTree: models.RawTree{
							Name:  "child1.1",
							Depth: 3,
						},
					},
					{
						RawTree: models.RawTree{
							Name:  "child1.2",
							Depth: 3,
						},
						Children: []models.Tree{
							{
								RawTree: models.RawTree{
									Name:  "child1.2.1",
									Depth: 4,
								},
							},
						},
					},
				},
			},
			{
				RawTree: models.RawTree{
					Name:  "child2",
					Depth: 2,
				},
				Children: []models.Tree{
					{
						RawTree: models.RawTree{
							Name:  "child2.1",
							Depth: 3,
						},
					},
				},
			},
		},
	}

	tmp := data.BuildRawTree()
	res := utils.Do(tmp, tmp[0].Depth)
	for _, t := range res {
		fmt.Println(t)
	}
}
