package utils

import (
	"github.com/udev-21/nested-set-go/pkg/models"
)

// BuildNestedCategory builds nestedCategory from nodes that ordered by lft asc
func BuildNestedCategory(nodes []models.Category) []models.Category {
	if len(nodes) == 0 {
		return []models.Category{}
	}
	return _buildNestedCategory(nodes, nodes[0].Depth)
}

func _buildNestedCategory(nodes []models.Category, depth int64) []models.Category {
	if len(nodes) == 0 {
		return []models.Category{}
	}
	var res []models.Category
	for idx, node := range nodes {
		if node.Depth == depth {
			childNode := node
			childLines := []models.Category{}
			for i := idx + 1; i < len(nodes); i++ {
				if nodes[i].Depth == depth {
					break
				}
				childLines = append(childLines, nodes[i])
			}
			childNode.Children = _buildNestedCategory(childLines, depth+1)
			res = append(res, childNode)
		}
	}
	return res
}

func PrintNestedCategory(node []models.Category, prefix string) string {
	res := ""
	for _, node := range node {
		res += prefix + node.Name
		res += "\n" + PrintNestedCategory(node.Children, prefix+"  ")
	}
	return res
}
