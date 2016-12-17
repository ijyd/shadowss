package factory

import (
	"gofreezer/pkg/storage"
	"gofreezer/pkg/storage/mysqls/mysql"
	"gofreezer/pkg/storage/storagebackend"

	_ "github.com/go-sql-driver/mysql"
	dbmysql "github.com/jinzhu/gorm"
)

//connectionStr: user:password@tcp(host:port)/dbname
func newMysqlClient(connectionStr string) (*dbmysql.DB, error) {
	var err error
	connStr := string(connectionStr) + string("?parseTime=True")
	//connStr := string(connectionStr)
	db, err := dbmysql.Open(string("mysql"), connStr)
	if err != nil {
		return nil, err
	}
	//db = db.Debug()

	return db, db.DB().Ping()
}

func newMysqlStorage(c storagebackend.Config) (storage.Interface, DestroyFunc, error) {
	endpoints := c.Mysql.ServerList

	client, err := newMysqlClient(endpoints[0])
	if err != nil {
		return nil, nil, err
	}

	destroyFunc := func() {
		client.Close()
	}

	return mysql.New(client, c.Codec, c.StorageVersion), destroyFunc, nil
}
