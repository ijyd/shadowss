package user

var (
	client = NewClient()
)

func NewClient() Client {
	var client Client
	return client
}

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
	GetMethod() string
	UpdatetTraffic() error
}
