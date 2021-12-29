package mysql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/udev-21/nested-set-go/internal/models"
	"github.com/udev-21/nested-set-go/internal/repository"
	"github.com/udev-21/nested-set-go/internal/utils"
)

type Category struct {
	db *sqlx.DB
}

var _ repository.Category = (*Category)(nil)

func NewCategoryRepo(db *sqlx.DB) *Category {
	return &Category{
		db: db,
	}
}

func (r *Category) FetchByID(ctx context.Context, id int64) (models.Category, error) {
	if err := r.db.Ping(); err != nil {
		return models.Category{}, err
	}
	var node models.Category
	err := r.db.GetContext(ctx, &node, "SELECT * FROM category WHERE category_id = ?", id)
	if err != nil {
		return models.Category{}, err
	}

	return node, nil
}

func (r *Category) FetchAllChildren(ctx context.Context, parent models.Category) ([]models.Category, error) {
	sql := `SELECT * FROM category WHERE lft BETWEEN ? AND ? AND root_id = ? ORDER BY lft`
	return r.fetchAll(ctx, sql, parent.Left, parent.Right, parent.RootID)
}

func (r *Category) FetchAllChildrenWithDepth(ctx context.Context, parent models.Category, depth int64) ([]models.Category, error) {
	sql := `SELECT * FROM category WHERE lft BETWEEN ? AND ? AND root_id = ? AND depth <= ? ORDER BY lft`
	return r.fetchAll(ctx, sql, parent.Left, parent.Right, parent.RootID, parent.Depth+depth)
}

func (r *Category) FetchLeafs(ctx context.Context, node models.Category) ([]models.Category, error) {
	sql := `SELECT * FROM category WHERE lft BETWEEN ? AND ? AND root_id = ? AND lft + 1 = rgt ORDER BY lft`
	return r.fetchAll(ctx, sql, node.Left, node.Right, node.RootID)
}

func (r *Category) AppendRoot(ctx context.Context, node models.Category) error {
	if err := r.db.Ping(); err != nil {
		return err
	}

	tx := r.db.MustBeginTx(ctx, nil)
	raws := node.BuildRaw()
	cnt := len(raws)
	// INSERT INTO category(name, lft, rgt) VALUES('AFTER TELEVISIONS', @myRight + 1, @myRight + 2);
	sql := ""
	for idx, raw := range raws {
		sql += fmt.Sprintf("(%q, %d, %d, %d, %q)", raw.Name, raw.Left, raw.Right, raw.Depth, raw.RootID)
		if idx != cnt-1 {
			sql += ","
		}
	}

	_, err := tx.ExecContext(ctx, "INSERT INTO category (name, lft, rgt, depth, root_id) VALUES "+sql)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *Category) AppendChildren(ctx context.Context, parentNode models.Category, children []models.Category) error {
	if err := r.db.Ping(); err != nil {
		return err
	}
	if len(children) == 0 {
		return nil
	}
	tx := r.db.MustBeginTx(ctx, nil)
	// -- add new node inside specific node (after all, if children exists):
	rawChilds := utils.BuildRaw(children)
	childCnt := len(rawChilds)
	width := childCnt * 2
	// UPDATE category SET rgt = rgt + 2 WHERE rgt > @myLeft;
	_, err := tx.ExecContext(ctx, "UPDATE category SET rgt = rgt + ? WHERE rgt > ? AND root_id = ?", width, parentNode.Left, parentNode.RootID)
	if err != nil {
		tx.Rollback()
		return err
	}
	// UPDATE category SET lft = lft + 2 WHERE lft > @myLeft;
	_, err = tx.ExecContext(ctx, "UPDATE category SET lft = lft + ? WHERE lft > ? AND root_id = ?", width, parentNode.Left, parentNode.RootID)
	if err != nil {
		tx.Rollback()
		return err
	}
	// INSERT INTO category(name, lft, rgt) VALUES('FRS', @myLeft + 1, @myLeft + 2);
	sql := ""
	for idx, node := range rawChilds {
		sql += fmt.Sprintf("(%q, %d, %d, %d, %q)", node.Name, node.Left+parentNode.Left, node.Right+parentNode.Left, node.Depth+parentNode.Depth, parentNode.RootID)
		if idx != childCnt-1 {
			sql += ","
		}
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO category (name, lft, rgt, depth, root_id) VALUES "+sql)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *Category) AppendAfter(ctx context.Context, afterNode models.Category, children []models.Category) error {
	if err := r.db.Ping(); err != nil {
		return err
	}

	if afterNode.IsRoot() {
		return fmt.Errorf("after node must not be root node")
	}

	if len(children) == 0 {
		return nil
	}

	rawChilds := utils.BuildRaw(children)
	childCnt := len(rawChilds)
	width := childCnt * 2

	tx := r.db.MustBeginTx(ctx, nil)
	// UPDATE category SET rgt = rgt + 2 WHERE rgt > @myRight;
	_, err := tx.ExecContext(ctx, "UPDATE category SET rgt = rgt + ? WHERE rgt > ? AND root_id = ?", width, afterNode.Right, afterNode.RootID)
	if err != nil {
		tx.Rollback()
		return err
	}
	// UPDATE category SET lft = lft + 2 WHERE lft > @myRight;
	_, err = tx.ExecContext(ctx, "UPDATE category SET lft = lft + ? WHERE lft > ? AND root_id = ?", width, afterNode.Right, afterNode.RootID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// INSERT INTO category(name, lft, rgt) VALUES('AFTER TELEVISIONS', @myRight + 1, @myRight + 2);
	sql := ""
	for idx, node := range rawChilds {
		sql += fmt.Sprintf("(%q, %d, %d, %d, %q)", node.Name, node.Left+afterNode.Right, node.Right+afterNode.Right, node.Depth+afterNode.Depth-1, afterNode.RootID)
		if idx != childCnt-1 {
			sql += ","
		}
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO category (name, lft, rgt, depth, root_id) VALUES "+sql)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *Category) AppendBefore(ctx context.Context, beforeNode models.Category, children []models.Category) error {
	if err := r.db.Ping(); err != nil {
		return err
	}

	if beforeNode.IsRoot() {
		return fmt.Errorf("after node must not be root node")
	}

	if len(children) == 0 {
		return nil
	}

	rawChilds := utils.BuildRaw(children)
	childCnt := len(rawChilds)
	width := childCnt * 2

	tx := r.db.MustBeginTx(ctx, nil)
	// UPDATE category SET rgt = rgt + 2 WHERE rgt >= @myLeft;
	_, err := tx.ExecContext(ctx, "UPDATE category SET rgt = rgt + ? WHERE rgt >= ? AND root_id = ?", width, beforeNode.Left, beforeNode.RootID)
	if err != nil {
		tx.Rollback()
		return err
	}
	// UPDATE category SET lft = lft + 2 WHERE lft >= @myLeft;
	_, err = tx.ExecContext(ctx, "UPDATE category SET lft = lft + ? WHERE lft >= ? AND root_id = ?", width, beforeNode.Left, beforeNode.RootID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// INSERT INTO category(name, lft, rgt) VALUES('BEFORE TELEVISIONS', @myLeft, @myLeft + 1);
	sql := ""
	for idx, node := range rawChilds {
		sql += fmt.Sprintf("(%q, %d, %d, %d, %q)", node.Name, node.Left+beforeNode.Left, node.Right+beforeNode.Left, node.Depth+beforeNode.Depth-1, beforeNode.RootID)
		if idx != childCnt-1 {
			sql += ","
		}
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO category (name, lft, rgt, depth, root_id) VALUES "+sql)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *Category) MoveAfter(ctx context.Context, afterNode, target models.Category) error {
	if err := r.db.Ping(); err != nil {
		return err
	}

	if afterNode.IsRoot() {
		return fmt.Errorf("after node must not be root node")
	} else if afterNode.RootID != target.RootID {
		return fmt.Errorf("afterNode and target must be common root_id")
	}

	width := target.Right - target.Left + 1
	// SELECT rgt, lft, rgt - lft + 1
	// INTO @myRgt, @myLft, @myWidth
	// FROM category
	// WHERE name = 'PORTABLE ELECTRONICS';
	tx := r.db.MustBeginTx(ctx, nil)
	// SELECT rgt, lft, rgt - lft + 1
	// INTO @myRgt, @myLft, @myWidth
	// FROM category
	// WHERE name = 'TELEVISIONS';

	// SELECT GROUP_CONCAT(category_id) INTO @tmpIds FROM category where lft >= @myLft && rgt <= @myRgt;
	_, err := tx.ExecContext(ctx, "SELECT GROUP_CONCAT(category_id) INTO @tmpIds FROM category WHERE lft >= ? AND rgt <= ? AND root_id = ?", target.Left, target.Right, target.RootID)
	if err != nil {
		return err
	}

	// UPDATE category SET lft = lft - @myWidth WHERE lft > @myRgt AND NOT FIND_IN_SET(category_id, @tmpIds) ;
	_, err = tx.ExecContext(ctx, "UPDATE category SET lft = lft - ? WHERE lft > ? AND NOT FIND_IN_SET(category_id, @tmpIds) AND root_id = ?", width, target.Right, target.RootID)
	if err != nil {
		return err
	}
	// UPDATE category SET rgt = rgt - @myWidth WHERE rgt > @myRgt AND NOT FIND_IN_SET(category_id, @tmpIds) ;
	_, err = tx.ExecContext(ctx, "UPDATE category SET rgt = rgt - ? WHERE rgt > ? AND NOT FIND_IN_SET(category_id, @tmpIds) AND root_id = ?", width, target.Right, target.RootID)
	if err != nil {
		return err
	}
	// SELECT rgt INTO @afterNodeRight
	// FROM category
	// WHERE name = 'AFTER PORTABLE ELECTRONICS';
	err = tx.GetContext(ctx, &afterNode, "SELECT * FROM category WHERE category_id = ?", afterNode.ID)
	if err != nil {
		return err
	}

	// UPDATE category SET lft = lft + @myWidth WHERE lft > @afterNodeRight AND NOT FIND_IN_SET(category_id, @tmpIds) ;
	_, err = tx.ExecContext(ctx, "UPDATE category SET lft = lft + ? WHERE lft > ? AND NOT FIND_IN_SET(category_id, @tmpIds) AND root_id = ?", width, afterNode.Right, target.RootID)
	if err != nil {
		return err
	}

	// UPDATE category SET rgt = rgt + @myWidth WHERE rgt > @afterNodeRight AND NOT FIND_IN_SET(category_id, @tmpIds) ;
	_, err = tx.ExecContext(ctx, "UPDATE category SET rgt = rgt + ? WHERE rgt > ? AND NOT FIND_IN_SET(category_id, @tmpIds) AND root_id = ?", width, afterNode.Right, target.RootID)
	if err != nil {
		return err
	}
	// SELECT @afterNodeRight - @myLft + 1 INTO @diff;
	diff := afterNode.Right - target.Left + 1

	diffDepth := afterNode.Depth - target.Depth

	// SELECT @diff;
	// UPDATE category SET lft = lft + @diff, rgt = rgt + @diff WHERE FIND_IN_SET(category_id, @tmpIds) ;
	_, err = tx.ExecContext(ctx, "UPDATE category SET lft = lft + ?, rgt = rgt + ?, depth = depth + ? WHERE FIND_IN_SET(category_id, @tmpIds)  AND root_id = ?", diff, diff, diffDepth, target.RootID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *Category) MoveBefore(ctx context.Context, beforeNode, target models.Category) error {
	if err := r.db.Ping(); err != nil {
		return err
	}

	if beforeNode.IsRoot() {
		return fmt.Errorf("after node must not be root node")
	} else if beforeNode.RootID != target.RootID {
		return fmt.Errorf("beforeNode and target must be common root_id")
	}

	width := target.Right - target.Left + 1
	// SELECT rgt, lft, rgt - lft + 1
	// INTO @myRgt, @myLft, @myWidth
	// FROM category
	// WHERE name = 'PORTABLE ELECTRONICS';
	tx := r.db.MustBeginTx(ctx, nil)

	// SELECT GROUP_CONCAT(category_id) INTO @tmpIds FROM category where lft >= @myLft && rgt <= @myRgt;
	_, err := tx.ExecContext(ctx, "SELECT GROUP_CONCAT(category_id) INTO @tmpIds FROM category WHERE lft >= ? AND rgt <= ? AND root_id = ?", target.Left, target.Right, target.RootID)
	if err != nil {
		return err
	}

	// UPDATE category SET lft = lft - @myWidth WHERE lft > @myRgt AND NOT FIND_IN_SET(category_id, @tmpIds) ;
	_, err = tx.ExecContext(ctx, "UPDATE category SET lft = lft - ? WHERE lft > ? AND NOT FIND_IN_SET(category_id, @tmpIds) AND root_id = ?", width, target.Right, target.RootID)
	if err != nil {
		return err
	}
	// UPDATE category SET rgt = rgt - @myWidth WHERE rgt > @myRgt AND NOT FIND_IN_SET(category_id, @tmpIds) ;
	_, err = tx.ExecContext(ctx, "UPDATE category SET rgt = rgt - ? WHERE rgt > ? AND NOT FIND_IN_SET(category_id, @tmpIds) AND root_id = ?", width, target.Right, target.RootID)
	if err != nil {
		return err
	}
	// SELECT lft INTO @beforeLeft
	// FROM category
	// WHERE name = 'TELEVISIONS';
	err = tx.GetContext(ctx, &beforeNode, "SELECT * FROM category WHERE category_id = ?", beforeNode.ID)
	if err != nil {
		return err
	}

	// UPDATE category SET lft = lft + @myWidth WHERE lft >= @beforeLeft AND NOT FIND_IN_SET(category_id, @tmpIds);
	_, err = tx.ExecContext(ctx, "UPDATE category SET lft = lft + ? WHERE lft >= ? AND NOT FIND_IN_SET(category_id, @tmpIds) AND root_id = ?", width, beforeNode.Left, target.RootID)
	if err != nil {
		return err
	}
	// UPDATE category SET rgt = rgt + @myWidth WHERE rgt >= @beforeLeft AND NOT FIND_IN_SET(category_id, @tmpIds);
	_, err = tx.ExecContext(ctx, "UPDATE category SET rgt = rgt + ? WHERE rgt >= ? AND NOT FIND_IN_SET(category_id, @tmpIds) AND root_id = ?", width, beforeNode.Left, target.RootID)
	if err != nil {
		return err
	}
	// SELECT @beforeLeft-@myLft INTO @diff;
	diff := beforeNode.Left - target.Left

	diffDepth := beforeNode.Depth - target.Depth

	// UPDATE category SET lft = lft + @diff, rgt = rgt + @diff WHERE FIND_IN_SET(category_id, @tmpIds) ;
	_, err = tx.ExecContext(ctx, "UPDATE category SET lft = lft + ?, rgt = rgt + ?, depth = depth + ? WHERE FIND_IN_SET(category_id, @tmpIds)  AND root_id = ?", diff, diff, diffDepth, target.RootID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *Category) MoveInto(ctx context.Context, beforeNode, target models.Category) error {
	return nil
}

func (r *Category) Delete(ctx context.Context, node models.Category) error {
	if err := r.db.Ping(); err != nil {
		return err
	} else if node.IsRoot() {
		return fmt.Errorf("can't delete just root, use deleteWithChildren method")
	}

	// DELETE FROM tree WHERE lft = @myLeft;
	tx := r.db.MustBeginTx(ctx, nil)
	_, err := tx.ExecContext(ctx, "DELETE FROM category WHERE lft = ? AND root_id = ?", node.Left, node.RootID)

	if err != nil {
		tx.Rollback()
		return err
	}

	// UPDATE tree SET rgt = rgt - 1, lft = lft - 1 WHERE lft BETWEEN @myLeft AND @myRight;
	_, err = tx.ExecContext(ctx, "UPDATE category SET rgt = rgt - 1, lft = lft - 1, depth = depth - 1 WHERE lft BETWEEN ? AND ? AND root_id = ?", node.Left, node.Right, node.RootID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// UPDATE tree SET rgt = rgt - 2 WHERE rgt > @myRight;
	_, err = tx.ExecContext(ctx, "UPDATE category SET rgt = rgt - 2 WHERE rgt > ? AND root_id = ?", node.Right, node.RootID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// UPDATE tree SET lft = lft - 2 WHERE lft > @myRight;
	_, err = tx.ExecContext(ctx, "UPDATE category SET lft = lft - 2 WHERE lft > ? AND root_id = ?", node.Right, node.RootID)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *Category) DeleteWithChildren(ctx context.Context, node models.Category) error {

	// SELECT @myLeft := lft, @myRight := rgt, @myWidth := rgt - lft + 1
	// FROM tree
	// WHERE name = 'GAME CONSOLES';

	// DELETE FROM tree WHERE lft BETWEEN @myLeft AND @myRight;
	tx := r.db.MustBeginTx(ctx, nil)
	_, err := tx.ExecContext(ctx, "DELETE FROM category WHERE lft BETWEEN ? AND ? AND root_id = ?", node.Left, node.Right, node.RootID)
	if err != nil {
		tx.Rollback()
		return err
	}
	var width = node.Right - node.Left + 1
	// UPDATE tree SET rgt = rgt - @myWidth WHERE rgt > @myRight;
	_, err = tx.ExecContext(ctx, "UPDATE category SET rgt = rgt - ? WHERE rgt > ? AND root_id = ?", width, node.Right, node.RootID)
	if err != nil {
		tx.Rollback()
		return err
	}
	// UPDATE tree SET lft = lft - @myWidth WHERE lft > @myRight;
	_, err = tx.ExecContext(ctx, "UPDATE category SET lft = lft - ? WHERE lft > ? AND root_id = ?", width, node.Right, node.RootID)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *Category) RecalculateDepth(ctx context.Context) error {
	sql := `				
			WITH tmp(category_id, depth)
			as (
				SELECT category_id, 
					(
						select COUNT(*) + 1 FROM category where lft < c.lft AND rgt > c.rgt
					) as depth 
				FROM category as c
			)
			UPDATE category
			JOIN tmp USING(category_id)
			SET category.depth = tmp.depth;`
	_, err := r.db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func (r *Category) GetAllRootIds(ctx context.Context) ([]string, error) {
	if err := r.db.Ping(); err != nil {
		return []string{}, err
	}
	rows, err := r.db.Query("SELECT DISTINCT root_id FROM category")
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()
	var ids = []string{}
	var rootId string
	for rows.Next() {
		if err := rows.Scan(&rootId); err != nil {
			return []string{}, err
		}
		ids = append(ids, rootId)
	}
	return ids, nil
}

func (r *Category) fetchAll(ctx context.Context, sql string, args ...interface{}) ([]models.Category, error) {
	var children []models.Category
	err := r.db.SelectContext(ctx, &children, sql, args...)
	if err != nil {
		return []models.Category{}, err
	}
	return children, nil
}
