package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var (
	tableName string
)

func genConnStr(user, password, host, dbname string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=True", user, password, host, dbname)
}

func (c *Client) Boot(dbType, user, password, host, dbname string) error {
	var err error
	db, err := gorm.Open(dbType, genConnStr(user, password, host, dbname))
	if err != nil {
		return err
	}
	c.SetDb(&db)
	c.db.DB()
	return c.db.DB().Ping()
}

func (c *Client) SetTable(table string) {
	tableName = table
	c.table = table
}
