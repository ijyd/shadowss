package mysql

import (
	"github.com/jinzhu/gorm"
	"github.com/orvice/shadowsocks-go/mu/log"
	"github.com/orvice/shadowsocks-go/mu/user"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
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

func (c *Client) SetDb(db *gorm.DB) {
	c.db = db
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

func (u *User) GetCipher() (*ss.Cipher, error) {
	return ss.NewCipher(u.method, u.passwd)
}

func (u *User) UpdatetTraffic() error {
	return nil
}

func (u *User) GetUserInfo() user.UserInfo {
	return user.UserInfo{
		Passwd: u.passwd,
		Port:   u.port,
		Method: u.method,
	}
}

func (c *Client) GetUsers() ([]user.User, error) {
	log.Log.Info("get mysql users")
	var datas []*User
	rows, err := c.db.Model(User{}).Where("enable = ?", "1").Select("passwd, port, method").Rows()
	if err != nil {
		log.Log.Error(err)
		var users []user.User
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		var data User
		err := rows.Scan(&data.passwd, &data.port, &data.method)
		if err != nil {
			log.Log.Error(err)
			continue
		}
		datas = append(datas, &data)
	}
	log.Log.Info(len(datas))
	users := make([]user.User, len(datas))
	for k, v := range datas {
		users[k] = v
	}
	return users, nil
}
