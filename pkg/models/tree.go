package models

import "sort"

type Tree struct {
	RawTree
	Children []Tree
}

func (node *Tree) BuildRawTree() []RawTree {
	tmp := node._BuildRawTree(1, 1)
	sort.SliceStable(tmp, func(i, j int) bool {
		return tmp[i].Left < tmp[j].Left
	})
	return tmp
}

func (node *Tree) _BuildRawTree(left, depth int64) []RawTree {
	res := make([]RawTree, 0)

	parent := RawTree{
		Name:  node.Name,
		Left:  left,
		Right: left + 1,
		Depth: depth,
	}

	for _, child := range node.Children {
		childs := child._BuildRawTree(left+1, depth+1)
		if len(childs) > 0 {
			res = append(res, childs...)
			left = childs[len(childs)-1].Right
		}
	}
	if len(res) > 0 {
		parent.Right = res[len(res)-1].Right + 1
	}
	res = append(res, parent)

	return res
}
