package factory

import (
	"shadowsocks-go/pkg/storage/mysql"

	"shadowsocks-go/pkg/storage"
	"shadowsocks-go/pkg/storage/storagebackend"

	_ "github.com/go-sql-driver/mysql"
	dbmysql "github.com/jinzhu/gorm"
)

// func genConnStr(user, password, host, dbname string) string {
// 	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=True", user, password, host, dbname)
// }

//connectionStr: user:password@tcp(host:port)/dbname
func newMysqlClient(connectionStr string) (*dbmysql.DB, error) {
	var err error
	connStr := string(connectionStr) + string("?parseTime=True")
	db, err := dbmysql.Open(string("mysql"), connStr)
	if err != nil {
		return nil, err
	}

	return db, db.DB().Ping()
}

// func (c *Client) SetTable(table string) {
// 	tableName = table
// 	c.table = table
// }

func newMysqlStorage(c storagebackend.Config) (storage.Interface, error) {
	endpoints := c.ServerList

	client, err := newMysqlClient(endpoints[0])
	if err != nil {
		return nil, err
	}

	return mysql.New(client), nil
}
