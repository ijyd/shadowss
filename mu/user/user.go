package user

var (
	client = NewMysqlClient()
)

func GetClient() Client {
	return client
}

func SetClient(c Client) {
	client = c
}

type Client interface {
	GetUsers() ([]User, error)
}

type User interface {
	GetPort() int
	GetPasswd() string
	UpdatetTraffic() error
}
