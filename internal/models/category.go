package models

import (
	"fmt"
	"sort"
)

type Category struct {
	ID       int64  `db:"category_id"`
	Name     string `db:"name"`
	Left     int64  `db:"lft"`
	Right    int64  `db:"rgt"`
	Depth    int64  `db:"depth"`
	RootID   string `db:"root_id"`
	Children []Category
}

func (t Category) String() string {
	return fmt.Sprintf("left: %d, right: %d, name: %s, depth: %d, children: %d ", t.Left, t.Right, t.Name, t.Depth, len(t.Children))
}

func (node *Category) BuildRaw() []Category {
	tmp := node._BuildRaw(1, 1, node.RootID)
	sort.SliceStable(tmp, func(i, j int) bool {
		return tmp[i].Left < tmp[j].Left
	})
	return tmp
}

func (node Category) IsRoot() bool {
	return node.Left == 1
}

func (node Category) IsLeaf() bool {
	return node.Left+1 == node.Right
}

func (node *Category) _BuildRaw(left, depth int64, rootId string) []Category {
	res := make([]Category, 0)

	parent := Category{
		Depth:  depth,
		Left:   left,
		Right:  left + 1,
		Name:   node.Name,
		RootID: rootId,
	}

	for _, child := range node.Children {
		childs := child._BuildRaw(left+1, depth+1, rootId)
		if len(childs) > 0 {
			res = append(res, childs...)
			left = childs[len(childs)-1].Right
		}
	}
	if len(res) > 0 {
		parent.Right = res[len(res)-1].Right + 1
	}

	return append(res, parent)
}
