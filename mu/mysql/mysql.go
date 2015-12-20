package mysql

import (
	"github.com/orvice/shadowsocks-go/mu/user"
	"github.com/jinzhu/gorm"
)

func NewClient() *Client {
	mclient := new(Client)
	return mclient
}

type Client struct {
	db    *gorm.DB
	table string
}

type User struct {
	port   int
	passwd string
	method string
}

func (c User) TableName() string {
	return tableName
}

func (u *User) GetPort() int {
	return u.port
}

func (u *User) GetPasswd() string {
	return u.passwd
}

func (u *User) GetMethod() string {
	return u.method
}

func (u *User) UpdatetTraffic() error {
	return nil
}

func (c *Client) GetUsers() ([]user.User, error) {
	var datas []*User
	rows, err := c.db.Model(User{}).Where("enable = ?", "1").Select("passwd, prot, method").Rows()
	if err != nil {
		var users []user.User
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		var data *User
		rows.Scan(data.passwd, data.port, data.method)
		datas = append(datas, data)
	}
	users := make([]user.User, len(datas))
	for k, v := range datas {
		users[k] = v
	}
	return users, nil
}
