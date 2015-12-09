package user

func NewMysqlClient() Client {
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

func (u *MysqlUser) UpdatetTraffic() error {
	return nil
}

func (c *MysqlClient) GetUsers() ([]User, error) {
	var datas []*MysqlUser
	users := make([]User, len(datas))
	for k, v := range datas {
		users[k] = v
	}
	return users, nil
}
