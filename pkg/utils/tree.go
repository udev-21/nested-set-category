package utils

import (
	"github.com/udev-21/nested-set-go/pkg/models"
)

func BuildNestedTree(nodes []models.RawTree, level int64) []models.Tree {
	if len(nodes) == 0 {
		return []models.Tree{}
	}
	var res []models.Tree
	for idx, node := range nodes {
		if node.Depth == level {
			childNode := models.Tree{RawTree: models.RawTree{Name: node.Name}}
			childLines := []models.RawTree{}
			for i := idx + 1; i < len(nodes); i++ {
				if nodes[i].Depth == level {
					break
				}
				childLines = append(childLines, nodes[i])
			}
			childNode.Children = BuildNestedTree(childLines, level+1)
			res = append(res, childNode)
		}
	}
	return res
}
