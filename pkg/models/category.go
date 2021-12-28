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
	return fmt.Sprintf("name: %s, depth: %d, children: %d", t.Name, t.Depth, len(t.Children))
}

func (node *Category) BuildRaw() []Category {
	tmp := node._BuildRaw(1, 1)
	sort.SliceStable(tmp, func(i, j int) bool {
		return tmp[i].Left < tmp[j].Left
	})
	return tmp
}

func (node *Category) _BuildRaw(left, depth int64) []Category {
	res := make([]Category, 0)

	parent := node

	for _, child := range node.Children {
		childs := child._BuildRaw(left+1, depth+1)
		if len(childs) > 0 {
			res = append(res, childs...)
			left = childs[len(childs)-1].Right
		}
	}
	if len(res) > 0 {
		parent.Right = res[len(res)-1].Right + 1
	}

	return append(res, *parent)
}
