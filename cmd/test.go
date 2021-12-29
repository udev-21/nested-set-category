package main

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/udev-21/nested-set-go/internal/models"
	"github.com/udev-21/nested-set-go/internal/repository/mysql"
	"github.com/udev-21/nested-set-go/internal/utils"

	"github.com/jmoiron/sqlx"
)

func test(db *sqlx.DB, treeRepo *mysql.Category) {
	data := []models.Category{
		{
			Name: "child1",
			Children: []models.Category{
				{
					Name: "child1.1",
					Children: []models.Category{
						{
							Name: "child1.1.1",
						},
					},
				},
			},
		},
		{
			Name: "child2",
			Children: []models.Category{
				{
					Name: "child2.1",
					Children: []models.Category{
						{
							Name: "child2.1.1",
						},
					},
				},
				{
					Name: "child2.2",
					Children: []models.Category{
						{
							Name: "child2.2.1",
						},
					},
				},
			},
		},
	}
	parent := models.Category{
		Name:     "root",
		Children: data,
		RootID:   "23",
	}
	treeRepo.AppendRoot(context.Background(), parent)
	// parent, _ := treeRepo.FetchByID(context.Background(), 12)
	// err := treeRepo.Delete(context.Background(), parent)
	// if err != nil {
	// 	panic(err)
	// }
	parent, _ = treeRepo.FetchByID(context.Background(), 12)
	childs, _ := treeRepo.FetchAllChildren(context.Background(), parent)
	fmt.Println(utils.PrintNestedCategory(utils.BuildNestedCategory(childs), ""))
	// fmt.Println("ok")
}

func testMoveBefore(db *sqlx.DB, treeRepo *mysql.Category) {
	target, _ := treeRepo.FetchByID(context.Background(), 2)
	beforeNode, err := treeRepo.FetchByID(context.Background(), 10)
	err = treeRepo.MoveBefore(context.Background(), beforeNode, target)
	if err != nil {
		panic(err)
	}
}

func testMoveAfter(db *sqlx.DB, treeRepo *mysql.Category) {
	target, _ := treeRepo.FetchByID(context.Background(), 11)
	afterNode, err := treeRepo.FetchByID(context.Background(), 10)
	err = treeRepo.MoveAfter(context.Background(), afterNode, target)
	if err != nil {
		panic(err)
	}
}

func testMoveInto(db *sqlx.DB, treeRepo *mysql.Category) {
	target, _ := treeRepo.FetchByID(context.Background(), 6)
	parentNode, err := treeRepo.FetchByID(context.Background(), 11)
	err = treeRepo.MoveInto(context.Background(), parentNode, target)
	if err != nil {
		panic(err)
	}
}

func main() {
	db := ConnectDB()
	treeRepo := mysql.NewCategoryRepo(db)
	if err := treeRepo.RecalculateDepth(context.Background()); err != nil {
		fmt.Println("RecalculateDepth: ", err.Error())
	}
	// test(db, treeRepo)
	testMoveBefore(db, treeRepo)
	// testMoveAfter(db, treeRepo)
	// testMoveInto(db, treeRepo)
	parent, _ := treeRepo.FetchByID(context.Background(), 1)
	childs, _ := treeRepo.FetchAllChildren(context.Background(), parent)
	fmt.Println(utils.PrintNestedCategory(utils.BuildNestedCategory(childs), ""))
	return

	node, _ := treeRepo.FetchByID(context.Background(), 1)
	children, _ := treeRepo.FetchLeafs(context.Background(), node)
	for _, v := range children {
		fmt.Println(v)
	}
	fmt.Println(utils.PrintNestedCategory(utils.BuildNestedCategory(children), ""))
	defer db.Close()
}

func ConnectDB() *sqlx.DB {
	user := "gav"
	pass := "Sok0l"
	host := "127.0.0.1"
	port := "3366"
	db := "parsers"
	return sqlx.MustConnect("mysql", user+":"+pass+"@("+host+":"+port+")/"+db+"?charset=utf8mb4,utf8&parseTime=true")
}
