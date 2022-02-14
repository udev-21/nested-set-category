package utils

import (
	"sort"
	"strconv"

	"github.com/udev-21/nested-set-go/models"
)

// BuildNestedCategory builds nestedCategory from nodes that ordered by lft asc
func BuildNestedCategory(nodes []models.Category) []models.Category {
	if len(nodes) == 0 {
		return []models.Category{}
	}
	return _buildNestedCategory(nodes, nodes[0].Depth)
}

func BuildRaw(nodes []models.Category) []models.Category {
	if len(nodes) == 0 {
		return []models.Category{}
	}
	tmp := _BuildRaw(nodes, 1, 1)
	sort.SliceStable(tmp, func(i, j int) bool {
		return tmp[i].Left < tmp[j].Left
	})
	return tmp
}

func PrintNestedCategory(node []models.Category, prefix string) string {
	res := ""
	for _, node := range node {
		res += prefix + strconv.Itoa(int(node.ID)) + " " + node.Name + " " + strconv.Itoa(int(node.Left)) + " " + strconv.Itoa(int(node.Right))
		res += "\n" + PrintNestedCategory(node.Children, prefix+"  ")
	}
	return res
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

func _BuildRaw(nodes []models.Category, left, depth int64) []models.Category {
	res := make([]models.Category, 0)
	for _, child := range nodes {
		child.Left = left
		child.Depth = depth
		childs := _BuildRaw(child.Children, left+1, depth+1)
		if len(childs) > 0 {
			res = append(res, childs...)
			left = childs[len(childs)-1].Right
		}
		child.Right = left + 1
		res = append(res, child)
		left = left + 2
	}
	return res
}
