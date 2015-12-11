package user

import (
	"github.com/jinzhu/gorm"
)

func NewMysqlClient() *MysqlClient {
	var client *MysqlClient
	return client
}

type MysqlClient struct {
	db    *gorm.DB
	table string
}

type MysqlUser struct {
	port   int
	passwd string
	method string
}

func (c MysqlUser) TableName() string {
	return tableName
}

func (u *MysqlUser) GetPort() int {
	return u.port
}

func (u *MysqlUser) GetPasswd() string {
	return u.passwd
}

func (u *MysqlUser) GetMethod() string {
	return u.method
}

func (u *MysqlUser) UpdatetTraffic() error {
	return nil
}

func (c *MysqlClient) GetUsers() ([]User, error) {
	var datas []*MysqlUser
	rows, err := c.db.Model(MysqlUser{}).Where("enable = ?", "1").Select("passwd, prot, method").Rows()
	if err != nil {
		var users []User
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		var data *MysqlUser
		rows.Scan(data.passwd, data.port, data.method)
		datas = append(datas, data)
	}
	users := make([]User, len(datas))
	for k, v := range datas {
		users[k] = v
	}
	return users, nil
}
