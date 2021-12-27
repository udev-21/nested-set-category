package utils

import (
	"github.com/udev-21/nested-set-go/pkg/models"
)

func Do(nodes []models.RawTree, depth int64) []models.Tree {
	if len(nodes) == 0 {
		return []models.Tree{}
	}
	var res []models.Tree
	for idx, node := range nodes {
		if node.Depth == depth {
			childNode := models.Tree{RawTree: node}
			childLines := []models.RawTree{}
			for i := idx + 1; i < len(nodes); i++ {
				if nodes[i].Depth == depth {
					break
				}
				childLines = append(childLines, nodes[i])
			}
			childNode.Children = Do(childLines, depth+1)
			res = append(res, childNode)
		}
	}
	return res
}
