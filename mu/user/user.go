package user

var (
	client = NewMysqlClient()
)

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
