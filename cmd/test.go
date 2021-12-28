package main

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/udev-21/nested-set-go/internal/repository/mysql"
	"github.com/udev-21/nested-set-go/pkg/utils"

	"github.com/jmoiron/sqlx"
)

func main() {
	db := ConnectDB()
	treeRepo := mysql.NewCategoryRepo(db)
	if err := treeRepo.RecalculateDepth(context.Background()); err != nil {
		fmt.Println("RecalculateDepth: ", err.Error())
	}

	node, _ := treeRepo.FetchByID(context.Background(), 1)
	children, _ := treeRepo.FetchAllChildren(context.Background(), node)
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
