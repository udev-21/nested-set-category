package mysql

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/udev-21/nested-set-go/internal/repository"
	"github.com/udev-21/nested-set-go/pkg/models"
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
	sql := `
		SELECT * FROM category WHERE lft BETWEEN ? AND ? AND root_id = ? ORDER BY lft
	`
	var children []models.Category
	err := r.db.SelectContext(ctx, &children, sql, parent.Left, parent.Right, parent.RootID)
	if err != nil {
		return []models.Category{}, err
	}
	return children, nil
}

func (r *Category) FetchAllChildrenWithDepth(ctx context.Context, parent models.Category, depth int64) ([]models.Category, error) {
	sql := `
		SELECT * FROM category WHERE lft BETWEEN ? AND ? AND root_id = ? AND depth <= ? ORDER BY lft
	`
	var children []models.Category
	err := r.db.SelectContext(ctx, &children, sql, parent.Left, parent.Right, parent.RootID, parent.Depth+depth)
	if err != nil {
		return []models.Category{}, err
	}
	return children, nil
}

func (r *Category) FetchLeafs(ctx context.Context, node models.Category) ([]models.Category, error) {
	return []models.Category{}, nil
}

func (r *Category) AppendRoot(ctx context.Context, node models.Category) error {
	return nil
}

func (r *Category) AppendChildren(ctx context.Context, parentNode models.Category, children models.Category) error {
	return nil
}

func (r *Category) AppendAfter(ctx context.Context, afterNode models.Category, children models.Category) error {
	return nil
}

func (r *Category) AppendBefore(ctx context.Context, beforeNode models.Category, children models.Category) error {
	return nil
}

func (r *Category) MoveAfter(ctx context.Context, afterNode, target models.Category) error {
	return nil
}

func (r *Category) MoveBefore(ctx context.Context, beforeNode, target models.Category) error {
	return nil
}

func (r *Category) Delete(ctx context.Context, node models.Category) error {

	// SELECT @myLeft := lft, @myRight := rgt, @myWidth := rgt - lft + 1
	// FROM tree
	// WHERE name = 'PORTABLE ELECTRONICS';

	// DELETE FROM tree WHERE lft = @myLeft;
	tx := r.db.MustBeginTx(ctx, nil)
	_, err := tx.ExecContext(ctx, "DELETE FROM category WHERE lft = ? AND root_id = ?", node.Left, node.RootID)

	if err != nil {
		tx.Rollback()
		return err
	}

	// UPDATE tree SET rgt = rgt - 1, lft = lft - 1 WHERE lft BETWEEN @myLeft AND @myRight;
	_, err = tx.ExecContext(ctx, "UPDATE category SET rgt = rgt - 1, lft = lft - 1 WHERE lft BETWEEN ? AND ? AND root_id = ?", node.Left, node.Right, node.RootID)
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
