package user

func NewMysqlClient() *MysqlClient {
	var client *MysqlClient
	return client
}

type MysqlClient struct {
}

type MysqlUser struct {
	port   int
	passwd string
	method string
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

func (c *MysqlClient) GetClient() ([]User, error) {
	var users []MysqlUser
	return users, error()
}
