package mysql

import (
	"github.com/jinzhu/gorm"
	"github.com/orvice/shadowsocks-go/mu/log"
	"github.com/orvice/shadowsocks-go/mu/user"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
)

var client *Client

func SetClient(c *Client){
	client = c
}

func NewClient() *Client {
	mclient := new(Client)
	return mclient
}

type Client struct {
	db    *gorm.DB
	table string
}

type User struct {
	port           int
	passwd         string
	method         string
	enable         int
	transferEnable int
	u              int
	d              int
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

func (u *User) IsEnable() bool {
	if u.enable == 0 {
		return false
	}
	if u.u+u.d > u.transferEnable {
		return false
	}
	return true
}

func (u *User) GetCipher() (*ss.Cipher, error) {
	return ss.NewCipher(u.method, u.passwd)
}

func (u *User) UpdatetTraffic(storageSize int) error {
	return client.db.Model(u).Update("d", gorm.Expr("d  + ?", storageSize)).Error
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
	rows, err := c.db.Model(User{}).Select("passwd, port, method,enable,transfer_enable,u,d").Rows()
	if err != nil {
		log.Log.Error(err)
		var users []user.User
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		var data User
		err := rows.Scan(&data.passwd, &data.port, &data.method, &data.enable, &data.transferEnable, &data.u, &data.d)
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
