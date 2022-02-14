package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/udev-21/nested-set-go/config/database/mysql"
)

func NewMysqlDatabase(cnf *mysql.Config) *sqlx.DB {
	return sqlx.MustConnect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cnf.User, cnf.Password, cnf.Host, cnf.Port, cnf.DBName))
}
