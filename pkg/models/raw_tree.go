package models

import "fmt"

type RawTree struct {
	ID    int64  `db:"id"`
	Name  string `db:"name"`
	Left  int64  `db:"lft"`
	Right int64  `db:"rgt"`
	Depth int64  `db:"depth"`
}

func (r RawTree) String() string {
	return fmt.Sprintf("%s (%d, %d, %d)", r.Name, r.Left, r.Right, r.Depth)
}
