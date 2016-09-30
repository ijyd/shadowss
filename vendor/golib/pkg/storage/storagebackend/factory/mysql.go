package factory

import (
	"golib/pkg/storage"
	"golib/pkg/storage/mysql"
	"golib/pkg/storage/storagebackend"

	_ "github.com/go-sql-driver/mysql"
	dbmysql "github.com/jinzhu/gorm"
)

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

func newMysqlStorage(c storagebackend.Config) (storage.Interface, error) {
	endpoints := c.ServerList

	client, err := newMysqlClient(endpoints[0])
	if err != nil {
		return nil, err
	}

	return mysql.New(client), nil
}
